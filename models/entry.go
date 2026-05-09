package models

import "time"

type Entry struct {
	Site      string    `json:"site"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	TOTP      string    `json:"totp,omitempty"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Vault struct {
	Version    int    `json:"version"`
	Salt       string `json:"salt"`
	CipherText []byte `json:"cipher_text"`
}
