package request

type LoginByEmail struct {
	Email    string `json:"email"`
	Pass     string `json:"pass"`
	Platform string `json:"platform"`
}
