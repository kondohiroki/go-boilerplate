package response

// Standard Response
type CommonResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Errors    any    `json:"errors,omitempty"`
	Data      any    `json:"data,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}
