package entity

import "time"

type Token struct {
	Value string
	Exp   time.Duration
}
