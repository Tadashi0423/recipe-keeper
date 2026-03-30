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

	// Register routes
	router.HandleFunc("/recipes", handlers.GetAllRecipes).Methods(http.MethodGet)
	router.HandleFunc("/recipes/{id}", handlers.GetRecipeByID).Methods(http.MethodGet)
	router.HandleFunc("/recipes", handlers.CreateRecipe).Methods(http.MethodPost)
	router.HandleFunc("/recipes/{id}", handlers.UpdateRecipe).Methods(http.MethodPut)
	router.HandleFunc("/recipes/{id}", handlers.DeleteRecipe).Methods(http.MethodDelete)
	router.HandleFunc("/recipes/search", handlers.SearchRecipes).Methods(http.MethodGet)
	router.HandleFunc("/recipes/cuisine/{cuisine}", handlers.GetRecipesByCuisine).Methods(http.MethodGet)

	// Start server
	port := ":8080"
	log.Printf("Server starting on http://localhost%s", port)

	// Graceful shutdown
	go func() {
		if err := http.ListenAndServe(port, router); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

}
