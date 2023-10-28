package persistence

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const recipesTable = "recipes"

type RecipesRepository struct {
	conn    *pgxpool.Pool
	logger  *zap.SugaredLogger
	builder squirrel.StatementBuilderType
}

func NewRecipesRepository(
	conn *pgxpool.Pool,
	logger *zap.SugaredLogger,
) *RecipesRepository {
	l := logger.With("component", "recipesRepository")

	return &RecipesRepository{
		conn:    conn,
		logger:  l,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *RecipesRepository) GetRecipes(ctx context.Context) ([]Recipe, error) {
	recipes := []Recipe{}

	stmt, _, err := r.builder.
		Select(
			"id",
			"name",
			"description",
			"ingredients",
		).
		From(recipesTable).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("error during get all recipes: %w", err)
	}

	err = pgxscan.Select(ctx, r.conn, &recipes, stmt)
	if err != nil {
		return nil, fmt.Errorf("error during get all recipes: %w", err)
	}

	return recipes, err
}
