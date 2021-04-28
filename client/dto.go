package client

import (
	"fmt"
)

// User is the base representation of a user on whose behalf registration and authentication are performed with the
// Hanko Authentication API.
type User struct {
	ID          string `json:"id"`             // unique user id
	Name        string `json:"name,omitempty"` // unique user name, related to the ID
	DisplayName string `json:"displayName,omitempty"`
}

// ApiError is the representation for all errors returned before, during, or after making requests to the Hanko
// Authentication API.
//
// The API will always provide values for StatusCode and StatusCode. A DebugMessage is
// available on requests resulting in a Bad Request (status 401). A fully populated ApiError will only be
// returned when the API is running in development mode.
type ApiError struct {
	Message      string `json:"message"`       // contains a string which generally describes the error
	Details      string `json:"details"`       // optionally contains a string which adds details to the Message
	DebugMessage string `json:"debug_message"` // optionally contains a technical error message
	StatusText   string `json:"status_text"`   // contains the http status text which corresponds to the StatusCode
	StatusCode   int    `json:"status_code"`   // contains the http status code
}

// Error fulfills the go error interface and returns all error details available.
func (e *ApiError) Error() string {
	str := fmt.Sprintf("%d - %s", e.StatusCode, e.StatusText)

	if e.Message != "" {
		str = fmt.Sprintf("%s: %s", str, e.Message)
	}
	if e.Details != "" {
		str = fmt.Sprintf("%s: %s", str, e.Details)
	}
	if e.DebugMessage != "" {
		str = fmt.Sprintf("%s: %s", str, e.DebugMessage)
	}

	return str
}

// WrapError wraps a given error and returns an ApiError. The resulting ApiError has an underlying Internal Server Error
// (500) per default.
func WrapError(err error) *ApiError {
	return &ApiError{
		Message:      "sdk error",
		Details:      "an error occurred while processing the request",
		DebugMessage: err.Error(),
		StatusText:   "Internal Server Error",
		StatusCode:   500,
	}
}
