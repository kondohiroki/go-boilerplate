package response_test

import (
	"encoding/json"
	"testing"

	. "github.com/kondohiroki/go-boilerplate/internal/interface/response"
)

func TestCommonResponse_JSONMarshalling(t *testing.T) {
	cr := CommonResponse{
		Code:      200,
		Message:   "Success",
		Errors:    nil,
		Data:      "Sample Data",
		RequestID: "12345",
	}

	expectedJSON := `{"code":200,"message":"Success","data":"Sample Data","request_id":"12345"}`

	marshalledJSON, err := json.Marshal(cr)
	if err != nil {
		t.Fatalf("Failed to marshal CommonResponse: %v", err)
	}

	if string(marshalledJSON) != expectedJSON {
		t.Errorf("Expected JSON: %s, got: %s", expectedJSON, string(marshalledJSON))
	}
}

func TestCommonResponse_JSONUnmarshalling(t *testing.T) {
	inputJSON := `{"code":200,"message":"Success","data":"Sample Data","request_id":"12345"}`

	expectedResponse := CommonResponse{
		Code:      200,
		Message:   "Success",
		Errors:    nil,
		Data:      "Sample Data",
		RequestID: "12345",
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
		Code:    200,
		Message: "Success",
	}

	expectedJSON := `{"code":200,"message":"Success"}`

	marshalledJSON, err := json.Marshal(cr)
	if err != nil {
		t.Fatalf("Failed to marshal CommonResponse: %v", err)
	}

	if string(marshalledJSON) != expectedJSON {
		t.Errorf("Expected JSON: %s, got: %s", expectedJSON, string(marshalledJSON))
	}
}
