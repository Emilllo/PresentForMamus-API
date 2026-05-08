package models

type Game struct {
	ID int `json:"id"`
}

type GameDetails struct {
	GameID         int     `json:"gameid"`
	RoundID        int     `json:"roundid"`
	CategoryID     int     `json:"categoryID"`
	CategoryName   string  `json:"categoryName"`
	CatDescription *string `json:"catDescription,omitempty"`
}