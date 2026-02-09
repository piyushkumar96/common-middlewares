package context

import (
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/piyushkumar96/app-monitoring/interfaces"
)

type CtxMeta struct {
	DeploymentID string
	UserID       string
	TraceParent  string
	TraceState   string
	TraceID      string
	ReqID        string
	Path         string
	Body         string
	UA           string
	QP           url.Values
	Time         time.Time
}

// ResponseMeta holds response writer and optional app metrics from piyushkumar96/app-monitoring.
type ResponseMeta struct {
	ResWriter  gin.ResponseWriter
	AppMetrics interfaces.AppMetricsInterface
}

type TraceMeta struct {
	Trace              []string
	Error              []string
	IdentifierMappings map[string]interface{}
}
