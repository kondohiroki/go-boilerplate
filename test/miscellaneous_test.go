package test

import (
	"net/http"
	"testing"
)

func TestMiscellaneous(t *testing.T) {
	type params struct{}

	tests := []struct {
		name               string
		params             params
		body               any
		expectedStatusCode int
		expectedSchema     string
		expectedCode       int
		expectedMessage    string
	}{
		{
			name:               "test not found",
			body:               "",
			expectedStatusCode: http.StatusNotFound,
			expectedSchema:     readJSONToString(t, "json_response_schema/misc_not_found.json"),
			expectedCode:       404,
			expectedMessage:    "this is not api you are looking for",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := fastHTTPTester(t, r.Handler())

			resp := e.GET("/api/v1/notfound").Expect()

			resp.Status(tt.expectedStatusCode)
			resp.JSON().Schema(tt.expectedSchema)
			resp.JSON().Object().Value("response_code").IsEqual(tt.expectedCode)
			resp.JSON().Object().Value("response_message").IsEqual(tt.expectedMessage)

		})
	}
}
