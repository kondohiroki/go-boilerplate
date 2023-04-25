package exception

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/kondohiroki/go-boilerplate/internal/interface/validation"
)

// NewValidationFailedErrors reads ValidationErrors and converts to format of ExceptionErrors.
func NewValidationFailedErrors(validationErrs validator.ValidationErrors) *ExceptionErrors {
	errItems := make([]*ExceptionError, 0, len(validationErrs))
	for _, validationErr := range validationErrs {
		errItems = append(errItems, &ExceptionError{
			Message:      validation.Translate(validationErr),
			Type:         ERROR_TYPE_VALIDATION_ERROR,
			ErrorSubcode: SUBCODE_VALIDATION_FAILED,
		})
	}
	return &ExceptionErrors{
		GlobalMessage:  "validation failed",
		HttpStatusCode: http.StatusUnprocessableEntity,
		ErrItems:       errItems,
	}
}

// IsEmpty reports whether there is at least one errItems in an instance.
func (cErrs *ExceptionErrors) IsEmpty() bool {
	return len(cErrs.ErrItems) == 0
}

// Append manually appends error item to current errors.
func (cErrs *ExceptionErrors) Append(cErr *ExceptionError) *ExceptionErrors {
	cErrs.ErrItems = append(cErrs.ErrItems, cErr)
	return cErrs
}

// Datasource Errors

func (cErrs *ExceptionErrors) AppendInputFieldIsNotConfigured(fieldName, productName string) *ExceptionErrors {
	cErrs.ErrItems = append(cErrs.ErrItems, &ExceptionError{
		Message:      fmt.Sprintf("input field: %s of product_name: %s is not configured", fieldName, productName),
		Type:         ERROR_TYPE_DATASOURCE_ERROR,
		ErrorSubcode: SUBCODE_INPUT_FIELD_IS_NOT_CONFIGURED,
	})
	return cErrs
}

func (cErrs *ExceptionErrors) AppendInvalidFieldValuesLength() *ExceptionErrors {
	cErrs.ErrItems = append(cErrs.ErrItems, &ExceptionError{
		Message:      "length of field values must be the same for every input fields",
		Type:         ERROR_TYPE_DATASOURCE_ERROR,
		ErrorSubcode: SUBCODE_NUM_MULTIPLE_VALUES_ERROR,
	})
	return cErrs
}

func (cErrs *ExceptionErrors) AppendInvalidFieldValue(fieldName, fieldValue string) *ExceptionErrors {
	cErrs.ErrItems = append(cErrs.ErrItems, &ExceptionError{
		Message:      fmt.Sprintf("invalid format value: %s of field: %s", fieldValue, fieldName),
		Type:         ERROR_TYPE_DATASOURCE_ERROR,
		ErrorSubcode: SUBCODE_INVALID_FIELD_VALUE_FORMAT,
	})
	return cErrs
}
