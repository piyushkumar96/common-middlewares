// Package main demonstrates the authentication middleware (static token in Authorization header).
// Run: go run github.com/piyushkumar96/common-middlewares/authentication/examples
// Then: curl -H "Authorization: my-secret-token" http://localhost:8082/ping
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/piyushkumar96/common-middlewares/authentication"
	"github.com/piyushkumar96/common-middlewares/context"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Optional: set request ID (wrap so it's a gin.HandlerFunc)
	r.Use(func(c *gin.Context) { context.InitRequestContext(c); c.Next() })
	r.Use(authentication.Auth(&authentication.AuthConfig{Token: "my-secret-token"}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	fmt.Println("Auth example: curl -H \"Authorization: my-secret-token\" http://localhost:8082/ping")
	if err := r.Run(":8082"); err != nil {
		log.Fatal(err)
	}
}
