// Package middlewares hosts all components which provides a convenient mechanism for filtering HTTP requests
// Cors defines the Cors policy of the application
package cors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	l "github.com/piyushkumar96/generic-logger"
)

type CORSHeaders struct {
	ContentType                   string
	AccessControlAllowOrigin      string
	AccessControlMaxAge           string
	AccessControlAllowMethods     string
	AccessControlAllowHeaders     string
	AccessControlAllowCredentials string
}

// CORS will handle the CORS middleware
func CORS(corsHeaders *CORSHeaders) gin.HandlerFunc {
	return func(c *gin.Context) {
		l.Logger.Debug("Middleware > Cors() logic starts")
		matched, err := MatchStringWithRegex(corsHeaders.AccessControlAllowOrigin, c.GetHeader("Origin"))
		if err != nil {
			l.Logger.Fatal("invalid allow origin regex pattern", "err", err.Error())
		}
		if !matched {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Header("Content-Type", "application/json")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Max-Age", "86400")
		c.Header("Access-Control-Allow-Methods", corsHeaders.AccessControlAllowMethods)
		c.Header("Access-Control-Allow-Headers", corsHeaders.AccessControlAllowHeaders)
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		} else {
			l.Logger.Debug("Middleware > Cors() logic ends")
			c.Next()
		}
	}
}
