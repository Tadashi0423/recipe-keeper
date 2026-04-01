package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"recipe-keeper/database"
	"recipe-keeper/handlers"

	"github.com/gorilla/mux"
)

// CORS middleware
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Connect to database
	dbPath := "recipe.db"
	if err := database.Connect(dbPath); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to run migration: %v", err)
	}

	// Create router
	router := mux.NewRouter()

	// IMPORTANT: Register specific routes BEFORE /{id} routes
	// Order matters in gorilla/mux - first match wins
	
	// GET routes
	router.HandleFunc("/recipes", handlers.GetAllRecipes).Methods(http.MethodGet)
	router.HandleFunc("/recipes/search", handlers.SearchRecipes).Methods(http.MethodGet)
	router.HandleFunc("/recipes/cuisine/{cuisine}", handlers.GetRecipesByCuisine).Methods(http.MethodGet)
	router.HandleFunc("/recipes/{id}", handlers.GetRecipeByID).Methods(http.MethodGet)

	// POST route
	router.HandleFunc("/recipes", handlers.CreateRecipe).Methods(http.MethodPost)

	// PUT route
	router.HandleFunc("/recipes/{id}", handlers.UpdateRecipe).Methods(http.MethodPut)

	// DELETE route
	router.HandleFunc("/recipes/{id}", handlers.DeleteRecipe).Methods(http.MethodDelete)

	// Start server with CORS wrapped around router
	port := ":8080"
	log.Printf("Server starting on http://localhost%s", port)

	// Graceful shutdown
	go func() {
		if err := http.ListenAndServe(port, withCORS(router)); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
