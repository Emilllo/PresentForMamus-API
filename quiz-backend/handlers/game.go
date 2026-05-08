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
	case http.MethodGet:
		GetGames(w, r)
	case http.MethodPost:
		CreateGame(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func GameByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetGame(w, r)
	case http.MethodDelete:
		DeleteGame(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func GetGames(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(
		context.Background(),
		`SELECT id FROM game ORDER BY id ASC`,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	games := []models.Game{}
	for rows.Next() {
		var game models.Game

		if err := rows.Scan(&game.ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		games = append(games, game)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
}

func GetGame(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var game models.Game

	err := db.DB.QueryRow(
		context.Background(),
		`SELECT id FROM game WHERE id = $1`,
		id,
	).Scan(&game.ID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			http.Error(w, "Game not found", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}

func CreateGame(w http.ResponseWriter, r *http.Request) {
	var game models.Game

	err := db.DB.QueryRow(
		context.Background(),
		`
		INSERT INTO game DEFAULT VALUES
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
		`DELETE FROM game WHERE id = $1`,
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
