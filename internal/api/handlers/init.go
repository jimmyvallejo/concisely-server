package handlers

type Handlers struct{
	GPTKEY string
}

func NewHandlers(key string) *Handlers {
	return &Handlers{
		GPTKEY: key,
	}
}
