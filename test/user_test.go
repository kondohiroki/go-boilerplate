package test

import (
	"net/http"
	"testing"
)

func TestGetUsers(t *testing.T) {
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
			name:               "test get all users",
			body:               "",
			expectedStatusCode: http.StatusOK,
			expectedSchema:     readJSONToString(t, "json_response_schema/get_users.json"),
			expectedCode:       0,
			expectedMessage:    "OK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := fastHTTPTester(t, r.Handler())

			resp := e.GET("/api/v1/users").Expect()

			resp.Status(tt.expectedStatusCode)
			resp.JSON().Schema(tt.expectedSchema)
			resp.JSON().Object().Value("code").IsEqual(tt.expectedCode)
			resp.JSON().Object().Value("message").IsEqual(tt.expectedMessage)

		})
	}
}

func TestGetUserByID(t *testing.T) {
	tests := []struct {
		name               string
		expectedStatusCode int
		expectedSchema     string
		expectedCode       int
		expectedMessage    string
	}{
		{
			name:               "test get user by id",
			expectedStatusCode: http.StatusOK,
			expectedSchema:     readJSONToString(t, "json_response_schema/get_user_by_id.json"),
			expectedCode:       0,
			expectedMessage:    "OK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := fastHTTPTester(t, r.Handler())

			resp := e.GET("/api/v1/users/1").Expect()

			resp.Status(tt.expectedStatusCode)
			resp.JSON().Schema(tt.expectedSchema)
			resp.JSON().Object().Value("code").IsEqual(tt.expectedCode)
			resp.JSON().Object().Value("message").IsEqual(tt.expectedMessage)
		})
	}
}
