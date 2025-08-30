package domain

type Message struct {
	Key     string            `json:"key"`
	Headers map[string]string `json:"headers"`
	Content []byte            `json:"content"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
