// Package main demonstrates the OpenAPI validator middleware.
// Requires an OpenAPI spec file (e.g. openapi.yaml or openapi.json) in the current directory or adjust paths.
// Run from repo root: go run github.com/piyushkumar96/common-middlewares/openapi/examples
// Or with a spec: openapi/examples/openapi.yaml
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/piyushkumar96/common-middlewares/context"
	"github.com/piyushkumar96/common-middlewares/openapi"
)

func main() {
	// Example: look for spec in examples dir or current dir
	dir := "."
	if p, err := os.Getwd(); err == nil {
		dir = p
	}
	specPath := filepath.Join(dir, "openapi.yaml")
	if _, err := os.Stat(specPath); err != nil {
		specPath = filepath.Join(dir, "openapi.json")
	}
	if _, err := os.Stat(specPath); err != nil {
		log.Printf("No openapi.yaml or openapi.json found in %s; create one to validate requests", dir)
		log.Printf("Skipping validator; running server without OpenAPI validation")
		runServer(nil)
		return
	}

	config := &openapi.SwaggerConfig{
		SwaggerFilePath: filepath.Dir(specPath) + "/",
		SwaggerFileName: filepath.Base(specPath),
	}
	router, err := openapi.NewOpenAPIValidatorRouter(config)
	if err != nil {
		log.Fatalf("OpenAPI router: %v", err)
	}
	// OpenAPIValidator expects routers.Router interface; router is *routers.Router.
	runServer(openapi.OpenAPIValidator(*router))
}

func runServer(validator gin.HandlerFunc) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(func(c *gin.Context) { context.InitRequestContext(c); c.Next() })
	if validator != nil {
		r.Use(validator)
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	fmt.Println("OpenAPI example: GET http://localhost:8084/ping")
	if err := r.Run(":8084"); err != nil {
		log.Fatal(err)
	}
}
