package entity

type Session struct {
	Subject     string `db:"subject"`
	Fingerprint string `db:"fingerprint"`
	ExpiresIn   int64  `db:"expires_in"`
}
