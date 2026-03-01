package errors

type statusCode int

const (
	StatusContinue                     statusCode = 100
	StatusSwitchingProtocols           statusCode = 101
	StatusProcessing                   statusCode = 102
	StatusEarlyHints                   statusCode = 103
	StatusOK                           statusCode = 200
	StatusCreated                      statusCode = 201
	StatusAccepted                     statusCode = 202
	StatusNonAuthoritativeInfo         statusCode = 203
	StatusNoContent                    statusCode = 204
	StatusResetContent                 statusCode = 205
	StatusPartialContent               statusCode = 206
	StatusMultiStatus                  statusCode = 207
	StatusAlreadyReported              statusCode = 208
	StatusIMUsed                       statusCode = 226
	StatusMultipleChoices              statusCode = 300
	StatusMovedPermanently             statusCode = 301
	StatusFound                        statusCode = 302
	StatusSeeOther                     statusCode = 303
	StatusNotModified                  statusCode = 304
	StatusUseProxy                     statusCode = 305
	StatusTemporaryRedirect            statusCode = 307
	StatusPermanentRedirect            statusCode = 308
	BadRequestError                    statusCode = 400
	UnauthorizedError                  statusCode = 401
	PaymentRequiredError               statusCode = 402
	ForbiddenError                     statusCode = 403
	NotFoundError                      statusCode = 404
	MethodNotAllowedError              statusCode = 405
	NotAcceptableError                 statusCode = 406
	ProxyAuthRequiredError             statusCode = 407
	RequestTimeoutError                statusCode = 408
	ConflictError                      statusCode = 409
	GoneError                          statusCode = 410
	LengthRequiredError                statusCode = 411
	PreconditionFailedError            statusCode = 412
	RequestEntityTooLargeError         statusCode = 413
	RequestURITooLongError             statusCode = 414
	UnsupportedMediaTypeError          statusCode = 415
	RequestedRangeNotSatisfiableError  statusCode = 416
	ExpectationFailedError             statusCode = 417
	TeapotError                        statusCode = 418
	MisdirectedRequestError            statusCode = 421
	UnprocessableEntityError           statusCode = 422
	LockedError                        statusCode = 423
	FailedDependencyError              statusCode = 424
	TooEarlyError                      statusCode = 425
	UpgradeRequiredError               statusCode = 426
	PreconditionRequiredError          statusCode = 428
	TooManyRequestsError               statusCode = 429
	RequestHeaderFieldsTooLargeError   statusCode = 431
	UnavailableForLegalReasonsError    statusCode = 451
	InternalServerErrorError           statusCode = 500
	NotImplementedError                statusCode = 501
	BadGatewayError                    statusCode = 502
	ServiceUnavailableError            statusCode = 503
	GatewayTimeoutError                statusCode = 504
	HTTPVersionNotSupportedError       statusCode = 505
	VariantAlsoNegotiatesError         statusCode = 506
	InsufficientStorageError           statusCode = 507
	LoopDetectedError                  statusCode = 508
	NotExtendedError                   statusCode = 510
	NetworkAuthenticationRequiredError statusCode = 511
)

type statusCodeErrorWrapping struct {
	sc  statusCode
	err error
}

func (wr *statusCodeErrorWrapping) Code() int {
	return int(wr.sc)
}

func (wr *statusCodeErrorWrapping) Unwrap() error {
	return wr.err
}

func (wr *statusCodeErrorWrapping) Error() string {
	if wr.err != nil {
		return wr.err.Error()
	}
	return ""
}

func wrapError(sc statusCode, err error) *statusCodeErrorWrapping {
	return &statusCodeErrorWrapping{
		sc:  sc,
		err: err,
	}
}

type StstusCodeErrorWrapper func(err error) error

var (
	BadRequest                    StstusCodeErrorWrapper = func(err error) error { return wrapError(BadRequestError, err) }
	Unauthorized                  StstusCodeErrorWrapper = func(err error) error { return wrapError(UnauthorizedError, err) }
	PaymentRequired               StstusCodeErrorWrapper = func(err error) error { return wrapError(PaymentRequiredError, err) }
	Forbidden                     StstusCodeErrorWrapper = func(err error) error { return wrapError(ForbiddenError, err) }
	NotFound                      StstusCodeErrorWrapper = func(err error) error { return wrapError(NotFoundError, err) }
	MethodNotAllowed              StstusCodeErrorWrapper = func(err error) error { return wrapError(MethodNotAllowedError, err) }
	NotAcceptable                 StstusCodeErrorWrapper = func(err error) error { return wrapError(NotAcceptableError, err) }
	ProxyAuthRequired             StstusCodeErrorWrapper = func(err error) error { return wrapError(ProxyAuthRequiredError, err) }
	RequestTimeout                StstusCodeErrorWrapper = func(err error) error { return wrapError(RequestTimeoutError, err) }
	Conflict                      StstusCodeErrorWrapper = func(err error) error { return wrapError(ConflictError, err) }
	Gone                          StstusCodeErrorWrapper = func(err error) error { return wrapError(GoneError, err) }
	LengthRequired                StstusCodeErrorWrapper = func(err error) error { return wrapError(LengthRequiredError, err) }
	PreconditionFailed            StstusCodeErrorWrapper = func(err error) error { return wrapError(PreconditionFailedError, err) }
	RequestEntityTooLarge         StstusCodeErrorWrapper = func(err error) error { return wrapError(RequestEntityTooLargeError, err) }
	RequestURITooLong             StstusCodeErrorWrapper = func(err error) error { return wrapError(RequestURITooLongError, err) }
	UnsupportedMediaType          StstusCodeErrorWrapper = func(err error) error { return wrapError(UnsupportedMediaTypeError, err) }
	RequestedRangeNotSatisfiable  StstusCodeErrorWrapper = func(err error) error { return wrapError(RequestedRangeNotSatisfiableError, err) }
	ExpectationFailed             StstusCodeErrorWrapper = func(err error) error { return wrapError(ExpectationFailedError, err) }
	Teapot                        StstusCodeErrorWrapper = func(err error) error { return wrapError(TeapotError, err) }
	MisdirectedRequest            StstusCodeErrorWrapper = func(err error) error { return wrapError(MisdirectedRequestError, err) }
	UnprocessableEntity           StstusCodeErrorWrapper = func(err error) error { return wrapError(UnprocessableEntityError, err) }
	Locked                        StstusCodeErrorWrapper = func(err error) error { return wrapError(LockedError, err) }
	FailedDependency              StstusCodeErrorWrapper = func(err error) error { return wrapError(FailedDependencyError, err) }
	TooEarly                      StstusCodeErrorWrapper = func(err error) error { return wrapError(TooEarlyError, err) }
	UpgradeRequired               StstusCodeErrorWrapper = func(err error) error { return wrapError(UpgradeRequiredError, err) }
	PreconditionRequired          StstusCodeErrorWrapper = func(err error) error { return wrapError(PreconditionRequiredError, err) }
	TooManyRequests               StstusCodeErrorWrapper = func(err error) error { return wrapError(TooManyRequestsError, err) }
	RequestHeaderFieldsTooLarge   StstusCodeErrorWrapper = func(err error) error { return wrapError(RequestHeaderFieldsTooLargeError, err) }
	UnavailableForLegalReasons    StstusCodeErrorWrapper = func(err error) error { return wrapError(UnavailableForLegalReasonsError, err) }
	InternalServerError           StstusCodeErrorWrapper = func(err error) error { return wrapError(InternalServerErrorError, err) }
	NotImplemented                StstusCodeErrorWrapper = func(err error) error { return wrapError(NotImplementedError, err) }
	BadGateway                    StstusCodeErrorWrapper = func(err error) error { return wrapError(BadGatewayError, err) }
	ServiceUnavailable            StstusCodeErrorWrapper = func(err error) error { return wrapError(ServiceUnavailableError, err) }
	GatewayTimeout                StstusCodeErrorWrapper = func(err error) error { return wrapError(GatewayTimeoutError, err) }
	HTTPVersionNotSupported       StstusCodeErrorWrapper = func(err error) error { return wrapError(HTTPVersionNotSupportedError, err) }
	VariantAlsoNegotiates         StstusCodeErrorWrapper = func(err error) error { return wrapError(VariantAlsoNegotiatesError, err) }
	InsufficientStorage           StstusCodeErrorWrapper = func(err error) error { return wrapError(InsufficientStorageError, err) }
	LoopDetected                  StstusCodeErrorWrapper = func(err error) error { return wrapError(LoopDetectedError, err) }
	NotExtended                   StstusCodeErrorWrapper = func(err error) error { return wrapError(NotExtendedError, err) }
	NetworkAuthenticationRequired StstusCodeErrorWrapper = func(err error) error { return wrapError(NetworkAuthenticationRequiredError, err) }
)
