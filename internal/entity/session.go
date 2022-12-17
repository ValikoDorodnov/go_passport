package entity

type Session struct {
	Subject     string `redis:"subject"`
	Fingerprint string `redis:"fingerprint"`
}
