package gecko

// ================================================================ //
// ========== CONSTS ============================================== //

// Status codes
const (
	Status200OK                   = 200
	Status201Created              = 201
	Status202Accepted             = 202
	Status203NonAuthoritativeInfo = 203
	Status204NoContent            = 204
	Status205ResetContent         = 205
	Status206PartialContent       = 206
	Status207MultiStatus          = 207
	Status208AlreadyReported      = 208

	Status300MultipleChoices   = 300
	Status301MovedPermanently  = 301
	Status302Found             = 302
	Status303SeeOther          = 303
	Status304NotModified       = 304
	Status305UseProxy          = 305
	Status307TemporaryRedirect = 307
	Status308PermanentRedirect = 308

	Status400BadRequest                   = 400
	Status401Unauthorized                 = 401
	Status402PaymentRequired              = 402
	Status403Forbidden                    = 403
	Status404NotFound                     = 404
	Status405MethodNotAllowed             = 405
	Status406NotAcceptable                = 406
	Status407ProxyAuthRequired            = 407
	Status408RequestTimeout               = 408
	Status409Conflict                     = 409
	Status410Gone                         = 410
	Status411LengthRequired               = 411
	Status412PreconditionFailed           = 412
	Status413RequestEntityTooLarge        = 413
	Status414RequestURITooLong            = 414
	Status415UnsupportedMediaType         = 415
	Status416RequestedRangeNotSatisfiable = 416
	Status417ExpectationFailed            = 417
	Status418Teapot                       = 418
	Status421MisdirectedRequest           = 421
	Status422UnprocessableEntity          = 422
	Status423Locked                       = 423
	Status424FailedDependency             = 424
	Status425TooEarly                     = 425
	Status426UpgradeRequired              = 426
	Status428PreconditionRequired         = 428
	Status429TooManyRequests              = 429
	Status431RequestHeaderFieldsTooLarge  = 431
	Status451UnavailableForLegalReasons   = 451

	Status500InternalServerError           = 500
	Status501NotImplemented                = 501
	Status502BadGateway                    = 502
	Status503ServiceUnavailable            = 503
	Status504GatewayTimeout                = 504
	Status505HTTPVersionNotSupported       = 505
	Status506VariantAlsoNegotiates         = 506
	Status507InsufficientStorage           = 507
	Status508LoopDetected                  = 508
	Status510NotExtended                   = 510
	Status511NetworkAuthenticationRequired = 511
)

// MIME types
const (
	charsetUTF8 = "charset=UTF-8"

	MIMEApplicationJSON                  = "application/json"
	MIMEApplicationJSONCharsetUTF8       = MIMEApplicationJSON + "; " + charsetUTF8
	MIMEApplicationJavaScript            = "application/javascript"
	MIMEApplicationJavaScriptCharsetUTF8 = MIMEApplicationJavaScript + "; " + charsetUTF8
	MIMEApplicationXML                   = "application/xml"
	MIMEApplicationXMLCharsetUTF8        = MIMEApplicationXML + "; " + charsetUTF8
	MIMETextXML                          = "text/xml"
	MIMETextXMLCharsetUTF8               = MIMETextXML + "; " + charsetUTF8
	MIMEApplicationForm                  = "application/x-www-form-urlencoded"
	MIMEApplicationProtobuf              = "application/protobuf"
	MIMEApplicationMsgpack               = "application/msgpack"
	MIMETextHTML                         = "text/html"
	MIMETextHTMLCharsetUTF8              = MIMETextHTML + "; " + charsetUTF8
	MIMETextPlain                        = "text/plain"
	MIMETextPlainCharsetUTF8             = MIMETextPlain + "; " + charsetUTF8
	MIMEMultipartForm                    = "multipart/form-data"
	MIMEOctetStream                      = "application/octet-stream"
)

// Headers
const (
	HeaderAccept         = "Accept"
	HeaderAcceptEncoding = "Accept-Encoding"
	// HeaderAllow is the name of the "Allow" header field used to list the set of methods
	// advertised as supported by the target resource. Returning an Allow header is mandatory
	// for status 405 (method not found) and useful for the OPTIONS method in responses.
	// See RFC 7231: https://datatracker.ietf.org/doc/html/rfc7231#section-7.4.1
	HeaderAllow               = "Allow"
	HeaderAuthorization       = "Authorization"
	HeaderContentDisposition  = "Content-Disposition"
	HeaderContentEncoding     = "Content-Encoding"
	HeaderContentLength       = "Content-Length"
	HeaderContentType         = "Content-Type"
	HeaderCookie              = "Cookie"
	HeaderSetCookie           = "Set-Cookie"
	HeaderIfModifiedSince     = "If-Modified-Since"
	HeaderLastModified        = "Last-Modified"
	HeaderLocation            = "Location"
	HeaderRetryAfter          = "Retry-After"
	HeaderUpgrade             = "Upgrade"
	HeaderVary                = "Vary"
	HeaderWWWAuthenticate     = "WWW-Authenticate"
	HeaderXForwardedFor       = "X-Forwarded-For"
	HeaderXForwardedProto     = "X-Forwarded-Proto"
	HeaderXForwardedProtocol  = "X-Forwarded-Protocol"
	HeaderXForwardedSsl       = "X-Forwarded-Ssl"
	HeaderXUrlScheme          = "X-Url-Scheme"
	HeaderXHTTPMethodOverride = "X-HTTP-Method-Override"
	HeaderXRealIP             = "X-Real-Ip"
	HeaderXRequestID          = "X-Request-Id"
	HeaderXCorrelationID      = "X-Correlation-Id"
	HeaderXRequestedWith      = "X-Requested-With"
	HeaderServer              = "Server"
	HeaderOrigin              = "Origin"
	HeaderCacheControl        = "Cache-Control"
	HeaderConnection          = "Connection"

	// Access control
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"

	// Security
	HeaderStrictTransportSecurity         = "Strict-Transport-Security"
	HeaderXContentTypeOptions             = "X-Content-Type-Options"
	HeaderXXSSProtection                  = "X-XSS-Protection"
	HeaderXFrameOptions                   = "X-Frame-Options"
	HeaderContentSecurityPolicy           = "Content-Security-Policy"
	HeaderContentSecurityPolicyReportOnly = "Content-Security-Policy-Report-Only"
	HeaderXCSRFToken                      = "X-CSRF-Token"
	HeaderReferrerPolicy                  = "Referrer-Policy"
)

const (
	HxLocation           = "HX-Location"             // Response header: allows you to do a client-side redirect that does not do a full page reload.
	HxPushUrl            = "HX-Push-Url"             // Response header: pushes a new url into the history stack.
	HxRedirect           = "HX-Redirect"             // Response header: can be used to do a client-side redirect to a new location.
	HxRefresh            = "HX-Refresh"              // Response header: if set to “true” the client-side will do a full refresh of the page.
	HxReplaceUrl         = "HX-Replace-Url"          // Response header: replaces the current URL in the location bar.
	HxReswap             = "HX-Reswap"               // Response header: allows you to specify how the response will be swapped. See hx-swap for possible values.
	HxRetarget           = "HX-Retarget"             // Response header: a CSS selector that updates the target of the content update to a different element on the page.
	HxReselect           = "HX-Reselect"             // Response header: a CSS selector that allows you to choose which part of the response is used to be swapped in. Overrides an existing hx-select on the triggering element.
	HxTriggerRes         = "HX-Trigger"              // Response header: allows you to trigger client-side events.
	HxTriggerAfterSettle = "HX-Trigger-After-Settle" // Response header: allows you to trigger client-side events after the settle step.
	HxTriggerAfterSwap   = "HX-Trigger-After-Swap"   // Response header: allows you to trigger client-side events after the swap step.

	HxBoosted               = "HX-Boosted"                 // Request header: indicates that the request is via an element using hx-boost
	HxCurrentUrl            = "HX-Current-URL"             // Request header: the current URL of the browser
	HxHistoryRestoreRequest = "HX-History-Restore-Request" // Request header: “true” if the request is for history restoration after a miss in the local history cache
	HxPrompt                = "HX-Prompt"                  // Request header: the user response to an hx-prompt
	HxRequest               = "HX-Request"                 // Request header: always “true”
	HxTarget                = "HX-Target"                  // Request header: the id of the target element if it exists
	HxTriggerName           = "HX-Trigger-Name"            // Request header: the name of the triggered element if it exists
	HxTriggerReq            = "HX-Trigger"                 // Request header: the id of the triggered element if it exists

	HxPromptEncoded = "HX-Prompt-Encoded" // Request header: (gecko) la respuesta url encoded a un hx-prompt
	HxAskfor        = "Hx-Askfor"         // Request header: (gecko) lo que el cliente solicita como respuesta mediante hx-askfor
)
