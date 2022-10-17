package entity

type Session struct {
	Subject   int    `db:"subject"`
	Platform  string `db:"platform"`
	ExpiresIn int64  `db:"expires_in"`
}
