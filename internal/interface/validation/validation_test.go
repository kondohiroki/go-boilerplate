package validation_test

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-playground/validator/v10"
	. "github.com/kondohiroki/go-boilerplate/internal/interface/validation"
)

type testStruct struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"min=18"`
}

func TestInitValidator(t *testing.T) {
	InitValidator()
	validate, trans := GetValidator()

	require.NotNil(t, validate)
	require.NotNil(t, trans)
}

func TestGetValidator(t *testing.T) {
	validate, trans := GetValidator()

	require.NotNil(t, validate)
	require.NotNil(t, trans)
}

func TestGetValidationErrors(t *testing.T) {
	validate, _ := GetValidator()

	testData := testStruct{
		Username: "",
		Email:    "invalid_email",
		Age:      16,
	}

	err := validate.Struct(testData)
	require.Error(t, err)

	validationErrors, ok := err.(validator.ValidationErrors)
	require.True(t, ok)

	errors := GetValidationErrors(validationErrors)
	require.Len(t, errors, 3)

	expectedErrors := []map[string]string{
		{"username": "username is a required field"},
		{"email": "email must be a valid email address"},
		{"age": "age must be 18 or greater"},
	}

	for idx, err := range errors {
		errMap := make(map[string]string)
		errJson, err := json.Marshal(err)
		require.NoError(t, err)
		err = json.Unmarshal(errJson, &errMap)
		require.NoError(t, err)

		assert.Equal(t, expectedErrors[idx], errMap)
	}
}

func TestValidatorTranslation(t *testing.T) {
	tests := []struct {
		name          string
		fieldName     string
		fieldValue    interface{}
		tag           string
		expectedError string
	}{
		{
			name:          "required",
			fieldName:     "username",
			fieldValue:    "",
			tag:           "required",
			expectedError: "username is a required field",
		},
		{
			name:          "email",
			fieldName:     "email",
			fieldValue:    "invalid_email",
			tag:           "email",
			expectedError: "email must be a valid email address",
		},
		{
			name:          "min",
			fieldName:     "age",
			fieldValue:    16,
			tag:           "min",
			expectedError: "age must be 18 or greater",
		},
	}

	validate, trans := GetValidator()
	require.NotNil(t, validate)
	require.NotNil(t, trans)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testData := testStruct{
				Username: "valid_username",
				Email:    "valid@example.com",
				Age:      18,
			}

			value := reflect.ValueOf(&testData).Elem().FieldByNameFunc(func(fieldName string) bool {
				return strings.ToLower(fieldName) == tt.fieldName
			})
			require.True(t, value.IsValid())
			value.Set(reflect.ValueOf(tt.fieldValue))

			err := validate.Struct(testData)
			require.Error(t, err)

			validationErrors, ok := err.(validator.ValidationErrors)
			require.True(t, ok)

			require.Len(t, validationErrors, 1)
			validationError := validationErrors[0]

			assert.Equal(t, tt.fieldName, validationError.Field())
			assert.Equal(t, tt.expectedError, validationError.Translate(trans))
		})
	}
}
