package responses

import "github.com/cdimonaco/ephimeral-pr-env-demo/api/internal/persistence"

type RecipeResponse struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Ingredients []string `json:"recipe_ingredients"`
}

func MapRecipeEntityToRecipeResponse(entity persistence.Recipe) RecipeResponse {
	return RecipeResponse{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		Ingredients: entity.Ingredients,
	}
}

func MapRecipesEntitiesToRecipeResponseList(entities []persistence.Recipe) []RecipeResponse {
	result := []RecipeResponse{}
	for _, e := range entities {
		result = append(result, MapRecipeEntityToRecipeResponse(e))
	}

	return result
}
