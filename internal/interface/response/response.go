package response

import (
	"github.com/bytedance/sonic"
	"github.com/kondohiroki/go-boilerplate/pkg/exception"
)

type DataUnwrapper interface {
	UnwrapData(interface{}) error
}

// Standard Response
type CommonResponse struct {
	ResponseCode    int                        `json:"response_code"`
	ResponseMessage string                     `json:"response_message"`
	Errors          *exception.ExceptionErrors `json:"errors,omitempty"`
	Data            any                        `json:"data,omitempty"`
	RequestID       string                     `json:"request_id,omitempty"`
}

func (resp *CommonResponse) UnwrapData(target interface{}) error {
	bs, err := sonic.Marshal(resp.Data)
	if err != nil {
		return err
	}

	if err := sonic.Unmarshal(bs, target); err != nil {
		return err
	}

	return nil
}
