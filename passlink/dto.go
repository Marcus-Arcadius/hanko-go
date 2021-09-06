package passlink

import (
	"github.com/google/uuid"
	"time"
)

type LinkRequest struct {
	Email        string `json:"email"`
	Template     string `json:"template"`
	TemplateLang string `json:"lang"`
	TTL          string `json:"ttl"`
	RedirectTo   string `json:"redirect_to"`
}

type Link struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Status       string    `json:"status"`
	RedirectTo   string    `json:"redirect_to"`
	ValidUntil   time.Time `json:"valid_until"`
	Template     string    `json:"template"`
	TemplateLang string    `json:"template_lang"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}