package test

import (
	"net/http"
	"testing"
)

func TestHealthz(t *testing.T) {
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
			name:               "test get healthz",
			body:               "",
			expectedStatusCode: http.StatusOK,
			expectedSchema:     readJSONToString(t, "json_response_schema/get_healthz.json"),
			expectedCode:       0,
			expectedMessage:    "OK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := fastHTTPTester(t, r.Handler())

			resp := e.GET("/api/healthz").Expect()

			resp.Status(tt.expectedStatusCode)
			resp.JSON().Schema(tt.expectedSchema)
			resp.JSON().Object().Value("response_code").IsEqual(tt.expectedCode)
			resp.JSON().Object().Value("response_message").IsEqual(tt.expectedMessage)

		})
	}
}
