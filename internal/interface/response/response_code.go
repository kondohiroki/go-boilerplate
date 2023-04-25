package response

type ResponseCode int

const (
	SYSTEM_OPERATION_SUCCESS ResponseCode = 0

	// 7xx client errors
	BAD_REQUEST        ResponseCode = 700
	UNAUTHORIZED       ResponseCode = 701
	FORBIDDEN          ResponseCode = 703
	NOT_FOUND          ResponseCode = 704
	METHOD_NOT_ALLOWED ResponseCode = 705
	VALIDATION_FAILED  ResponseCode = 760

	// 8xx server errors
	INTERNAL_SERVER_ERROR ResponseCode = 800
	NOT_IMPLEMENTED       ResponseCode = 801
	SERVUCE_UNAVAILABLE   ResponseCode = 803

	// custom errors
	BACKGROUND_JOB_FAILED ResponseCode = 810
)
