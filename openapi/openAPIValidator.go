package openapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"unsafe"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/legacy"
	"github.com/gin-gonic/gin"
	ae "github.com/piyushkumar96/app-error"
	cx "github.com/piyushkumar96/common-middlewares/context"
	l "github.com/piyushkumar96/generic-logger"
)

type SwaggerConfig struct {
	SwaggerFilePath string
	SwaggerFileName string
}

// loadOpenAPISpecs loading the Open API Specs
func loadOpenAPISpecs(filePath string) (*openapi3.T, error) {
	loader := openapi3.NewLoader()
	return loader.LoadFromFile(filePath)
}

func NewOpenAPIValidatorRouter(config *SwaggerConfig) (*routers.Router, error) {
	openAPIConfig, err := loadOpenAPISpecs(config.SwaggerFilePath + config.SwaggerFileName)
	if err != nil {
		if l.Logger != nil {
			l.Logger.Fatal(ErrLoadSpec.Message, "err", err.Error())
		}
		return nil, err
	}
	router, err := legacy.NewRouter(openAPIConfig)
	if err != nil {
		if l.Logger != nil {
			l.Logger.Fatal(ErrCreateRouter.Message, "err", err.Error())
		}
		return nil, err
	}
	return &router, nil
}

func OpenAPIValidator(router routers.Router) gin.HandlerFunc {
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
		gc.Next()
	}
}

func validateRequest(gc *gin.Context, pathParams map[string]string, route *routers.Route) error {
	// Validate the request against the OpenAPI specification
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    gc.Request,
		PathParams: pathParams,
		Route:      route,
		Options: &openapi3filter.Options{
			MultiError: true,
		},
	}
	validationErr := openapi3filter.ValidateRequest(gc.Request.Context(), requestValidationInput)
	return validationErr
}

// buildFinalValidationErrorMsg building the final request validation error msg basically combine all
func buildFinalValidationErrorMsg(ctx context.Context, validationMultiError *openapi3.MultiError) string {
	errsMsg := make([]string, 0)
	for _, vme := range *validationMultiError {
		reqValidationErr, ok := vme.(*openapi3filter.RequestError)
		if !ok {
			continue
		}
		multiErr, multiErrOK := reqValidationErr.Err.(openapi3.MultiError)
		if multiErrOK {
			errsMsg = append(errsMsg, handleMultiError(reqValidationErr, multiErr)...)
		}

		// Handling error related to the key data type
		parseErr, parseErrOK := reqValidationErr.Err.(*openapi3filter.ParseError)
		if parseErrOK {
			if reqValidationErr.Parameter != nil {
				errsMsg = append(errsMsg, fmt.Sprintf("%s-%s%s %s", reqValidationErr.Parameter.In, paramsKey, reqValidationErr.Parameter.Name, parseErr.Reason))
			} else {
				errsMsg = append(errsMsg, fmt.Sprintf("%s%s", requestBodyKey, parseErr.Reason))
			}
		}

		// Handling error related to unexpected key send in request body
		schemaErr, schemaErrOK := reqValidationErr.Err.(*openapi3.SchemaError)
		if schemaErrOK {
			originMultiErr, originMultiErrOK := schemaErr.Origin.(openapi3.MultiError)
			if originMultiErrOK {
				errsMsg = append(errsMsg, handleMultiError(reqValidationErr, originMultiErr)...)
			} else if reqValidationErr.Parameter != nil {
				// check for schema error in case with Origin
				errsMsg = append(errsMsg, fmt.Sprintf("%s-%s %s %s", reqValidationErr.Parameter.In, paramsKey, reqValidationErr.Parameter.Name, schemaErr.Reason))
				continue
			} else {
				var wrapError reflect.Value
				// Unwrapping the nested SchemaError to wrapError
				nestedSErr, nestedSErrOK := schemaErr.Origin.(*openapi3.SchemaError)
				if nestedSErrOK {
					// Loop to find the innermost SchemaError
					for nestedSErrOK {
						_, nestedSErrOK = nestedSErr.Origin.(*openapi3.SchemaError)
					}
					// Set wrapError to the innermost Origin
					wrapError = reflect.ValueOf(nestedSErr.Origin)
				} else {
					// If not a nested SchemaError, set wrapError to the original Origin
					wrapError = reflect.ValueOf(schemaErr.Origin)
				}

				// Handling the case of schema miss match error for anyOf, oneOf
				if wrapError.Kind() == reflect.Ptr {
					wrapError = wrapError.Elem()
				}
				errField := wrapError.FieldByName("err")
				if errField.IsValid() {
					dmMErrs, dmMErrOK := convertToMultiErrorSlice(getUnexportedField(errField))
					if dmMErrOK {
						for _, dmMErr := range dmMErrs {
							errsMsg = append(errsMsg, handleMultiError(reqValidationErr, dmMErr)...)
						}
					}
				} else {
					msg := "could not access the unexported variable err under schemaErr.Origin"
					cx.AddTraceLog(ctx, msg)
					if l.Logger != nil {
						l.Logger.Debug(msg, "meta", map[string]interface{}{"open_api_error": schemaErr.Origin.Error()})
					}
				}
				continue
			}
		}

		if !multiErrOK && !parseErrOK {
			if reqValidationErr.Parameter != nil {
				errsMsg = append(errsMsg, fmt.Sprintf("%s-%s %s %s", reqValidationErr.Parameter.In, paramsKey, reqValidationErr.Parameter.Name, reqValidationErr.Reason))
			}
		}

	}

	errsMsgStr := ""
	if len(errsMsg) > 0 {
		// Process validationErrs to get all errors as a single string
		errsMsgStr = strings.ReplaceAll(strings.Join(errsMsg, ", "), `"`, "")
	}
	return errsMsgStr
}

func handleMultiError(reqValidationErr *openapi3filter.RequestError, multiErr openapi3.MultiError) []string {
	errsMsg := make([]string, 0)
	for _, ome := range multiErr {
		omeSchemaErr, omeSchemaErrOK := ome.(*openapi3.SchemaError)
		if omeSchemaErrOK {
			// to filter the property name as it is not coming in error reason
			errsMsg = append(errsMsg, formOpenAPISchemaErrorMessage(reqValidationErr, omeSchemaErr)...)
		}
	}
	return errsMsg
}

func formOpenAPISchemaErrorMessage(reqValidationErr *openapi3filter.RequestError, dmmeSErr *openapi3.SchemaError) []string {
	errsMsg := make([]string, 0)
	if reqValidationErr.Parameter != nil {
		errsMsg = append(errsMsg, fmt.Sprintf("%s-%s%s %s", reqValidationErr.Parameter.In, paramsKey, reqValidationErr.Parameter.Name, dmmeSErr.Reason))
	} else if reflect.TypeOf(dmmeSErr.Origin) == reflect.TypeOf(openapi3.MultiError{}) {
		errsMsg = append(errsMsg, getErrorCustomReasonForSchemaErrorInReqBody(dmmeSErr))
		errsMsg = append(errsMsg, handleMultiError(reqValidationErr, dmmeSErr.Origin.(openapi3.MultiError))...)
	} else {
		errsMsg = append(errsMsg, getErrorCustomReasonForSchemaErrorInReqBody(dmmeSErr))
	}
	return errsMsg
}

func getErrorCustomReasonForSchemaErrorInReqBody(schemaErr *openapi3.SchemaError) string {
	// to filter the property name as it is not coming in error reason
	reversePath := schemaErr.JSONPointer()
	propKey := ""
	switch len(reversePath) {
	case 1:
		propKey = reversePath[0]
	case 2:
		propKey = fmt.Sprintf("%v[%v]", reversePath[0], reversePath[1])
	case 3:
		propKey = fmt.Sprintf("%v[%v][%v]", reversePath[0], reversePath[1], reversePath[2])
	case 4:
		propKey = fmt.Sprintf("%v[%v][%v][%v]", reversePath[0], reversePath[1], reversePath[2], reversePath[3])
	}
	return fmt.Sprintf("%s%s %s", requestBodyKey, propKey, schemaErr.Reason)
}

func getUnexportedField(field reflect.Value) interface{} {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}

func convertToMultiErrorSlice(data interface{}) ([]openapi3.MultiError, bool) {
	// Use reflection to access the underlying slice
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice {
		// Get the length of the slice
		length := v.Len()
		// Create a slice to hold the converted MultiError elements
		result := make([]openapi3.MultiError, length)
		for i := 0; i < length; i++ {
			elem := v.Index(i).Interface()
			// Check if the element is a SchemaError and unwrap it if necessary
			nestedSErr, nestedSErrOK := elem.(*openapi3.SchemaError)
			for nestedSErrOK {
				elem = nestedSErr.Origin
				nestedSErr, nestedSErrOK = nestedSErr.Origin.(*openapi3.SchemaError)
			}
			// Check if the unwrapped element is a MultiError
			if castedElem, ok := elem.(openapi3.MultiError); ok {
				result[i] = castedElem
			} else {
				// Return false if any element is not a MultiError
				return nil, false
			}
		}
		return result, true
	}
	return nil, false
}
