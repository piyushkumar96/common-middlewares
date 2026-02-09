package authentication

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	ae "github.com/piyushkumar96/app-error"
	"github.com/piyushkumar96/common-middlewares/context"
)

// AuthConfig configures the Auth middleware
type AuthConfig struct {
	Token string
}

// Auth returns a gin middleware that validates the request against AuthConfig.Token (e.g. Bearer or static token in Authorization header).
func Auth(authConfig *AuthConfig) gin.HandlerFunc {
	return func(gc *gin.Context) {
		headers := gc.Request.Header[string(context.HeaderAuthorization)]
		if len(headers) == 0 || authConfig.Token != headers[0] {
			ctx := context.GetRequestContext(gc)
			appErr := ae.GetAppErr(ctx, errors.New(ErrUnauthorized.Message), ErrUnauthorized, http.StatusUnauthorized)
			context.RespondJSON(gc, http.StatusUnauthorized, context.MessageFailure(appErr.GetMsg()))
			gc.Abort()
			return
		}
		gc.Next()
	}
}
