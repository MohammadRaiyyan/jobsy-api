package utils

type Status string

const (
	Success Status = "success"
	Error   Status = "error"
)

type Response struct {
	Status  Status      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
