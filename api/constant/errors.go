package constant

type APIError struct {
	HTTPStatusCode int   `json:"http_status_code"`
	ErrorMessage   error `json:"error_message"`
}

func (apiError *APIError) Error() string {
	return apiError.ErrorMessage.Error()
}

func APIErrorBadRequest(err error) *APIError {
	return &APIError{
		HTTPStatusCode: HTTPResponseCodeBadRequest,
		ErrorMessage:   err,
	}
}

func APIErrorForbidden(err error) *APIError {
	return &APIError{
		HTTPStatusCode: HTTPResponseCodeForbidden,
		ErrorMessage:   err,
	}
}

func APIErrorInternalServerError(err error) *APIError {
	return &APIError{
		HTTPStatusCode: HTTPResponseCodeInternalServerError,
		ErrorMessage:   err,
	}
}
