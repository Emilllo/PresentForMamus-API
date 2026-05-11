package models

import "time"

type Players struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Scores    int       `json:"scores"`
	Token     string    `json:"token"`
	IsBlocked bool      `json:"is_blocked"`
	CreatedAt *time.Time `json:"created_at"`
}