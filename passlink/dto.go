package passlink

import (
	"github.com/google/uuid"
	"time"
)

// LinkRequest is a representation of the request body to be used for initializing a Passlink authentication flow with
// the Hanko Authentication API.
type LinkRequest struct {
	// The ID of the user to initialize a Passlink for
	UserID     string `json:"user_id"`

	// Determines the communication channel through which Passlinks are delivered to the user.
	// Currently, only `email` is supported.
	Transport  string `json:"transport"`

	// The recipient address the message containing the Passlink should be sent to
	Email      string `json:"email"`

	// The name of the template to use for the message containing the Passlink sent to the user
	Template   string `json:"template"`

	// Locale string of the form `Language_Territory` where `Language` is a two-letter ISO 639-1 language code,
	// and `Territory` is a two-letter ISO 3166-1 alpha-2 territory code.
	// The `Language` part determines the requested language of the template that is used to construct the message
	// sent to the user. If no template is found for the given language the API will return an
	// error response.
	// The `Territory` part is used to appropriately format datetime strings in messages sent to the user. For these
	// purposes the `locale` value should be one of the values listed above. If an unknown value is used
	// or the `locale` attribute is omitted, it will default to `en_GB`.
	Locale     string `json:"locale"`

	// Timezone name of the form `Area/Location` as defined in the IANA Time Zone Database.
	// If provided, it is used to appropriately format datetime strings (e.g. the expiration date
	// of a Passlink) in messages sent to the user. Defaults to `UTC` if the given timezone name
	// cannot be resolved.
	Timezone   string `json:"timezone"`

	// The salutation to use in the Passlink message. If omitted, a default salutation as configured for the
	// given template (see `template` attribute) and language (see `locale` attribute) will be used.
	Salutation string `json:"salutation"`

	// Determines how long the Passlink should be valid. Must be a duration string consisting of a sequence of decimal
	// numbers and a unit suffix. Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h". If omitted, defaults
	// to 15m.
	TTL        string `json:"ttl"`

	// The relying party URL to redirect to after a user has confirmed (clicked) a Passlink. Must be
	// a URL that has been configured by the relying party as a valid redirect URL in the Hanko Console.
	RedirectTo string `json:"redirect_to"`
}

// NewEmailLinkRequest constructs a new LinkRequest for the given userId and recipient email address. The transport
// is automatically set to "email".
func NewEmailLinkRequest(userId string, recipient string) LinkRequest {
	return LinkRequest{
		UserID:    userId,
		Transport: "email",
		Email:     recipient,
	}
}

// WithTemplate allows you to specify the name of the LinkRequest.Template to use for the message sent to the user.
func (r LinkRequest) WithTemplate(name string) LinkRequest {
	r.Template = name
	return r
}

// WithLocale allows you to specify the LinkRequest.Locale (i.e. template language and date formatting) for the message
// sent to the user.
func (r LinkRequest) WithLocale(locale string) LinkRequest {
	r.Locale = locale
	return r
}

// WithTimezone allows you to specify the LinkRequest.Timezone for formatting datetime values in the message sent to the
// user.
func (r LinkRequest) WithTimezone(timezone string) LinkRequest {
	r.Timezone = timezone
	return r
}

// WithSalutation allows you to set a custom LinkRequest.Salutation to use in the message sent to the user
func (r LinkRequest) WithSalutation(salutation string) LinkRequest {
	r.Salutation = salutation
	return r
}

// WithTTL allows you to set a custom LinkRequest.TTL for a Passlink.
func (r LinkRequest) WithTTL(ttl string) LinkRequest {
	r.TTL = ttl
	return r
}

// WithRedirectTo allows you to set a custom LinkRequest.RedirectTo URL for a Passlink.
func (r LinkRequest) WithRedirectTo(redirectTo string) LinkRequest {
	r.RedirectTo = redirectTo
	return r
}

// Link is a representation of a Passlink.
type Link struct {
	ID         uuid.UUID `json:"id"`
	UserID     string    `json:"user_id"`
	Status     string    `json:"status"`
	ValidUntil time.Time `json:"valid_until"`
}
