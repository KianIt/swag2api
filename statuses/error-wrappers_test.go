package statuses

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusCodeWrappingError_Code(t *testing.T) {
	wr := statusCodeWrappingError{
		sc:  statusOK,
		err: errors.New("test"),
	}
	assert.Equal(t, int(statusOK), wr.Code())
}

func TestStatusCodeWrappingError_Unwrap(t *testing.T) {
	wr := statusCodeWrappingError{
		sc:  statusOK,
		err: errors.New("test"),
	}
	assert.Equal(t, errors.New("test"), wr.Unwrap())
}

func TestStatusCodeWrappingError_Error(t *testing.T) {
	wr := statusCodeWrappingError{
		sc:  statusOK,
		err: errors.New("test"),
	}
	assert.Equal(t, "test", wr.Error())
}

func TestStstusCodeErrorWrapper(t *testing.T) {
	tt := []struct {
		name    string
		sc      statusCode
		wrapper StatusCodeErrorWrapper
	}{
		{name: "Continue", sc: statusContinue, wrapper: Continue},
		{name: "SwitchingProtocols", sc: statusSwitchingProtocols, wrapper: SwitchingProtocols},
		{name: "Processing", sc: statusProcessing, wrapper: Processing},
		{name: "EarlyHints", sc: statusEarlyHints, wrapper: EarlyHints},
		{name: "OK", sc: statusOK, wrapper: OK},
		{name: "Created", sc: statusCreated, wrapper: Created},
		{name: "Accepted", sc: statusAccepted, wrapper: Accepted},
		{name: "NonAuthoritativeInfo", sc: statusNonAuthoritativeInfo, wrapper: NonAuthoritativeInfo},
		{name: "NoContent", sc: statusNoContent, wrapper: NoContent},
		{name: "ResetContent", sc: statusResetContent, wrapper: ResetContent},
		{name: "PartialContent", sc: statusPartialContent, wrapper: PartialContent},
		{name: "MultiStatus", sc: statusMultiStatus, wrapper: MultiStatus},
		{name: "AlreadyReported", sc: statusAlreadyReported, wrapper: AlreadyReported},
		{name: "IMUsed", sc: statusIMUsed, wrapper: IMUsed},
		{name: "MultipleChoices", sc: statusMultipleChoices, wrapper: MultipleChoices},
		{name: "MovedPermanently", sc: statusMovedPermanently, wrapper: MovedPermanently},
		{name: "Found", sc: statusFound, wrapper: Found},
		{name: "SeeOther", sc: statusSeeOther, wrapper: SeeOther},
		{name: "NotModified", sc: statusNotModified, wrapper: NotModified},
		{name: "UseProxy", sc: statusUseProxy, wrapper: UseProxy},
		{name: "TemporaryRedirect", sc: statusTemporaryRedirect, wrapper: TemporaryRedirect},
		{name: "PermanentRedirect", sc: statusPermanentRedirect, wrapper: PermanentRedirect},
		{name: "BadRequestError", sc: statusBadRequestError, wrapper: BadRequestError},
		{name: "UnauthorizedError", sc: statusUnauthorizedError, wrapper: UnauthorizedError},
		{name: "PaymentRequiredError", sc: statusPaymentRequiredError, wrapper: PaymentRequiredError},
		{name: "ForbiddenError", sc: statusForbiddenError, wrapper: ForbiddenError},
		{name: "NotFoundError", sc: statusNotFoundError, wrapper: NotFoundError},
		{name: "MethodNotAllowedError", sc: statusMethodNotAllowedError, wrapper: MethodNotAllowedError},
		{name: "NotAcceptableError", sc: statusNotAcceptableError, wrapper: NotAcceptableError},
		{name: "ProxyAuthRequiredError", sc: statusProxyAuthRequiredError, wrapper: ProxyAuthRequiredError},
		{name: "RequestTimeoutError", sc: statusRequestTimeoutError, wrapper: RequestTimeoutError},
		{name: "ConflictError", sc: statusConflictError, wrapper: ConflictError},
		{name: "GoneError", sc: statusGoneError, wrapper: GoneError},
		{name: "LengthRequiredError", sc: statusLengthRequiredError, wrapper: LengthRequiredError},
		{name: "PreconditionFailedError", sc: statusPreconditionFailedError, wrapper: PreconditionFailedError},
		{name: "RequestEntityTooLargeError", sc: statusRequestEntityTooLargeError, wrapper: RequestEntityTooLargeError},
		{name: "RequestURITooLongError", sc: statusRequestURITooLongError, wrapper: RequestURITooLongError},
		{name: "UnsupportedMediaTypeError", sc: statusUnsupportedMediaTypeError, wrapper: UnsupportedMediaTypeError},
		{name: "RequestedRangeNotSatisfiableError", sc: statusRequestedRangeNotSatisfiableError, wrapper: RequestedRangeNotSatisfiableError},
		{name: "ExpectationFailedError", sc: statusExpectationFailedError, wrapper: ExpectationFailedError},
		{name: "TeapotError", sc: statusTeapotError, wrapper: TeapotError},
		{name: "MisdirectedRequestError", sc: statusMisdirectedRequestError, wrapper: MisdirectedRequestError},
		{name: "UnprocessableEntityError", sc: statusUnprocessableEntityError, wrapper: UnprocessableEntityError},
		{name: "LockedError", sc: statusLockedError, wrapper: LockedError},
		{name: "FailedDependencyError", sc: statusFailedDependencyError, wrapper: FailedDependencyError},
		{name: "TooEarlyError", sc: statusTooEarlyError, wrapper: TooEarlyError},
		{name: "UpgradeRequiredError", sc: statusUpgradeRequiredError, wrapper: UpgradeRequiredError},
		{name: "PreconditionRequiredError", sc: statusPreconditionRequiredError, wrapper: PreconditionRequiredError},
		{name: "TooManyRequestsError", sc: statusTooManyRequestsError, wrapper: TooManyRequestsError},
		{name: "RequestHeaderFieldsTooLargeError", sc: statusRequestHeaderFieldsTooLargeError, wrapper: RequestHeaderFieldsTooLargeError},
		{name: "UnavailableForLegalReasonsError", sc: statusUnavailableForLegalReasonsError, wrapper: UnavailableForLegalReasonsError},
		{name: "InternalServerErrorError", sc: statusInternalServerErrorError, wrapper: InternalServerErrorError},
		{name: "NotImplementedError", sc: statusNotImplementedError, wrapper: NotImplementedError},
		{name: "BadGatewayError", sc: statusBadGatewayError, wrapper: BadGatewayError},
		{name: "ServiceUnavailableError", sc: statusServiceUnavailableError, wrapper: ServiceUnavailableError},
		{name: "GatewayTimeoutError", sc: statusGatewayTimeoutError, wrapper: GatewayTimeoutError},
		{name: "HTTPVersionNotSupportedError", sc: statusHTTPVersionNotSupportedError, wrapper: HTTPVersionNotSupportedError},
		{name: "VariantAlsoNegotiatesError", sc: statusVariantAlsoNegotiatesError, wrapper: VariantAlsoNegotiatesError},
		{name: "InsufficientStorageError", sc: statusInsufficientStorageError, wrapper: InsufficientStorageError},
		{name: "LoopDetectedError", sc: statusLoopDetectedError, wrapper: LoopDetectedError},
		{name: "NotExtendedError", sc: statusNotExtendedError, wrapper: NotExtendedError},
		{name: "NetworkAuthenticationRequiredError", sc: statusNetworkAuthenticationRequiredError, wrapper: NetworkAuthenticationRequiredError},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.sc, tc.wrapper(nil).(*statusCodeWrappingError).sc) //nolint:errorlint // Need to assret error type.
		})
	}
}
