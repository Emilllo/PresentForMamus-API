package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"quiz-backend/db"
	"quiz-backend/models"
)

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetCategories(w, r)

	case http.MethodPost:
		CreateCategory(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func CategoryByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		UpdateCategory(w, r)
	case http.MethodDelete:
		DeleteCategory(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func CategoryRoundHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		SetCategoryRound(w, r)
	case http.MethodDelete:
		ClearCategoryRound(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(
		context.Background(),
		`
		SELECT id, name, description, round_id
		FROM categories
		ORDER BY name ASC
	`,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	categories := []models.Categories{}

	for rows.Next() {
		var category models.Categories

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.RoundID,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		categories = append(categories, category)
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(categories)
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Categories

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.DB.QueryRow(
		context.Background(),
		`
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
		RETURNING id
		`,
		category.Name,
		category.Description,
	).Scan(&category.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(category)
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var category models.Categories
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	commandTag, err := db.DB.Exec(
		context.Background(),
		`
		UPDATE categories
		SET name = $1, description = $2
		WHERE id = $3
		`,
		category.Name,
		category.Description,
		id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	category.ID = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	commandTag, err := db.DB.Exec(
		context.Background(),
		`DELETE FROM categories WHERE id = $1`,
		id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func SetCategoryRound(w http.ResponseWriter, r *http.Request) {
	categoryID, ok := parseID(w, r)
	if !ok {
		return
	}

	var payload struct {
		RoundID *int `json:"round_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if payload.RoundID == nil || *payload.RoundID <= 0 {
		http.Error(w, "Invalid round_id", http.StatusBadRequest)
		return
	}

	var roundID int
	err := db.DB.QueryRow(
		context.Background(),
		`SELECT id FROM round WHERE id = $1`,
		*payload.RoundID,
	).Scan(&roundID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Round not found", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	commandTag, err := db.DB.Exec(
		context.Background(),
		`UPDATE categories SET round_id = $1 WHERE id = $2`,
		*payload.RoundID,
		categoryID,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			http.Error(w, "Round not found", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		ID      int `json:"id"`
		RoundID int `json:"round_id"`
	}{
		ID:      categoryID,
		RoundID: *payload.RoundID,
	})
}

func ClearCategoryRound(w http.ResponseWriter, r *http.Request) {
	categoryID, ok := parseID(w, r)
	if !ok {
		return
	}

	commandTag, err := db.DB.Exec(
		context.Background(),
		`UPDATE categories SET round_id = NULL WHERE id = $1`,
		categoryID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
