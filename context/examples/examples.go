// Package main demonstrates the context package: request ID, context meta, and response helpers.
// Run: go run github.com/piyushkumar96/common-middlewares/context/examples
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/piyushkumar96/common-middlewares/context"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		context.InitRequestContext(c)
		c.Next()
	})

	r.GET("/ping", func(c *gin.Context) {
		ctx := context.GetRequestContext(c)
		reqID := context.GetRequestID(ctx)
		c.JSON(http.StatusOK, gin.H{"message": "pong", "request_id": reqID})
	})

	r.GET("/fail", func(c *gin.Context) {
		context.RespondJSON(c, http.StatusBadRequest, context.MessageFailure("something went wrong"))
	})

	fmt.Println("Context example: GET http://localhost:8083/ping or /fail")
	if err := r.Run(":8083"); err != nil {
		log.Fatal(err)
	}
}
