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
		`
		SELECT g.id AS gameid, r.id AS roundid, cat.id AS categoryID, cat.name AS categoryName, cat.description AS catDescription
		FROM game g
		JOIN round r ON r.game_id = g.id
		JOIN categories cat ON cat.round_id = r.id
		ORDER BY g.id ASC, r.id ASC, cat.id ASC
		`,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	games := []models.GameDetails{}
	for rows.Next() {
		var game models.GameDetails

		if err := rows.Scan(
			&game.GameID,
			&game.RoundID,
			&game.CategoryID,
			&game.CategoryName,
			&game.CatDescription,
		); err != nil {
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

	rows, err := db.DB.Query(
		context.Background(),
		`
		SELECT g.id AS gameid, r.id AS roundid, cat.id AS categoryID, cat.name AS categoryName, cat.description AS catDescription
		FROM game g
		JOIN round r ON r.game_id = g.id
		JOIN categories cat ON cat.round_id = r.id
		WHERE g.id = $1
		ORDER BY r.id ASC, cat.id ASC
		`,
		id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	games := []models.GameDetails{}
	for rows.Next() {
		var game models.GameDetails

		if err := rows.Scan(
			&game.GameID,
			&game.RoundID,
			&game.CategoryID,
			&game.CategoryName,
			&game.CatDescription,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		games = append(games, game)
	}

	if len(games) == 0 {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
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

	tx, err := db.DB.Begin(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(
		context.Background(),
		`UPDATE categories SET round_id = NULL WHERE round_id IN (SELECT id FROM round WHERE game_id = $1)`,
		id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(
		context.Background(),
		`UPDATE round SET game_id = NULL WHERE game_id = $1`,
		id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	commandTag, err := tx.Exec(
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

	if err := tx.Commit(context.Background()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
