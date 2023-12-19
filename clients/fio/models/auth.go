package models

import "time"

type AuthRequest struct {
	Username string `json:"UserName"`
	Password string `json:"Password"`
}

type AuthResponse struct {
	AuthToken       string    `json:"AuthToken"`
	Expiry          time.Time `json:"Expiry"`
	IsAdministrator bool      `json:"IsAdministrator"`
}
