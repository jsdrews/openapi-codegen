// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	strictgin "github.com/oapi-codegen/runtime/strictmiddleware/gin"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Server Health Check
	// (GET /health)
	GetHealth(c *gin.Context)
	// Get API Version
	// (GET /version)
	GetVersion(c *gin.Context)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// GetHealth operation middleware
func (siw *ServerInterfaceWrapper) GetHealth(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetHealth(c)
}

// GetVersion operation middleware
func (siw *ServerInterfaceWrapper) GetVersion(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetVersion(c)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.GET(options.BaseURL+"/health", wrapper.GetHealth)
	router.GET(options.BaseURL+"/version", wrapper.GetVersion)
}

type GetHealthRequestObject struct {
}

type GetHealthResponseObject interface {
	VisitGetHealthResponse(w http.ResponseWriter) error
}

type GetHealth200JSONResponse HealthCheck

func (response GetHealth200JSONResponse) VisitGetHealthResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetHealth500JSONResponse ServerError

func (response GetHealth500JSONResponse) VisitGetHealthResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type GetHealth503JSONResponse ServerError

func (response GetHealth503JSONResponse) VisitGetHealthResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(503)

	return json.NewEncoder(w).Encode(response)
}

type GetVersionRequestObject struct {
}

type GetVersionResponseObject interface {
	VisitGetVersionResponse(w http.ResponseWriter) error
}

type GetVersion200JSONResponse VersionCheck

func (response GetVersion200JSONResponse) VisitGetVersionResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetVersion500JSONResponse ServerError

func (response GetVersion500JSONResponse) VisitGetVersionResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Server Health Check
	// (GET /health)
	GetHealth(ctx context.Context, request GetHealthRequestObject) (GetHealthResponseObject, error)
	// Get API Version
	// (GET /version)
	GetVersion(ctx context.Context, request GetVersionRequestObject) (GetVersionResponseObject, error)
}

type StrictHandlerFunc = strictgin.StrictGinHandlerFunc
type StrictMiddlewareFunc = strictgin.StrictGinMiddlewareFunc

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
}

// GetHealth operation middleware
func (sh *strictHandler) GetHealth(ctx *gin.Context) {
	var request GetHealthRequestObject

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetHealth(ctx, request.(GetHealthRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetHealth")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(GetHealthResponseObject); ok {
		if err := validResponse.VisitGetHealthResponse(ctx.Writer); err != nil {
			ctx.Error(err)
		}
	} else if response != nil {
		ctx.Error(fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetVersion operation middleware
func (sh *strictHandler) GetVersion(ctx *gin.Context) {
	var request GetVersionRequestObject

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetVersion(ctx, request.(GetVersionRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetVersion")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(GetVersionResponseObject); ok {
		if err := validResponse.VisitGetVersionResponse(ctx.Writer); err != nil {
			ctx.Error(err)
		}
	} else if response != nil {
		ctx.Error(fmt.Errorf("unexpected response type: %T", response))
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9RUT28TPxD9Kqv5/Q4greK0FSLaW4WgFKiKKIUD6sFxJlk3Xtu1Z1NCtd8djb3Jtmmr",
	"cgAkLom98+fNvOeZG1Cu8c6ipQjVDURVYyPT8S1KQ/WrGtWSrz44j4E0Zj+S1KYTfpeNNwgVnL6HEmjt",
	"+RwpaLuArish4FWrA86g+rYJu9j6ueklKoKuhDMMKwyvQ3DhPhxuPg9ox5YwWGmKHFjkyKcKyIkewv+C",
	"IWpnH+l3la13S9gbjUfjJyE3ofdBuxIiqjZoWp8x7xnqyLmFwdPDlup9vs+Nu04G2VLtgv4hiet0M7z3",
	"8TwYqKAm8rESQirlWktxtEgZR8o1wgnHEfuCf6GEqJzPsAHlDCo4CtJSLPhWSKUwRijhOmjCwZiuG2tX",
	"Arkl7mBnlB5Zeh0TenJMfW+YSG6cY4oyYOCmuZh8e+NCIwkqePf1c6qVOYKqtw7EMyh0nFbbueN40pQU",
	"IozEmpSDgL1qXQnOo5VeQwUHvZBeUp3IEHV6/HxcIPEfP4bE8XFiCSmPB7DW0TsbM4v74zH/KWcJbQqU",
	"3hutUqi4jPkJ5Snj0/8B51DBf2IYQ9HPoLg9gKm7GUYVtKfcR2yTAPPWFNviuK0Xv7GE20P5QAkPz2Cq",
	"4eBv1dBDn1u5ktrIqcE8WG3TyLAeHDKbRaazBJKLyNPZC33BMeLWlD8me78m/qTudzbRPyX8Hd6PkIrD",
	"j8fFwNiG8+1G7PIK5BRsuNlB+dRaq+2iME5JY9bF3IVihis0zjfcYQntdulUQiS32kWqJpPJREivxWoP",
	"unI37YlTywyaMl61Wi1/Oe/Lx/N+YL/iRFu9bKdYPDuR6vluLrYEi4RxNHNqiWGkeyq3eS+6nwEAAP//",
	"+IR4cpUHAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
