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
	router.HandleFunc("/categories/{id}/round", handlers.CategoryRoundHandler)
	router.HandleFunc("/questions", handlers.QuestionsHandler)
	router.HandleFunc("/questions/{id}", handlers.QuestionByIDHandler)
	router.HandleFunc("/games", handlers.GameHandler)
	router.HandleFunc("/games/{id}", handlers.GameByIDHandler)
	router.HandleFunc("/rounds", handlers.RoundHandler)
	router.HandleFunc("/rounds/{id}", handlers.RoundByIDHandler)

	return router
}
