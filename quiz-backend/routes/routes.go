package routes

import (
	"net/http"

	"quiz-backend/handlers"
)

func SetupRoutes() *http.ServeMux {
	router := http.NewServeMux()

	//players
	router.HandleFunc("POST /create-player", handlers.CreatePlayer)
	//categories
	router.HandleFunc("/categories", handlers.CategoriesHandler)
	router.HandleFunc("/categories/{id}", handlers.CategoryByIDHandler)
	router.HandleFunc("/categories/{id}/round", handlers.CategoryRoundHandler)
	//questions
	router.HandleFunc("/questions", handlers.QuestionsHandler)
	router.HandleFunc("/questions/{id}", handlers.QuestionByIDHandler)
	//games
	router.HandleFunc("/games", handlers.GameHandler)
	router.HandleFunc("/games/{id}", handlers.GameByIDHandler)
	//rounds
	router.HandleFunc("/rounds", handlers.RoundHandler)
	router.HandleFunc("/rounds/{id}", handlers.RoundByIDHandler)


	router.HandleFunc("/games/{game_id}/rounds/{round_id}/categories", handlers.CategoriesByRoundHandler)

	return router
}
