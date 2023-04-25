package exception

type errorType string

const (
	ERROR_TYPE_UNKNOWN_ERROR          errorType = "UnknownError"
	ERROR_TYPE_BAD_REQUEST            errorType = "BadRequest"
	ERROR_TYPE_NOT_FOUND              errorType = "NotFound"
	ERROR_TYPE_UNAUTHORIZED           errorType = "Unauthorized"
	ERROR_TYPE_VALIDATION_ERROR       errorType = "ValidationError"
	ERROR_TYPE_JOB_ERROR              errorType = "JobError"
	ERROR_TYPE_EXTERNAL_SERVICE_ERROR errorType = "ExternalServiceError"
	ERROR_TYPE_DATASOURCE_ERROR       errorType = "DatasourceError"
)
