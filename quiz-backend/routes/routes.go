package routes

import (
	"net/http"

	"quiz-backend/handlers"
)

func SetupRoutes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/users", handlers.UsersHandler)
	router.HandleFunc("/users/{id}", handlers.UserByIDHandler)
	router.HandleFunc("/categories", handlers.CategoriesHandler)
	router.HandleFunc("/categories/{id}", handlers.CategoryByIDHandler)
	router.HandleFunc("/questions", handlers.QuestionsHandler)
	router.HandleFunc("/questions/{id}", handlers.QuestionByIDHandler)

	return router
}
