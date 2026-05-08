package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"quiz-backend/db"
	"quiz-backend/models"
)

func QuestionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetQuestions(w, r)

	case http.MethodPost:
		CreateQuestion(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func QuestionByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		UpdateQuestion(w, r)
	case http.MethodDelete:
		DeleteQuestion(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func GetQuestions(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(
		context.Background(),
		`
		SELECT id, question, points_per_question, image, answer, category_id
		FROM questions
		ORDER BY id ASC
	`,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	questions := []models.Questions{}

	for rows.Next() {
		var question models.Questions

		err := rows.Scan(
			&question.ID,
			&question.Question,
			&question.PointsPerQuestion,
			&question.ImageURL,
			&question.Answer,
			&question.CategoryID,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		questions = append(questions, question)
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(questions)
}

func CreateQuestion(w http.ResponseWriter, r *http.Request) {
	var question models.Questions

	err := json.NewDecoder(r.Body).Decode(&question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.DB.QueryRow(
		context.Background(),
		`
		INSERT INTO questions (question, points_per_question, image, answer, category_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
		`,
		question.Question,
		question.PointsPerQuestion,
		question.ImageURL,
		question.Answer,
		question.CategoryID,
	).Scan(&question.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(question)
}

func UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var question models.Questions
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	commandTag, err := db.DB.Exec(
		context.Background(),
		`
		UPDATE questions
		SET question = $1, points_per_question = $2, image = $3, answer = $4, category_id = $5
		WHERE id = $6
		`,
		question.Question,
		question.PointsPerQuestion,
		question.ImageURL,
		question.Answer,
		question.CategoryID,
		id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Error(w, "Question not found", http.StatusNotFound)
		return
	}

	question.ID = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(question)
}

func DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	commandTag, err := db.DB.Exec(
		context.Background(),
		`DELETE FROM questions WHERE id = $1`,
		id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Error(w, "Question not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
