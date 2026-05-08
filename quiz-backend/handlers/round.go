package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"quiz-backend/db"
	"quiz-backend/models"
)

func RoundHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		CreateRound(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func RoundByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		DeleteRound(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func CreateRound(w http.ResponseWriter, r *http.Request) {
    var round models.Round

    if err := json.NewDecoder(r.Body).Decode(&round); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if round.GameId == nil || *round.GameId <= 0 {
        http.Error(w, "Invalid game_id", http.StatusBadRequest)
        return
    }

    err := db.DB.QueryRow(
        context.Background(),
        `
        INSERT INTO round(game_id) VALUES ($1)
        RETURNING id
        `,
        round.GameId,
    ).Scan(&round.ID)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)

    json.NewEncoder(w).Encode(round)
}

func DeleteRound(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	commandTag, err := db.DB.Exec(
		context.Background(),
		`DELETE FROM round WHERE id = $1`,
		id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Error(w, "Round not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
