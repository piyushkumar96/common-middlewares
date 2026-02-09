package openapi

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/gin-gonic/gin"
	ae "github.com/piyushkumar96/app-error"
	cx "github.com/piyushkumar96/common-middlewares/context"
	l "github.com/piyushkumar96/generic-logger"
)

// responseRecorder captures status code and body written by the handler for response validation.
type responseRecorder struct {
	gin.ResponseWriter
	status int
	body   *bytes.Buffer
}

func newResponseRecorder(w gin.ResponseWriter) *responseRecorder {
	return &responseRecorder{
		ResponseWriter: w,
		status:         http.StatusOK,
		body:           &bytes.Buffer{},
	}
}

func (r *responseRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r *responseRecorder) Header() http.Header {
	return r.ResponseWriter.Header()
}

// OpenAPIValidatorRequestAndResponse runs request validation, then captures and validates the response against the OpenAPI spec.
// Response validation failures are logged only (response is already sent to the client).
func OpenAPIValidatorRequestAndResponse(router routers.Router) gin.HandlerFunc {
	return func(gc *gin.Context) {
		ctx := cx.GetRequestContext(gc)
		route, pathParams, err := router.FindRoute(gc.Request)
		if err != nil {
			appErr := ae.GetAppErr(ctx, err, ErrRouteNotFound, http.StatusNotFound)
			cx.RespondJSON(gc, http.StatusNotFound, cx.MessageFailure(appErr.GetMsg()))
			gc.Abort()
			return
		}

		validationError := validateRequest(gc, pathParams, route)
		if validationError != nil {
			validationMultiError, ok := validationError.(openapi3.MultiError)
			if !ok {
				if l.Logger != nil {
					l.Logger.Error(ErrValidationUnexpectedType.Message, "meta", map[string]interface{}{"type": reflect.TypeOf(validationError)})
				}
				appErr := ae.GetAppErr(ctx, errors.New(ErrValidationUnexpectedType.Message), ErrValidationUnexpectedType, http.StatusInternalServerError)
				cx.RespondJSON(gc, http.StatusInternalServerError, cx.MessageFailure(appErr.GetMsg()))
				gc.Abort()
				return
			}
			finalValidationErrorMsg := buildFinalValidationErrorMsg(ctx, &validationMultiError)
			if finalValidationErrorMsg != "" {
				customErr := ae.GetCustomErr(ErrValidationFailed.Code, finalValidationErrorMsg, false)
				appErr := ae.GetAppErr(ctx, errors.New(customErr.Message), customErr, http.StatusBadRequest)
				cx.RespondJSON(gc, http.StatusBadRequest, cx.MessageFailure(appErr.GetMsg()))
				gc.Abort()
				return
			}
		}

		rec := newResponseRecorder(gc.Writer)
		gc.Writer = rec
		gc.Next()

		// Response validation: log only (response already sent)
		if errResp := validateResponse(gc.Request.Context(), route, pathParams, gc.Request, rec); errResp != nil {
			if l.Logger != nil {
				l.Logger.Error(ErrResponseValidationFailed.Message, "err", errResp.Error(), "path", gc.Request.URL.Path, "method", gc.Request.Method, "status", rec.status)
			}
		}
	}
}

// validateResponse checks the recorded response against the OpenAPI spec for the route.
func validateResponse(ctx context.Context, route *routers.Route, pathParams map[string]string, req *http.Request, rec *responseRecorder) error {
	reqInput := &openapi3filter.RequestValidationInput{
		Request:    req,
		PathParams: pathParams,
		Route:      route,
		Options: &openapi3filter.Options{
			MultiError: true,
		},
	}
	respInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: reqInput,
		Status:                 rec.status,
		Header:                 rec.Header().Clone(),
		Body:                   io.NopCloser(bytes.NewReader(rec.body.Bytes())),
		Options: &openapi3filter.Options{
			MultiError: true,
		},
	}
	return openapi3filter.ValidateResponse(ctx, respInput)
}
