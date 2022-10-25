package entity

import "time"

type ParsedToken struct {
	Subject, Jwt string
	ExpTtl       time.Duration
}
