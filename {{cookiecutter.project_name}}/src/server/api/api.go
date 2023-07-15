// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.0 DO NOT EDIT.
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

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Server Health Check
	// (GET /health)
	GetHealth(c *gin.Context)
	// Get Users
	// (GET /users)
	GetUsers(c *gin.Context, params GetUsersParams)
	// Add Users
	// (POST /users)
	AddUsers(c *gin.Context)
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

// GetUsers operation middleware
func (siw *ServerInterfaceWrapper) GetUsers(c *gin.Context) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetUsersParams

	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", c.Request.URL.Query(), &params.Offset)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter offset: %s", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", c.Request.URL.Query(), &params.Limit)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter limit: %s", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "sort" -------------

	err = runtime.BindQueryParameter("form", true, false, "sort", c.Request.URL.Query(), &params.Sort)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter sort: %s", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetUsers(c, params)
}

// AddUsers operation middleware
func (siw *ServerInterfaceWrapper) AddUsers(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.AddUsers(c)
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
	router.GET(options.BaseURL+"/users", wrapper.GetUsers)
	router.POST(options.BaseURL+"/users", wrapper.AddUsers)
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

type GetHealth500Response struct {
}

func (response GetHealth500Response) VisitGetHealthResponse(w http.ResponseWriter) error {
	w.WriteHeader(500)
	return nil
}

type GetHealth503Response struct {
}

func (response GetHealth503Response) VisitGetHealthResponse(w http.ResponseWriter) error {
	w.WriteHeader(503)
	return nil
}

type GetUsersRequestObject struct {
	Params GetUsersParams
}

type GetUsersResponseObject interface {
	VisitGetUsersResponse(w http.ResponseWriter) error
}

type GetUsers200JSONResponse Users

func (response GetUsers200JSONResponse) VisitGetUsersResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type AddUsersRequestObject struct {
	Body *AddUsersJSONRequestBody
}

type AddUsersResponseObject interface {
	VisitAddUsersResponse(w http.ResponseWriter) error
}

type AddUsers200JSONResponse UsersAdded

func (response AddUsers200JSONResponse) VisitAddUsersResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type AddUsers400Response struct {
}

func (response AddUsers400Response) VisitAddUsersResponse(w http.ResponseWriter) error {
	w.WriteHeader(400)
	return nil
}

type AddUsers401Response struct {
}

func (response AddUsers401Response) VisitAddUsersResponse(w http.ResponseWriter) error {
	w.WriteHeader(401)
	return nil
}

type AddUsers500Response struct {
}

func (response AddUsers500Response) VisitAddUsersResponse(w http.ResponseWriter) error {
	w.WriteHeader(500)
	return nil
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Server Health Check
	// (GET /health)
	GetHealth(ctx context.Context, request GetHealthRequestObject) (GetHealthResponseObject, error)
	// Get Users
	// (GET /users)
	GetUsers(ctx context.Context, request GetUsersRequestObject) (GetUsersResponseObject, error)
	// Add Users
	// (POST /users)
	AddUsers(ctx context.Context, request AddUsersRequestObject) (AddUsersResponseObject, error)
}

type StrictHandlerFunc = runtime.StrictGinHandlerFunc
type StrictMiddlewareFunc = runtime.StrictGinMiddlewareFunc

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
		ctx.Error(fmt.Errorf("Unexpected response type: %T", response))
	}
}

// GetUsers operation middleware
func (sh *strictHandler) GetUsers(ctx *gin.Context, params GetUsersParams) {
	var request GetUsersRequestObject

	request.Params = params

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetUsers(ctx, request.(GetUsersRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetUsers")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(GetUsersResponseObject); ok {
		if err := validResponse.VisitGetUsersResponse(ctx.Writer); err != nil {
			ctx.Error(err)
		}
	} else if response != nil {
		ctx.Error(fmt.Errorf("Unexpected response type: %T", response))
	}
}

// AddUsers operation middleware
func (sh *strictHandler) AddUsers(ctx *gin.Context) {
	var request AddUsersRequestObject

	var body AddUsersJSONRequestBody
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.Status(http.StatusBadRequest)
		ctx.Error(err)
		return
	}
	request.Body = &body

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.AddUsers(ctx, request.(AddUsersRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "AddUsers")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(AddUsersResponseObject); ok {
		if err := validResponse.VisitAddUsersResponse(ctx.Writer); err != nil {
			ctx.Error(err)
		}
	} else if response != nil {
		ctx.Error(fmt.Errorf("Unexpected response type: %T", response))
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7RVXW/bNhT9K8Td3qZESbthgZ6WdV3bDduAOkEfUj/Q0pXFmCKZy0s3XuH/PpC05W8s",
	"AZonUdL9PucefoXa9s4aNOyh+gq+7rCX6fgepebuTYf1LL46sg6JFWY7lhzSCR9l7zRCBf/8CQXwwsWz",
	"Z1JmCstlAYQPQRE2UN2t3caDnZ3cY82wLODWIx3mwV4qvZvm3nbml/T9vLb9QcoCHs+sdOqstg1O0Zzh",
	"I5M8YzlNAedSq0ZydFgXVuQksdZWkee/ZY+7Kf+wnYECnGRGMlDB58+//fAtMkvtOpkya3k0sezRv2jm",
	"PYA2A9gqqVjBcAq2XR7cDaiB76zu0e+gtTViGHVI2taz7WQVvE9OsCw2ge6/SPbWnAy0AmgryKfkAMtx",
	"AYqxTyV+T9hCBd+VG86XK8KXiX7LoUFJJBdDf9dNE8ezT04T+uHP0P7lEEMZxmmOakL/9lF5fqLtR3wI",
	"eGD96pj1sUX0oa7R+6du417KnWqLTZOH6Mf0WAdSvBjFMeapTFAS0nXgbvP2u6VecgTq0w0UWWVipPx3",
	"U2jH7DIplWlt9GfFqacb9JypPkfyysZduDy/OL+IQ7AOjXQKKnidPqWF6VI1ZZdkLB6nyPEREZSsrPnQ",
	"QAXvkLPQQZyNd9b43Mari4v4qK1hNMlROqdVnVzL+0iuQS//j1zbUpq6a9DXpBznPlZ4tUGLobjY1k+5",
	"hF3jDyYqgdRihDRHEm+JLGXr14fWK6NbI+dSaTnRmGELfS9psTHIFYpcYgFZOe5gNbxx9CnDetVPDTJr",
	"QRw+yR45Wd/tV3TToTChnyAJ24q0m4Kt8DPlxARbSyg8S4pgx++11RprFtyhIPRBs/DIEAkCFTwEpEUk",
	"aV5627b55waW/Y2JqnKqIL9TESEHMicyadWr5yYaWWIxWQjulBetQt2cCO4tHY09rPH4BbmaQXwyS3fI",
	"9A5ZrEmwplBkDYyXBTjrj9DmumnWHpQl6FfbLF6+mzeEklFIYfCLcMhCmcQxz5YQtnWSKeDypQeeJfZI",
	"naMT2vDjcW1IV71QxgXOVpeHVrdGBu4sqX9jyufJzA7c101zCu58N0TPYxrwMRgT91vbWmq9EK0l0eAc",
	"tXV9nGcBgfTqOqjKMpl11nN1dXV1VUqnyvklHC7YX7ae5aQp4kNQ9ezJcX8e4o6X/wUAAP//N7V8shYL",
	"AAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
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
	var res = make(map[string]func() ([]byte, error))
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
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
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