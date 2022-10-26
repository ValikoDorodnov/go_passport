package entity

import "time"

type ParsedToken struct {
	Subject string
	Jwt     string
	ExpTtl  time.Duration
}
