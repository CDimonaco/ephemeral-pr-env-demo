package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cdimonaco/ephimeral-pr-env-demo/api/internal/http/handlers"
	"github.com/cdimonaco/ephimeral-pr-env-demo/api/internal/persistence"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"moul.io/chizap"
)

type envConfig struct {
	DBConnection string
}

func mustGetEnvConfig() envConfig {
	db := os.Getenv("DB_CONNECTION")
	if db == "" {
		panic("DB_CONNECTION env is required")
	}
	return envConfig{DBConnection: db}
}

func main() {
	envConfig := mustGetEnvConfig()

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Encoding = "console"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	ul, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}
	l := ul.Sugar()

	appCtx := context.Background()

	dbpool, err := pgxpool.New(appCtx, envConfig.DBConnection)
	if err != nil {
		l.Fatalf("Unable to create connection pool: %v\n", err)
	}

	defer dbpool.Close()

	recipesRepository := persistence.NewRecipesRepository(dbpool, l)
	recipeHandlers := handlers.NewRecipesHandler(l, recipesRepository)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(chizap.New(ul, &chizap.Opts{WithUserAgent: true}))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.AllowContentType("application/json"))

	r.Get("/recipes", recipeHandlers.GetAllRecipes)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, handlers.HttpError{
			Error:     "route not found",
			Code:      http.StatusNotFound,
			Temporary: false,
		})
	})

	server := &http.Server{Addr: "0.0.0.0:4000", Handler: r}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second) //nolint:lostcancel

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				l.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			l.Fatal(err)
		}
		serverStopCtx()
	}()

	l.Info("server started, listening on :4000")

	// Run the server
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		l.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	l.Info("server shutdown")
}
