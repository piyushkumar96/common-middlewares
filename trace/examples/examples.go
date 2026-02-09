// Package main demonstrates the trace middleware: initializes request context (context meta, response meta, trace meta) and sets request ID.
// Run: go run github.com/piyushkumar96/common-middlewares/trace/examples
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/piyushkumar96/common-middlewares/context"
	"github.com/piyushkumar96/common-middlewares/trace"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Pass nil for appMetrics to skip metrics; or use app-monitoring's NewPromAppMetrics / NewNoOpPromAppMetrics.
	r.Use(trace.Trace(nil))

	r.GET("/ping", func(c *gin.Context) {
		ctx := context.GetRequestContext(c)
		meta := context.GetContextMeta(ctx)
		reqID := meta.ReqID
		c.JSON(http.StatusOK, gin.H{"message": "pong", "request_id": reqID})
	})

	fmt.Println("Trace example: GET http://localhost:8080/ping")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
