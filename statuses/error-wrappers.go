package statuses

// statusCode is an HTTP status code.
type statusCode int

const (
	statusContinue                           statusCode = 100
	statusSwitchingProtocols                 statusCode = 101
	statusProcessing                         statusCode = 102
	statusEarlyHints                         statusCode = 103
	statusOK                                 statusCode = 200
	statusCreated                            statusCode = 201
	statusAccepted                           statusCode = 202
	statusNonAuthoritativeInfo               statusCode = 203
	statusNoContent                          statusCode = 204
	statusResetContent                       statusCode = 205
	statusPartialContent                     statusCode = 206
	statusMultiStatus                        statusCode = 207
	statusAlreadyReported                    statusCode = 208
	statusIMUsed                             statusCode = 226
	statusMultipleChoices                    statusCode = 300
	statusMovedPermanently                   statusCode = 301
	statusFound                              statusCode = 302
	statusSeeOther                           statusCode = 303
	statusNotModified                        statusCode = 304
	statusUseProxy                           statusCode = 305
	statusTemporaryRedirect                  statusCode = 307
	statusPermanentRedirect                  statusCode = 308
	statusBadRequestError                    statusCode = 400
	statusUnauthorizedError                  statusCode = 401
	statusPaymentRequiredError               statusCode = 402
	statusForbiddenError                     statusCode = 403
	statusNotFoundError                      statusCode = 404
	statusMethodNotAllowedError              statusCode = 405
	statusNotAcceptableError                 statusCode = 406
	statusProxyAuthRequiredError             statusCode = 407
	statusRequestTimeoutError                statusCode = 408
	statusConflictError                      statusCode = 409
	statusGoneError                          statusCode = 410
	statusLengthRequiredError                statusCode = 411
	statusPreconditionFailedError            statusCode = 412
	statusRequestEntityTooLargeError         statusCode = 413
	statusRequestURITooLongError             statusCode = 414
	statusUnsupportedMediaTypeError          statusCode = 415
	statusRequestedRangeNotSatisfiableError  statusCode = 416
	statusExpectationFailedError             statusCode = 417
	statusTeapotError                        statusCode = 418
	statusMisdirectedRequestError            statusCode = 421
	statusUnprocessableEntityError           statusCode = 422
	statusLockedError                        statusCode = 423
	statusFailedDependencyError              statusCode = 424
	statusTooEarlyError                      statusCode = 425
	statusUpgradeRequiredError               statusCode = 426
	statusPreconditionRequiredError          statusCode = 428
	statusTooManyRequestsError               statusCode = 429
	statusRequestHeaderFieldsTooLargeError   statusCode = 431
	statusUnavailableForLegalReasonsError    statusCode = 451
	statusInternalServerErrorError           statusCode = 500
	statusNotImplementedError                statusCode = 501
	statusBadGatewayError                    statusCode = 502
	statusServiceUnavailableError            statusCode = 503
	statusGatewayTimeoutError                statusCode = 504
	statusHTTPVersionNotSupportedError       statusCode = 505
	statusVariantAlsoNegotiatesError         statusCode = 506
	statusInsufficientStorageError           statusCode = 507
	statusLoopDetectedError                  statusCode = 508
	statusNotExtendedError                   statusCode = 510
	statusNetworkAuthenticationRequiredError statusCode = 511
)

// statusCodeWrappingError is a status code error wrapper.
//
// Holds an error with a status code.
type statusCodeWrappingError struct {
	sc  statusCode
	err error
}

// Code returns the status code.
func (wr *statusCodeWrappingError) Code() int {
	return int(wr.sc)
}

// Unwrap returns the wrapped error.
func (wr *statusCodeWrappingError) Unwrap() error {
	return wr.err
}

// Error returns the wrapped error message.
func (wr *statusCodeWrappingError) Error() string {
	if wr.err != nil {
		return wr.err.Error()
	}
	return ""
}

// wrapError wraps an error with a status code.
func wrapError(sc statusCode, err error) *statusCodeWrappingError {
	return &statusCodeWrappingError{
		sc:  sc,
		err: err,
	}
}

// StatusCodeErrorWrapper is a status code error wrapper.
//
// Wraps an error with some status code.
type StatusCodeErrorWrapper func(err error) error

var (
	Continue                           StatusCodeErrorWrapper = func(err error) error { return wrapError(statusContinue, err) }
	SwitchingProtocols                 StatusCodeErrorWrapper = func(err error) error { return wrapError(statusSwitchingProtocols, err) }
	Processing                         StatusCodeErrorWrapper = func(err error) error { return wrapError(statusProcessing, err) }
	EarlyHints                         StatusCodeErrorWrapper = func(err error) error { return wrapError(statusEarlyHints, err) }
	OK                                 StatusCodeErrorWrapper = func(err error) error { return wrapError(statusOK, err) }
	Created                            StatusCodeErrorWrapper = func(err error) error { return wrapError(statusCreated, err) }
	Accepted                           StatusCodeErrorWrapper = func(err error) error { return wrapError(statusAccepted, err) }
	NonAuthoritativeInfo               StatusCodeErrorWrapper = func(err error) error { return wrapError(statusNonAuthoritativeInfo, err) }
	NoContent                          StatusCodeErrorWrapper = func(err error) error { return wrapError(statusNoContent, err) }
	ResetContent                       StatusCodeErrorWrapper = func(err error) error { return wrapError(statusResetContent, err) }
	PartialContent                     StatusCodeErrorWrapper = func(err error) error { return wrapError(statusPartialContent, err) }
	MultiStatus                        StatusCodeErrorWrapper = func(err error) error { return wrapError(statusMultiStatus, err) }
	AlreadyReported                    StatusCodeErrorWrapper = func(err error) error { return wrapError(statusAlreadyReported, err) }
	IMUsed                             StatusCodeErrorWrapper = func(err error) error { return wrapError(statusIMUsed, err) }
	MultipleChoices                    StatusCodeErrorWrapper = func(err error) error { return wrapError(statusMultipleChoices, err) }
	MovedPermanently                   StatusCodeErrorWrapper = func(err error) error { return wrapError(statusMovedPermanently, err) }
	Found                              StatusCodeErrorWrapper = func(err error) error { return wrapError(statusFound, err) }
	SeeOther                           StatusCodeErrorWrapper = func(err error) error { return wrapError(statusSeeOther, err) }
	NotModified                        StatusCodeErrorWrapper = func(err error) error { return wrapError(statusNotModified, err) }
	UseProxy                           StatusCodeErrorWrapper = func(err error) error { return wrapError(statusUseProxy, err) }
	TemporaryRedirect                  StatusCodeErrorWrapper = func(err error) error { return wrapError(statusTemporaryRedirect, err) }
	PermanentRedirect                  StatusCodeErrorWrapper = func(err error) error { return wrapError(statusPermanentRedirect, err) }
	BadRequestError                    StatusCodeErrorWrapper = func(err error) error { return wrapError(statusBadRequestError, err) }
	UnauthorizedError                  StatusCodeErrorWrapper = func(err error) error { return wrapError(statusUnauthorizedError, err) }
	PaymentRequiredError               StatusCodeErrorWrapper = func(err error) error { return wrapError(statusPaymentRequiredError, err) }
	ForbiddenError                     StatusCodeErrorWrapper = func(err error) error { return wrapError(statusForbiddenError, err) }
	NotFoundError                      StatusCodeErrorWrapper = func(err error) error { return wrapError(statusNotFoundError, err) }
	MethodNotAllowedError              StatusCodeErrorWrapper = func(err error) error { return wrapError(statusMethodNotAllowedError, err) }
	NotAcceptableError                 StatusCodeErrorWrapper = func(err error) error { return wrapError(statusNotAcceptableError, err) }
	ProxyAuthRequiredError             StatusCodeErrorWrapper = func(err error) error { return wrapError(statusProxyAuthRequiredError, err) }
	RequestTimeoutError                StatusCodeErrorWrapper = func(err error) error { return wrapError(statusRequestTimeoutError, err) }
	ConflictError                      StatusCodeErrorWrapper = func(err error) error { return wrapError(statusConflictError, err) }
	GoneError                          StatusCodeErrorWrapper = func(err error) error { return wrapError(statusGoneError, err) }
	LengthRequiredError                StatusCodeErrorWrapper = func(err error) error { return wrapError(statusLengthRequiredError, err) }
	PreconditionFailedError            StatusCodeErrorWrapper = func(err error) error { return wrapError(statusPreconditionFailedError, err) }
	RequestEntityTooLargeError         StatusCodeErrorWrapper = func(err error) error { return wrapError(statusRequestEntityTooLargeError, err) }
	RequestURITooLongError             StatusCodeErrorWrapper = func(err error) error { return wrapError(statusRequestURITooLongError, err) }
	UnsupportedMediaTypeError          StatusCodeErrorWrapper = func(err error) error { return wrapError(statusUnsupportedMediaTypeError, err) }
	RequestedRangeNotSatisfiableError  StatusCodeErrorWrapper = func(err error) error { return wrapError(statusRequestedRangeNotSatisfiableError, err) }
	ExpectationFailedError             StatusCodeErrorWrapper = func(err error) error { return wrapError(statusExpectationFailedError, err) }
	TeapotError                        StatusCodeErrorWrapper = func(err error) error { return wrapError(statusTeapotError, err) }
	MisdirectedRequestError            StatusCodeErrorWrapper = func(err error) error { return wrapError(statusMisdirectedRequestError, err) }
	UnprocessableEntityError           StatusCodeErrorWrapper = func(err error) error { return wrapError(statusUnprocessableEntityError, err) }
	LockedError                        StatusCodeErrorWrapper = func(err error) error { return wrapError(statusLockedError, err) }
	FailedDependencyError              StatusCodeErrorWrapper = func(err error) error { return wrapError(statusFailedDependencyError, err) }
	TooEarlyError                      StatusCodeErrorWrapper = func(err error) error { return wrapError(statusTooEarlyError, err) }
	UpgradeRequiredError               StatusCodeErrorWrapper = func(err error) error { return wrapError(statusUpgradeRequiredError, err) }
	PreconditionRequiredError          StatusCodeErrorWrapper = func(err error) error { return wrapError(statusPreconditionRequiredError, err) }
	TooManyRequestsError               StatusCodeErrorWrapper = func(err error) error { return wrapError(statusTooManyRequestsError, err) }
	RequestHeaderFieldsTooLargeError   StatusCodeErrorWrapper = func(err error) error { return wrapError(statusRequestHeaderFieldsTooLargeError, err) }
	UnavailableForLegalReasonsError    StatusCodeErrorWrapper = func(err error) error { return wrapError(statusUnavailableForLegalReasonsError, err) }
	InternalServerErrorError           StatusCodeErrorWrapper = func(err error) error { return wrapError(statusInternalServerErrorError, err) }
	NotImplementedError                StatusCodeErrorWrapper = func(err error) error { return wrapError(statusNotImplementedError, err) }
	BadGatewayError                    StatusCodeErrorWrapper = func(err error) error { return wrapError(statusBadGatewayError, err) }
	ServiceUnavailableError            StatusCodeErrorWrapper = func(err error) error { return wrapError(statusServiceUnavailableError, err) }
	GatewayTimeoutError                StatusCodeErrorWrapper = func(err error) error { return wrapError(statusGatewayTimeoutError, err) }
	HTTPVersionNotSupportedError       StatusCodeErrorWrapper = func(err error) error { return wrapError(statusHTTPVersionNotSupportedError, err) }
	VariantAlsoNegotiatesError         StatusCodeErrorWrapper = func(err error) error { return wrapError(statusVariantAlsoNegotiatesError, err) }
	InsufficientStorageError           StatusCodeErrorWrapper = func(err error) error { return wrapError(statusInsufficientStorageError, err) }
	LoopDetectedError                  StatusCodeErrorWrapper = func(err error) error { return wrapError(statusLoopDetectedError, err) }
	NotExtendedError                   StatusCodeErrorWrapper = func(err error) error { return wrapError(statusNotExtendedError, err) }
	NetworkAuthenticationRequiredError StatusCodeErrorWrapper = func(err error) error { return wrapError(statusNetworkAuthenticationRequiredError, err) }
)
