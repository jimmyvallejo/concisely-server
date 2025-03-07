package handlers

import "golang.org/x/time/rate"

type Handlers struct {
	GPTKEY  string
	Limiter *rate.Limiter
}

func NewHandlers(key string) *Handlers {
	return &Handlers{
		GPTKEY:  key,
		Limiter: rate.NewLimiter(rate.Limit(10), 20),
	}
}
