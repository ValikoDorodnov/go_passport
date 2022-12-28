package request

type LoginByEmail struct {
	Email       string `json:"email"       validate:"required,email"`
	Pass        string `json:"pass"        validate:"required"`
	Fingerprint string `json:"fingerprint" validate:"required"`
}
