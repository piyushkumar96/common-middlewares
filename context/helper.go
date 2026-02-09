package context

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/piyushkumar96/app-monitoring/interfaces"
)

func newRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b) + fmt.Sprintf("%d", time.Now().UnixNano())
}

// InitRequestContext creates a context with request ID and stores it in gin.
// Use GetRequestContext to retrieve it. Call this from Trace middleware.
func InitRequestContext(gc *gin.Context) context.Context {
	reqID := gc.GetHeader("X-Request-ID")
	if reqID == "" {
		reqID = newRequestID()
	}
	gc.Header("X-Request-ID", reqID)
	ctx := gc.Request.Context()
	ctx = context.WithValue(ctx, ReqIDKey, reqID)
	gc.Set(CtxKey, ctx)
	return ctx
}

// GetRequestContext returns the request context from gin, or the request's context.
func GetRequestContext(gc *gin.Context) context.Context {
	if v, ok := gc.Get(CtxKey); ok {
		if c, ok := v.(context.Context); ok {
			return c
		}
	}
	return gc.Request.Context()
}

func InitContextMeta(gc *gin.Context, body string) (string, *CtxMeta) {
	reqID := gc.GetHeader("X-Request-ID")
	if reqID == "" {
		reqID = newRequestID()
		gc.Header("X-Request-ID", reqID)
	}
	ctxMeta := CtxMeta{
		DeploymentID: gc.GetHeader("x-deployment-id"),
		UserID:       gc.GetHeader("x-user-id"),
		TraceParent:  gc.GetHeader("traceparent"),
		TraceState:   gc.GetHeader("tracestate"),
		TraceID:      fmt.Sprintf("%s:%s", gc.GetHeader("traceparent"), gc.GetHeader("tracestate")),
		ReqID:        reqID,
		Path:         gc.Request.RequestURI,
		Body:         body,
		UA:           gc.GetHeader("User-Agent"),
		QP:           gc.Request.URL.Query(),
		Time:         time.Now(),
	}
	return CtxMetaKey, &ctxMeta
}

func InitResponseMeta(gc *gin.Context, appMetrics interfaces.AppMetricsInterface) (string, *ResponseMeta) {
	return ResponseMetaKey, &ResponseMeta{ResWriter: gc.Writer, AppMetrics: appMetrics}
}

func InitTraceMeta() (string, *TraceMeta) {
	return TraceMetaKey, &TraceMeta{Trace: make([]string, 0)}
}

// GetContext returns the context from gin (with optional body update to CtxMeta). If not set, returns Background and sets it.
func GetContext(gc *gin.Context, body string) context.Context {
	var ctx context.Context
	ctxInterface, exists := gc.Get(CtxKey)
	if exists {
		ctx = ctxInterface.(context.Context)
		ctxMeta := GetContextMeta(ctx)
		ctxMeta.Body = body
	} else {
		ctx = context.Background()
		gc.Set(CtxKey, ctx)
	}
	return ctx
}

func GetContextMeta(ctx context.Context) *CtxMeta {
	meta := ctx.Value(CtxMetaKey)
	ctxMeta, ok := meta.(*CtxMeta)
	if !ok {
		return &CtxMeta{}
	}
	return ctxMeta
}

func GetResponseMeta(ctx context.Context) *ResponseMeta {
	meta := ctx.Value(ResponseMetaKey)
	respMeta, ok := meta.(*ResponseMeta)
	if !ok {
		return &ResponseMeta{}
	}
	return respMeta
}

func GetTraceMeta(ctx context.Context) *TraceMeta {
	trace := ctx.Value(TraceMetaKey)
	traceMeta, ok := trace.(*TraceMeta)
	if !ok {
		return &TraceMeta{}
	}
	return traceMeta
}

func AddTrace(ctx context.Context, msg ...string) *TraceMeta {
	if ctx == nil {
		return nil
	}
	trace := ctx.Value(TraceMetaKey)
	traceMeta, ok := trace.(*TraceMeta)
	if !ok {
		return nil
	}
	traceMeta.Trace = append(traceMeta.Trace, strings.Join(msg, UnderScore))
	return traceMeta
}

func AddTraceLog(ctx context.Context, errorMsg string) *TraceMeta {
	if ctx == nil {
		return nil
	}
	trace := ctx.Value(TraceMetaKey)
	traceMeta, ok := trace.(*TraceMeta)
	if !ok {
		return nil
	}
	traceMeta.Error = append(traceMeta.Error, errorMsg)
	return traceMeta
}

// RespondJSON writes a JSON response and aborts with the given status code.
func RespondJSON(gc *gin.Context, statusCode int, body interface{}) {
	gc.AbortWithStatusJSON(statusCode, body)
}

// MessageFailure returns a simple failure map for JSON responses.
func MessageFailure(message string) map[string]interface{} {
	return map[string]interface{}{
		"success": false,
		"message": message,
	}
}

// GetRequestID returns the request ID from context or empty string.
func GetRequestID(ctx context.Context) string {
	if v := ctx.Value(ReqIDKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
