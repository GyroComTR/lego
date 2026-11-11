package internal

import "fmt"

// APIError is the error returned by the API.
type APIError struct {
	ErrorInfo ErrorInfo `json:"error"`
}

func (a *APIError) Error() string {
	return fmt.Sprintf("%s: %s", a.ErrorInfo.Code, a.ErrorInfo.Message)
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ChallengeRequest struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ChallengeResponse struct {
	Data *Challenge `json:"data"`
}

type Challenge struct {
	ServiceID int    `json:"service_id"`
	Zone      string `json:"zone"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Values    int    `json:"values"`
	Changed   bool   `json:"changed"`
}
