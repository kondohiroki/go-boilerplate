package exception

import "github.com/bytedance/sonic"

// ExceptionErrors is used as our project error response.
// All error response will be in this format.
type ExceptionErrors struct {
	HttpStatusCode int

	GlobalMessage string
	ErrItems      []*ExceptionError
}

// Error implements go built-in error interface.
// This will output to CommonResponse for our project.
func (cErrs *ExceptionErrors) Error() string {
	return cErrs.GlobalMessage
}

// MarshalJSON implements JSON marshaller interface.
// This will marshal only property ErrItems.
func (cErrs *ExceptionErrors) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(cErrs.ErrItems)
}

type ExceptionError struct {
	Message      string       `json:"message"`
	Type         errorType    `json:"type"`
	ErrorSubcode errorSubcode `json:"error_subcode"`
}

func (cErr *ExceptionError) Error() string {
	return cErr.Message
}

// NewExceptionErrors allocates new empty error item ExceptionErrors
func NewExceptionErrors(httpStatusCode int, globalMessage string) *ExceptionErrors {
	return &ExceptionErrors{
		GlobalMessage:  globalMessage,
		HttpStatusCode: httpStatusCode,
	}
}

func createFixedExceptionErrors(
	httpStatusCode int,
	t errorType,
	esc errorSubcode,
	m string,
) *ExceptionErrors {
	return &ExceptionErrors{
		GlobalMessage:  m,
		HttpStatusCode: httpStatusCode,
		ErrItems: []*ExceptionError{
			{
				Message:      m,
				Type:         t,
				ErrorSubcode: esc,
			},
		},
	}
}
