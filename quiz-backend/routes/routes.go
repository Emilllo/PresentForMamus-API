package routes

import (
	"net/http"

	"quiz-backend/handlers"
)

func SetupRoutes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/users", handlers.UsersHandler)
	router.HandleFunc("/categories", handlers.CategoriesHandler)

	return router
}