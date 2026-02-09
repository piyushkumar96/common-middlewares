// Package main demonstrates the CORS middleware with configurable headers and origin regex.
// Run: go run github.com/piyushkumar96/common-middlewares/cors/examples
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/piyushkumar96/common-middlewares/cors"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Use regex for AccessControlAllowOrigin; ".*" allows any origin.
	headers := &cors.CORSHeaders{
		ContentType:                   "application/json",
		AccessControlAllowOrigin:      ".*",
		AccessControlMaxAge:           "86400",
		AccessControlAllowMethods:     "POST, GET, PUT, DELETE, UPDATE",
		AccessControlAllowHeaders:     "Content-Type, Authorization, X-Request-ID",
		AccessControlAllowCredentials: "true",
	}
	r.Use(cors.CORS(headers))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	fmt.Println("CORS example: GET http://localhost:8081/ping")
	if err := r.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}
