package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"quiz-backend/db"
	"quiz-backend/models"

	"github.com/google/uuid"
)

func CreatePlayer(w http.ResponseWriter, r *http.Request) {
	var player models.Players
	token := uuid.New().String()

	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	err = db.DB.QueryRow(
		context.Background(),	
		`
		INSERT INTO players (name, scores, token, is_blocked, created_at)
		VALUES ($1, $2, $3, false, now())
		RETURNING id, created_at
	`,
		player.Name,
		player.Scores,
		token,

	).Scan(&player.ID, &player.CreatedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(player)
}