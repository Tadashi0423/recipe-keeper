package services

import (
	"recipe-keeper/database"
	"recipe-keeper/models"
)

func CreateRecipe(recipe models.Recipe) (int64, error) {
	// Validation
	if recipe.Title == "" {
		return 0, ErrTitleRequired
	}

	// Use database layer to create
	id, err := database.CreateRecipe(recipe)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetAllRecipes return all recipes
func GetAllRecipes() ([]models.Recipe, error) {
	return database.GetAllRecipes()
}

// GetRecipeByID returns a single recipe by ID
func GetRecipeByID(id int64) (*models.Recipe, error) {
	return database.GetRecipeByID(id)
}

// UpdateRecipe validates and updates an exsiting recipe
func UpdateRecipe(id int64, recipe models.Recipe) error {
	// Check if recipe exists
	existing, err := database.GetRecipeByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrRecipeNotFound
	}

	// Validation
	if recipe.Title == "" {
		return ErrTitleRequired
	}

	return database.UpdateRecipe(id, recipe)
}

// DeleteRecipe deletes a recipe by ID
func DeleteRecipe(id int64) error {
	existing, err := database.GetRecipeByID(id)
	if err != nil {
		return err
	}

	if existing == nil {
		return ErrRecipeNotFound
	}

	return database.DeleteRecipe(id)
}

// SearchRecipes searches recipes by keyword
func SearchRecipes(keyword string) ([]models.Recipe, error) {
	if keyword == "" {
		return GetAllRecipes()
	}
	return database.SearchRecipes(keyword)
}

// GetRecipesByCuisine returns recipes filtered by cuisine
func GetRecipesByCuisine(cuisine string) ([]models.Recipe, error) {
	return database.SearchRecipes(cuisine)
}
