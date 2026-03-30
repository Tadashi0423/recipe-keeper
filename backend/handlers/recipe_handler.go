package handlers

import (
	"encoding/json"
	"net/http"
	"recipe-keeper/models"
	"recipe-keeper/services"
	"strconv"

	"github.com/gorilla/mux"
)

// Response helpers
func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	JSONResponse(w, status, map[string]string{"error": message})
}

// GET /recipes
func GetAllRecipes(w http.ResponseWriter, r *http.Request) {
	recipe, err := services.GetAllRecipes()
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSONResponse(w, http.StatusOK, recipe)
}

// GET /recipes/:id
func GetRecipeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	recipe, err := services.GetRecipeByID(id)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if recipe == nil {
		ErrorResponse(w, http.StatusNotFound, "Recipe not found")
		return
	}

	JSONResponse(w, http.StatusOK, recipe)
}

// POST /recipes
func CreateRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe models.Recipe
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&recipe)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	id, err := services.CreateRecipe(recipe)
	if err != nil {
		if err == services.ErrTitleRequired {
			ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusCreated, map[string]int64{"id": id})
}

// PUT /recipes/:id
func UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var recipe models.Recipe
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&recipe)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	err = services.UpdateRecipe(id, recipe)
	if err != nil {
		if err == services.ErrRecipeNotFound {
			ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		if err == services.ErrTitleRequired {
			ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, map[string]string{"message": "Updated successfully"})
}

// DELETE /recipes/:id
func DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	err = services.DeleteRecipe(id)
	if err != nil {
		if err == services.ErrRecipeNotFound {
			ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSONResponse(w, http.StatusOK, map[string]string{"message": "Deleted successfully"})

}

// GET /recipes/search?q=keyword
func SearchRecipes(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	recipes, err := services.SearchRecipes(query)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSONResponse(w, http.StatusOK, recipes)

}

// GET /recipes/cuisine/:cuisine
func GetRecipesByCuisine(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cuisine := vars["cuisine"]

	recipes, err := services.GetRecipesByCuisine(cuisine)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, recipes)
}
