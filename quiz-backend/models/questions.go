package models

type Questions struct {
	ID                int     `json:"id"`
	Question          string  `json:"question"`
	PointsPerQuestion int     `json:"points_per_question"`
	ImageURL          *string `json:"image,omitempty"`
	Answer            *string `json:"answer"`
	CategoryID        int     `json:"category_id"`
}