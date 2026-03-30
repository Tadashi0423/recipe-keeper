package database

import (
	"database/sql"
	"encoding/json"
	"recipe-keeper/models"
	"time"
)

func CreateRecipe(recipe models.Recipe) (int64, error) {
	ingredientsJSON, err := json.Marshal(recipe.Ingredients)
	if err != nil {
		return 0, err
	}

	query := `
	INSERT INTO recipes (title, url, cuisine, ingredients, instructions, cooking_time, servings, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := DB.Exec(query,
		recipe.Title,
		recipe.URL,
		recipe.Cuisine,
		string(ingredientsJSON),
		recipe.Instructions,
		recipe.CookingTime,
		recipe.Servings,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetAllRecipes() ([]models.Recipe, error) {
	query := `SELECT id, title, url, cuisine, ingredients, instructions, cooking_time, servings, created_at, updated_at FROM recipes ORDER BY created_at DESC`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []models.Recipe
	for rows.Next() {
		var r models.Recipe
		var ingredientsJSON string

		err := rows.Scan(&r.ID, &r.Title, &r.URL, &r.Cuisine, &ingredientsJSON, &r.Instructions, &r.CookingTime, &r.Servings, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(ingredientsJSON), &r.Ingredients)
		recipes = append(recipes, r)
	}

	return recipes, nil
}

func GetRecipeByID(id int64) (*models.Recipe, error) {
	query := `SELECT id, title, url, cuisine, ingredients, instructions, cooking_time, servings, created_at, updated_at FROM recipes WHERE id = ?`

	var r models.Recipe
	var ingredientsJSON string

	err := DB.QueryRow(query, id).Scan(&r.ID, &r.Title, &r.URL, &r.Cuisine, &ingredientsJSON, &r.Instructions, &r.CookingTime, &r.Servings, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	json.Unmarshal([]byte(ingredientsJSON), &r.Ingredients)
	return &r, nil
}

func UpdateRecipe(id int64, recipe models.Recipe) error {
	ingredientsJSON, err := json.Marshal(recipe.Ingredients)
	if err != nil {
		return err
	}

	query := `UPDATE recipes SET title = ?, url = ?, cuisine = ?, ingredients = ?, instructions = ?, cooking_time = ?, servings = ?, updated_at = ? WHERE id = ?`

	_, err = DB.Exec(query,
		recipe.Title,
		recipe.URL,
		recipe.Cuisine,
		string(ingredientsJSON),
		recipe.Instructions,
		recipe.CookingTime,
		recipe.Servings,
		time.Now(),
		id,
	)
	return err
}

func DeleteRecipe(id int64) error {
	query := `DELETE FROM recipes WHERE id = ?`
	_, err := DB.Exec(query, id)
	return err
}

func SearchRecipes(keyword string) ([]models.Recipe, error) {
	query := `SELECT id, title, url, cuisine, ingredients, instructions, cooking_time, servings, created_at, updated_at FROM recipes WHERE title LIKE ? OR cuisine LIKE ? OR ingredients LIKE ? ORDER BY created_at DESC`

	searchPattern := "%" + keyword + "%"

	rows, err := DB.Query(query, searchPattern, searchPattern, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []models.Recipe
	for rows.Next() {
		var r models.Recipe
		var ingredientsJSON string

		err := rows.Scan(&r.ID, &r.Title, &r.URL, &r.Cuisine, &ingredientsJSON, &r.Instructions, &r.CookingTime, &r.Servings, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(ingredientsJSON), &r.Ingredients)
		recipes = append(recipes, r)
	}

	return recipes, nil
}
