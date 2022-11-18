package dto

type Response struct {
	StatusCode int                    `json:"statusCode"`
	Status     string                 `json:"status"`
	Error      any                    `json:"error"`
	Data       map[string]interface{} `json:"data"`
}
