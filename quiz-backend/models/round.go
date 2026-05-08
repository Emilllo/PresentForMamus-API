package models

type Round struct {
	ID     int  `json:"id"`
	GameId *int `json:"game_id"`
}