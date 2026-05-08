package main

import (
	"log"
	"net/http"

	"quiz-backend/db"
	"quiz-backend/routes"
)

func main() {
	db.Connect()

	router := routes.SetupRoutes()

	log.Println("Server started on :3000")

	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatal(err)
	}
}