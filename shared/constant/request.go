package constant

type ctxKey string

const (
	AuthorizationHeader           = "Authorization"
	XRequestIDHeader              = "X-REQUEST-ID"
	XRequestIDHeaderCtxKey ctxKey = "X-REQUEST-ID"
)
