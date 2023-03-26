package response

// Standard Response eg. 200 OK or 404 Not Found
type Response struct {
	Code    uint64 `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// 422 Unprocessable Entity
type UnprocessableEntityError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Errors  any    `json:"errors"`
}
