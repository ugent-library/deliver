package httperror

import (
	"fmt"
	"net/http"
)

var (
	BadRequest                   = New(http.StatusBadRequest)
	Unauthorized                 = New(http.StatusUnauthorized)
	PaymentRequired              = New(http.StatusPaymentRequired)
	Forbidden                    = New(http.StatusForbidden)
	NotFound                     = New(http.StatusNotFound)
	MethodNotAllowed             = New(http.StatusMethodNotAllowed)
	NotAcceptable                = New(http.StatusNotAcceptable)
	ProxyAuthRequired            = New(http.StatusProxyAuthRequired)
	RequestTimeout               = New(http.StatusRequestTimeout)
	Conflict                     = New(http.StatusConflict)
	Gone                         = New(http.StatusGone)
	LengthRequired               = New(http.StatusLengthRequired)
	PreconditionFailed           = New(http.StatusPreconditionFailed)
	RequestEntityTooLarge        = New(http.StatusRequestEntityTooLarge)
	RequestURITooLong            = New(http.StatusRequestURITooLong)
	UnsupportedMediaType         = New(http.StatusUnsupportedMediaType)
	RequestedRangeNotSatisfiable = New(http.StatusRequestedRangeNotSatisfiable)
	ExpectationFailed            = New(http.StatusExpectationFailed)
	Teapot                       = New(http.StatusTeapot)
	MisdirectedRequest           = New(http.StatusMisdirectedRequest)
	UnprocessableEntity          = New(http.StatusUnprocessableEntity)
	Locked                       = New(http.StatusLocked)
	FailedDependency             = New(http.StatusFailedDependency)
	TooEarly                     = New(http.StatusTooEarly)
	UpgradeRequired              = New(http.StatusUpgradeRequired)
	PreconditionRequired         = New(http.StatusPreconditionRequired)
	TooManyRequests              = New(http.StatusTooManyRequests)
	RequestHeaderFieldsTooLarge  = New(http.StatusRequestHeaderFieldsTooLarge)
	UnavailableForLegalReasons   = New(http.StatusUnavailableForLegalReasons)

	InternalServerError           = New(http.StatusInternalServerError)
	NotImplemented                = New(http.StatusNotImplemented)
	BadGateway                    = New(http.StatusBadGateway)
	ServiceUnavailable            = New(http.StatusServiceUnavailable)
	GatewayTimeout                = New(http.StatusGatewayTimeout)
	HTTPVersionNotSupported       = New(http.StatusHTTPVersionNotSupported)
	VariantAlsoNegotiates         = New(http.StatusVariantAlsoNegotiates)
	InsufficientStorage           = New(http.StatusInsufficientStorage)
	LoopDetected                  = New(http.StatusLoopDetected)
	NotExtended                   = New(http.StatusNotExtended)
	NetworkAuthenticationRequired = New(http.StatusNetworkAuthenticationRequired)
)

type Error struct {
	Code int
}

func New(code int) *Error {
	return &Error{code}
}

func (e *Error) Error() string {
	return fmt.Sprintf("http error %d: %s", e.Code, http.StatusText(e.Code))
}
