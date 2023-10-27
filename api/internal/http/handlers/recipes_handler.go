package handlers

import (
	"net/http"

	"github.com/cdimonaco/ephimeral-pr-env-demo/api/internal/persistence"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type RecipesHandler struct {
	logger     *zap.SugaredLogger
	repository *persistence.RecipesRepository
}

func NewRecipesHandler(
	logger *zap.SugaredLogger,
	repository *persistence.RecipesRepository,
) *RecipesHandler {
	l := logger.With("component", "recipesHandler")

	return &RecipesHandler{logger: l, repository: repository}
}

func (h *RecipesHandler) GetAllRecipes(w http.ResponseWriter, r *http.Request) {
	recipes, err := h.repository.GetRecipes(r.Context())
	if err != nil {
		h.logger.Errorw(
			"error during getAllRecipes execution",
			"error",
			err,
			"reqId",
			middleware.GetReqID(r.Context()),
		)

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, HttpError{Error: "internal error, try later", Code: http.StatusInternalServerError})
		return
	}

	render.JSON(w, r, recipes)
}
