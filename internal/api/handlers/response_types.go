package handlers

type ErrorResponse struct {
	Error string `json:"error"`
}

type StatusResponse struct {
	Status string `json:"status"`
}


type ChatGTPResponse struct {
	Provider string `json:"provider"`
}