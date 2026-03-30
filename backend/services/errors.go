package services

import "errors"

var (
	ErrTitleRequired  = errors.New("recipe title is required")
	ErrRecipeNotFound = errors.New("recipe not found")
)
