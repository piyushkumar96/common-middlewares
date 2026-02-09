package openapi

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/legacy"
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
