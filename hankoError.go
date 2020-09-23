package hankoApiClient

import "fmt"

type HankoError struct {
	Code    int  `json:"code"`
	Message string `json:"message"`
}

func (e *HankoError) Error() string {
	return fmt.Sprintf("Code: %d Message: %s", e.Code, e.Message)
}
