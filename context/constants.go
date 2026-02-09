package context

const (
	CtxKey          = "ctx"
	CtxMetaKey      = "ctxMeta"
	ResponseMetaKey = "ResponseMeta"
	TraceMetaKey    = "TraceMeta"
	ReqIDKey        = "request_id"
)

// Trace separator used in AddTrace
const UnderScore = "_"

type TRequestHeaderKey string

const (
	HeaderAuthorization TRequestHeaderKey = "Authorization"
	HeaderResponseReqID TRequestHeaderKey = "req-id"
	HeaderContentType   TRequestHeaderKey = "Content-Type"
	HeaderAccountID     TRequestHeaderKey = "x-account-id"
	HeaderUserIDKey     TRequestHeaderKey = "x-user-id"
	HeaderAPIKey        TRequestHeaderKey = "x-api-key"
)

type TResponseContentType string

const (
	ApplicationJSON TResponseContentType = "application/json"
)
