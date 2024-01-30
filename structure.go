package golim

import (
	"time"
)

type Limiter struct {
	Duration time.Duration //Time of each session. Default: 1 second.
	Max int //Max requests of each session. Default 15 requests per a session.
}

func NewLimiter() *Limiter {
	return &Limiter{
		Duration: time.Second,
		Max: 15,
	}
}