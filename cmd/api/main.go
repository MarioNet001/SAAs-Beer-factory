package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"sistema-gestion-beer/src/domain/recipe"
	"sistema-gestion-beer/src/infrastructure/db/postgres"
	"sistema-gestion-beer/src/inventory"
	recipeapi "sistema-gestion-beer/src/recipe"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is required")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	invRepo := postgres.NewInventoryRepo(db)
	service := inventory.NewInventoryService(invRepo, db)
	handler := inventory.NewHandler(service)

	recipeRepo := postgres.NewRecipeRepo(db)
	recipeService := recipe.NewRecipeService(recipeRepo, invRepo)
	recipeHandler := recipeapi.NewHandler(recipeService)

	http.HandleFunc("/recipes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var req recipe.Recipe
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			if err := recipeHandler.HandleCreateRecipe(r.Context(), &req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusCreated)
		} else if r.Method == http.MethodGet {
			recipes, err := recipeHandler.HandleListRecipes(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(recipes)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/inventory/adjust", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			ProductID string  `json:"product_id"`
			Amount    float64 `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if err := handler.HandleAdjustStock(r.Context(), req.ProductID, req.Amount); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
