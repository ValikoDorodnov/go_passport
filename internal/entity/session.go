package entity

type Session struct {
	Subject   string `db:"subject"`
	Platform  string `db:"platform"`
	ExpiresIn int64  `db:"expires_in"`
}
