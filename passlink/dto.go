package passlink

import (
	"github.com/google/uuid"
	"time"
)

type LinkRequest struct {
	UserID     string `json:"user_id"`
	Transport  string `json:"transport"`
	Phone      string `json:"phone_no"`
	Email      string `json:"email"`
	Template   string `json:"template"`
	Locale     string `json:"locale"`
	Timezone   string `json:"timezone"`
	Salutation string `json:"salutation"`
	TTL        string `json:"ttl"`
	RedirectTo string `json:"redirect_to"`
}

type Link struct {
	ID         uuid.UUID `json:"id"`
	UserID     string    `json:"user_id"`
	Status     string    `json:"status"`
	ValidUntil time.Time `json:"valid_until"`
}
