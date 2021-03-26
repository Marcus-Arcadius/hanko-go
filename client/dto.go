package client

import (
	"fmt"
)

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	IconUrl     string `json:"icon"`
}

type ApiError struct {
	Message      string `json:"message"`
	Details      string `json:"details"`
	DebugMessage string `json:"debug_message"`
	StatusText   string `json:"status_text"`
	StatusCode   int    `json:"status_code"`
}

func (e *ApiError) Error() string {
	str := e.Message
	if e.Details != "" {
		str = fmt.Sprintf("%s: %s", str, e.Details)
	}
	if e.DebugMessage != "" {
		str = fmt.Sprintf("%s: %s", str, e.DebugMessage)
	}
	return str
}

func WrapError(err error) *ApiError {
	return &ApiError{
		Message:      "sdk error",
		Details:      "an error occurred while processing the request",
		DebugMessage: err.Error(),
		StatusText:   "Internal Server Error",
		StatusCode:   500,
	}
}
