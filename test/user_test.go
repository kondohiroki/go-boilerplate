package test

import (
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
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
			resp.JSON().Object().Value("response_code").IsEqual(tt.expectedCode)
			resp.JSON().Object().Value("response_message").IsEqual(tt.expectedMessage)

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
			resp.JSON().Object().Value("response_code").IsEqual(tt.expectedCode)
			resp.JSON().Object().Value("response_message").IsEqual(tt.expectedMessage)
		})
	}
}

func TestCreateUser(t *testing.T) {
	newAccount := map[string]interface{}{
		"name":  gofakeit.Name(),
		"email": gofakeit.Email(),
	}
	tests := []struct {
		name               string
		body               any
		expectedStatusCode int
		expectedSchema     string
		expectedCode       int
		expectedMessage    string
	}{
		{
			name:               "test create user success",
			body:               newAccount,
			expectedStatusCode: http.StatusCreated,
			expectedSchema:     readJSONToString(t, "json_response_schema/create_user.json"),
			expectedCode:       0,
			expectedMessage:    "OK",
		},
		{
			name:               "test create user duplicate email",
			body:               newAccount,
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedSchema:     readJSONToString(t, "json_response_schema/error_422.json"),
			expectedCode:       422,
			expectedMessage:    "user email already taken",
		},
		// err 400 bad request
		{
			name:               "test invalid request body",
			body:               "someinvalidbody",
			expectedStatusCode: http.StatusBadRequest,
			expectedSchema:     readJSONToString(t, "json_response_schema/invalid_request_body.json"),
			expectedCode:       400,
			expectedMessage:    "invalid request body",
		},
		// err 422 field validation required
		{
			name: "test create user field validation required",
			body: map[string]interface{}{
				"name":  "",
				"email": gofakeit.Email(),
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedSchema:     readJSONToString(t, "json_response_schema/error_422.json"),
			expectedCode:       422,
			expectedMessage:    "validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := fastHTTPTester(t, r.Handler())

			resp := e.POST("/api/v1/users").WithJSON(tt.body).Expect()

			resp.Status(tt.expectedStatusCode)
			resp.JSON().Schema(tt.expectedSchema)
			resp.JSON().Object().Value("response_code").IsEqual(tt.expectedCode)
			resp.JSON().Object().Value("response_message").IsEqual(tt.expectedMessage)

		})
	}
}
