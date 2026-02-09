package trace

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	im "github.com/piyushkumar96/app-monitoring/interfaces"
	cx "github.com/piyushkumar96/common-middlewares/context"
)

// Trace will generate req-id middleware
func Trace(appMetrics im.AppMetricsInterface) gin.HandlerFunc {
	return func(gc *gin.Context) {
		if gc.Request.Method == http.MethodOptions {
			gc.AbortWithStatus(http.StatusNoContent)
		} else {
			ctx := context.Background()

			/** initialize context meta */
			key, cm := cx.InitContextMeta(gc, "")
			ctx = context.WithValue(ctx, key, cm)

			/** initialize response meta */
			key, rm := cx.InitResponseMeta(gc, appMetrics)
			ctx = context.WithValue(ctx, key, rm)

			/** initialize trace meta */
			key, tm := cx.InitTraceMeta()
			ctx = context.WithValue(ctx, key, tm)

			/** Setting the context in GIN context */
			gc.Set(cx.CtxKey, ctx)
			gc.Next()
		}
	}
}
