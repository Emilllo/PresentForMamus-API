package models

type Categories struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	RoundID     *int    `json:"round_id,omitempty"`
}

type CategoriesByRound struct {
	CategoryId          int     `json:"cat_id"`
	CategoryName        string  `json:"cat_name"`
	CategoryDescription *string `json:"cat_desc,omitempty"` // * - необязательный. omitempty - не включать в JSON, если значение пустое
	Round               int     `json:"round_id"`
	Game                int     `json:"game_id"`
}
