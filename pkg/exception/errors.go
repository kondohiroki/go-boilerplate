package exception

import "net/http"

/*
Exported error instances.
Error instances should be equal to subcodes and classified by error type.
Error message can be any string.
*/
var (
	// UnknownError
	UnknownError *ExceptionErrors = createFixedExceptionErrors(
		http.StatusInternalServerError,
		ERROR_TYPE_UNKNOWN_ERROR,
		SUBCODE_UNKNOWN_ERROR,
		"an error is occurred",
	)

	// BadRequest
	BadRequestError *ExceptionErrors = createFixedExceptionErrors(
		http.StatusBadRequest,
		ERROR_TYPE_BAD_REQUEST,
		SUBCODE_BAD_REQUEST,
		"bad request",
	)
	InvalidRequestBodyError *ExceptionErrors = createFixedExceptionErrors(
		http.StatusBadRequest,
		ERROR_TYPE_BAD_REQUEST,
		SUBCODE_INVALID_REQUEST_BODY,
		"invalid request body",
	)
	InvalidRequestQueryParamError *ExceptionErrors = createFixedExceptionErrors(
		http.StatusBadRequest,
		ERROR_TYPE_BAD_REQUEST,
		SUBCODE_INVALID_REQUEST_BODY,
		"invalid request query parameter",
	)
	InvalidIDError *ExceptionErrors = createFixedExceptionErrors(
		http.StatusBadRequest,
		ERROR_TYPE_BAD_REQUEST,
		SUBCODE_INVALID_ID,
		"invalid ID",
	)

	// DataNotFound
	DataNotFoundError *ExceptionErrors = createFixedExceptionErrors(
		http.StatusNotFound,
		ERROR_TYPE_NOT_FOUND,
		SUBCODE_DATA_NOT_FOUND,
		"data is not found",
	)
	ApiNotFoundError *ExceptionErrors = createFixedExceptionErrors(
		http.StatusNotFound,
		ERROR_TYPE_NOT_FOUND,
		SUBCODE_API_NOTE_FOUND,
		"this is not api you are looking for",
	)

	// Unauthorized
	UnauthorizedError *ExceptionErrors = createFixedExceptionErrors(
		http.StatusUnauthorized,
		ERROR_TYPE_UNAUTHORIZED,
		SUBCODE_UNAUTHORIZED,
		"permission is not granted",
	)

	// ValidationError
	ValidationFailedError *ExceptionErrors = createFixedExceptionErrors(
		http.StatusUnprocessableEntity,
		ERROR_TYPE_VALIDATION_ERROR,
		SUBCODE_VALIDATION_FAILED,
		"validation failed",
	)
	UserEmailAlreadyTakenError *ExceptionErrors = createFixedExceptionErrors(
		http.StatusUnprocessableEntity,
		ERROR_TYPE_VALIDATION_ERROR,
		SUBCODE_USER_EMAIL_ALREADY_TAKEN,
		"user email already taken",
	)

	// JobError
	BackgroundJobFailedError *ExceptionErrors = createFixedExceptionErrors(
		http.StatusInternalServerError,
		ERROR_TYPE_JOB_ERROR,
		SUBCODE_BACKGROUND_JOB_FAILED,
		"background job failed",
	)
	CannotRunBatchDailyError *ExceptionErrors = createFixedExceptionErrors(
		http.StatusInternalServerError,
		ERROR_TYPE_JOB_ERROR,
		SUBCODE_CANNOT_RUN_BATCH_DAILY,
		"cannot run batch daily",
	)

	// ExternalServiceError
	CallNCBError *ExceptionErrors = createFixedExceptionErrors(
		http.StatusBadGateway,
		ERROR_TYPE_EXTERNAL_SERVICE_ERROR,
		SUBCODE_CALL_NCB_ERROR,
		"could not get NCB data",
	)
	ApplyProductError *ExceptionErrors = createFixedExceptionErrors(
		http.StatusBadGateway,
		ERROR_TYPE_EXTERNAL_SERVICE_ERROR,
		SUBCODE_APPLY_PRODUCT_ERROR,
		"could not apply for product",
	)
	AnswerQuestionError *ExceptionErrors = createFixedExceptionErrors(
		http.StatusBadGateway,
		ERROR_TYPE_EXTERNAL_SERVICE_ERROR,
		SUBCODE_ANSWER_QUESTION_ERROR,
		"could not answer question",
	)
	GetCaseDetailError *ExceptionErrors = createFixedExceptionErrors(
		http.StatusBadGateway,
		ERROR_TYPE_EXTERNAL_SERVICE_ERROR,
		SUBCODE_GET_CASE_DETAIL_ERROR,
		"could not get case detail",
	)

	// DatasourceError
	ResponseFieldNotFoundError *ExceptionErrors = createFixedExceptionErrors(
		http.StatusBadGateway,
		ERROR_TYPE_DATASOURCE_ERROR,
		SUBCODE_RESPONSE_FIELD_NOT_FOUND_ERROR,
		"result field not found",
	)
)
