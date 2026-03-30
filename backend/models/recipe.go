package models

import "time"

type Recipe struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	URL          string    `json:"url"`
	Cuisine      string    `json:"cuisine"`
	Ingredients  []string  `json:"ingredients"`
	Instructions string    `json:"instructions"`
	CookingTime  int       `json:"cooking_time"`
	Servings     int       `json:"servings"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
