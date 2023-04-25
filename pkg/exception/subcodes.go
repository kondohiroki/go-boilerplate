package exception

const _SYSTEM_CODE int = 10

type errorSubcode int

func newErrorSubcode(code int) errorSubcode {
	return errorSubcode(_SYSTEM_CODE*1000 + code)
}

var (
	// 7xx client errors
	SUBCODE_BAD_REQUEST                    errorSubcode = newErrorSubcode(700)
	SUBCODE_INVALID_REQUEST_BODY           errorSubcode = newErrorSubcode(798)
	SUBCODE_INVALID_ID                     errorSubcode = newErrorSubcode(799)
	SUBCODE_UNAUTHORIZED                   errorSubcode = newErrorSubcode(701)
	SUBCODE_DATA_NOT_FOUND                 errorSubcode = newErrorSubcode(704)
	SUBCODE_API_NOTE_FOUND                 errorSubcode = newErrorSubcode(705)
	SUBCODE_VALIDATION_FAILED              errorSubcode = newErrorSubcode(760)
	SUBCODE_USER_EMAIL_ALREADY_TAKEN       errorSubcode = newErrorSubcode(761)
	SUBCODE_INPUT_FIELD_IS_NOT_CONFIGURED  errorSubcode = newErrorSubcode(762)
	SUBCODE_INVALID_FIELD_VALUE_FORMAT     errorSubcode = newErrorSubcode(762)
	SUBCODE_NUM_MULTIPLE_VALUES_ERROR      errorSubcode = newErrorSubcode(763)
	SUBCODE_RESPONSE_FIELD_NOT_FOUND_ERROR errorSubcode = newErrorSubcode(764)

	// 8xx server errors
	SUBCODE_UNKNOWN_ERROR          errorSubcode = newErrorSubcode(800)
	SUBCODE_BACKGROUND_JOB_FAILED  errorSubcode = newErrorSubcode(810)
	SUBCODE_CANNOT_RUN_BATCH_DAILY errorSubcode = newErrorSubcode(811)
	SUBCODE_CALL_NCB_ERROR         errorSubcode = newErrorSubcode(812)
	SUBCODE_AUTH_CORE_ERROR        errorSubcode = newErrorSubcode(813)
	SUBCODE_APPLY_PRODUCT_ERROR    errorSubcode = newErrorSubcode(814)
	SUBCODE_ANSWER_QUESTION_ERROR  errorSubcode = newErrorSubcode(815)
	SUBCODE_GET_CASE_DETAIL_ERROR  errorSubcode = newErrorSubcode(816)
)
