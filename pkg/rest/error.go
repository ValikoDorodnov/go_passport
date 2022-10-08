package rest

type ErrorResponse struct {
	Errors []string `json:"errors"`
}
