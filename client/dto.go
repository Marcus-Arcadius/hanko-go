package client

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
	return e.Message
}

func WrapError(err error) *ApiError {
	return &ApiError{Message: err.Error()}
}
