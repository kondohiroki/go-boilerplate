package response_test

import (
	"encoding/json"
	"testing"

	. "github.com/kondohiroki/go-boilerplate/internal/interface/response"
	"github.com/stretchr/testify/assert"
)

func TestCommonResponse_JSONMarshalling(t *testing.T) {
	cr := CommonResponse{
		ResponseCode:    200,
		ResponseMessage: "Success",
		Errors:          nil,
		Data:            "Sample Data",
		RequestID:       "12345",
	}

	expectedJSON := `{"response_code":200,"response_message":"Success","data":"Sample Data","request_id":"12345"}`

	marshalledJSON, err := json.Marshal(cr)
	if err != nil {
		t.Fatalf("Failed to marshal CommonResponse: %v", err)
	}

	if string(marshalledJSON) != expectedJSON {
		t.Errorf("Expected JSON: %s, got: %s", expectedJSON, string(marshalledJSON))
	}
}

func TestCommonResponse_JSONUnmarshalling(t *testing.T) {
	inputJSON := `{"response_code":200,"response_message":"Success","data":"Sample Data","request_id":"12345"}`

	expectedResponse := CommonResponse{
		ResponseCode:    200,
		ResponseMessage: "Success",
		Errors:          nil,
		Data:            "Sample Data",
		RequestID:       "12345",
	}

	var cr CommonResponse
	err := json.Unmarshal([]byte(inputJSON), &cr)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON to CommonResponse: %v", err)
	}

	if cr != expectedResponse {
		t.Errorf("Expected CommonResponse: %+v, got: %+v", expectedResponse, cr)
	}
}

func TestCommonResponse_EmptyOptionalFields(t *testing.T) {
	cr := CommonResponse{
		ResponseCode:    200,
		ResponseMessage: "Success",
	}

	expectedJSON := `{"response_code":200,"response_message":"Success"}`

	marshalledJSON, err := json.Marshal(cr)
	if err != nil {
		t.Fatalf("Failed to marshal CommonResponse: %v", err)
	}

	if string(marshalledJSON) != expectedJSON {
		t.Errorf("Expected JSON: %s, got: %s", expectedJSON, string(marshalledJSON))
	}
}

type Student struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (t *Student) Object() interface{} {
	return &Student{}
}

type Error struct {
	Msg string `json:"error_message"`
}

func (t *Error) Object() interface{} {
	return &Error{}
}

type Objecter interface {
	Object() interface{}
}

func TestUnwrapper(t *testing.T) {

	type testCase struct {
		name     string
		input    DataUnwrapper
		expected Objecter
		isError  bool
	}

	testCases := []testCase{
		{
			name: "unwrap common response success 1",
			input: &CommonResponse{
				ResponseCode:    0,
				ResponseMessage: "success",
				Data: &Student{
					Name: "John",
					Age:  27,
				},
			},
			expected: &Student{
				Name: "John",
				Age:  27,
			},
		},
		{
			name: "unwrap common response success 2",
			input: &CommonResponse{
				ResponseCode:    1,
				ResponseMessage: "internal server error",
				Data: &Error{
					Msg: "test error",
				},
			},
			expected: &Error{
				Msg: "test error",
			},
		},
		{
			name: "unwrap common response marshal error",
			input: &CommonResponse{
				ResponseCode:    1,
				ResponseMessage: "internal server error",
				Data:            func() {},
			},
			expected: &Error{},
			isError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.expected.Object()
			err := tc.input.UnwrapData(actual)
			if tc.isError {
				assert.NotNil(t, err)
				assert.Empty(t, actual)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}
