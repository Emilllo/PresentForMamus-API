package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"quiz-backend/db"
	"quiz-backend/models"
)

func GameHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		CreateGame(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func GameByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		DeleteGame(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func CreateGame(w http.ResponseWriter, r *http.Request) {
	var game models.Game

	err := db.DB.QueryRow(
		context.Background(),
		`
		INSERT INTO games DEFAULT VALUES
		RETURNING id
		`,
	).Scan(&game.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(game)
}

func DeleteGame(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	commandTag, err := db.DB.Exec(
		context.Background(),
		`DELETE FROM games WHERE id = $1`,
		id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
