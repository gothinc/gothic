package httpclient

/**
 * @desc Status
 * @author zhaojiangwei
 * @date 2018-12-25 14:59
 */

type HttpStatusCode int
const (
	StatusContinue          		HttpStatusCode 	= 100 // RFC 7231, 6.2.1
	StatusSwitchingProtocols 	HttpStatusCode	= 101 // RFC 7231, 6.2.2
	StatusProcessing         		HttpStatusCode	= 102 // RFC 2518, 10.1

	StatusOK                   	HttpStatusCode	= 200 // RFC 7231, 6.3.1
	StatusCreated              	HttpStatusCode	= 201 // RFC 7231, 6.3.2
	StatusAccepted             	HttpStatusCode	= 202 // RFC 7231, 6.3.3
	StatusNonAuthoritativeInfo 	HttpStatusCode	= 203 // RFC 7231, 6.3.4
	StatusNoContent            	HttpStatusCode	= 204 // RFC 7231, 6.3.5
	StatusResetContent         	HttpStatusCode	= 205 // RFC 7231, 6.3.6
	StatusPartialContent       	HttpStatusCode	= 206 // RFC 7233, 4.1
	StatusMultiStatus          	HttpStatusCode	= 207 // RFC 4918, 11.1
	StatusAlreadyReported      	HttpStatusCode	= 208 // RFC 5842, 7.1
	StatusIMUsed               	HttpStatusCode	= 226 // RFC 3229, 10.4.1

	StatusMultipleChoices   		HttpStatusCode	= 300 // RFC 7231, 6.4.1
	StatusMovedPermanently  		HttpStatusCode	= 301 // RFC 7231, 6.4.2
	StatusFound             		HttpStatusCode	= 302 // RFC 7231, 6.4.3
	StatusSeeOther          		HttpStatusCode	= 303 // RFC 7231, 6.4.4
	StatusNotModified       		HttpStatusCode	= 304 // RFC 7232, 4.1
	StatusUseProxy          		HttpStatusCode	= 305 // RFC 7231, 6.4.5
	_                       		HttpStatusCode	= 306 // RFC 7231, 6.4.6 (Unused)
	StatusTemporaryRedirect 		HttpStatusCode	= 307 // RFC 7231, 6.4.7
	StatusPermanentRedirect 		HttpStatusCode	= 308 // RFC 7538, 3

	StatusBadRequest				HttpStatusCode	= 400 // RFC 7231, 6.5.1
	StatusUnauthorized           	HttpStatusCode	= 401 // RFC 7235, 3.1
	StatusPaymentRequired		HttpStatusCode	= 402 // RFC 7231, 6.5.2
	StatusForbidden				HttpStatusCode	= 403 // RFC 7231, 6.5.3
	StatusNotFound				HttpStatusCode	= 404 // RFC 7231, 6.5.4
	StatusMethodNotAllowed		HttpStatusCode	= 405 // RFC 7231, 6.5.5
	StatusNotAcceptable			HttpStatusCode	= 406 // RFC 7231, 6.5.6
	StatusProxyAuthRequired		HttpStatusCode	= 407 // RFC 7235, 3.2
	StatusRequestTimeout			HttpStatusCode	= 408 // RFC 7231, 6.5.7
	StatusConflict				HttpStatusCode	= 409 // RFC 7231, 6.5.8
	StatusGone						HttpStatusCode	= 410 // RFC 7231, 6.5.9
	StatusLengthRequired 			HttpStatusCode	= 411 // RFC 7231, 6.5.10
	StatusPreconditionFailed 	HttpStatusCode	= 412 // RFC 7232, 4.2
	StatusRequestEntityTooLarge	HttpStatusCode	= 413 // RFC 7231, 6.5.11
	StatusRequestURITooLong 		HttpStatusCode	= 414 // RFC 7231, 6.5.12
	StatusUnsupportedMediaType	HttpStatusCode	= 415 // RFC 7231, 6.5.13
	StatusExpectationFailed		HttpStatusCode	= 417 // RFC 7231, 6.5.14
	StatusTeapot					HttpStatusCode	= 418 // RFC 7168, 2.3.3
	StatusMisdirectedRequest 	HttpStatusCode	= 421 // RFC 7540, 9.1.2
	StatusUnprocessableEntity	HttpStatusCode	= 422 // RFC 4918, 11.2
	StatusLocked					HttpStatusCode	= 423 // RFC 4918, 11.3
	StatusFailedDependency		HttpStatusCode	= 424 // RFC 4918, 11.4
	StatusUpgradeRequired		HttpStatusCode	= 426 // RFC 7231, 6.5.15
	StatusPreconditionRequired	HttpStatusCode	= 428 // RFC 6585, 3
	StatusTooManyRequests		HttpStatusCode	= 429 // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge	HttpStatusCode	= 431 // RFC 6585, 5
	StatusUnavailableForLegalReasons	HttpStatusCode	= 451 // RFC 7725, 3

	StatusInternalServerError	HttpStatusCode	= 500 // RFC 7231, 6.6.1
	StatusNotImplemented			HttpStatusCode	= 501 // RFC 7231, 6.6.2
	StatusBadGateway				HttpStatusCode	= 502 // RFC 7231, 6.6.3
	StatusServiceUnavailable		HttpStatusCode	= 503 // RFC 7231, 6.6.4
	StatusGatewayTimeout			HttpStatusCode	= 504 // RFC 7231, 6.6.5
	StatusHTTPVersionNotSupported	HttpStatusCode	= 505 // RFC 7231, 6.6.6
	StatusVariantAlsoNegotiates		HttpStatusCode	= 506 // RFC 2295, 8.1
	StatusInsufficientStorage		HttpStatusCode	= 507 // RFC 4918, 11.5
	StatusLoopDetected				HttpStatusCode	= 508 // RFC 5842, 7.2
	StatusNotExtended					HttpStatusCode	= 510 // RFC 2774, 7
	StatusNetworkAuthenticationRequired HttpStatusCode	= 511 // RFC 6585, 6
)