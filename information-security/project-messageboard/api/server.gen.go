// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/go-chi/chi/v5"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /api/replies/{board})
	GetReplies(w http.ResponseWriter, r *http.Request, board Board, params GetRepliesParams)

	// (POST /api/replies/{board})
	CreateReply(w http.ResponseWriter, r *http.Request, board Board)

	// (GET /api/threads/{board})
	GetThreads(w http.ResponseWriter, r *http.Request, board Board)

	// (POST /api/threads/{board})
	CreateThread(w http.ResponseWriter, r *http.Request, board Board)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetReplies operation middleware
func (siw *ServerInterfaceWrapper) GetReplies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "board" -------------
	var board Board

	err = runtime.BindStyledParameterWithLocation("simple", false, "board", runtime.ParamLocationPath, chi.URLParam(r, "board"), &board)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "board", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetRepliesParams

	// ------------- Required query parameter "thread_id" -------------

	if paramValue := r.URL.Query().Get("thread_id"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "thread_id"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "thread_id", r.URL.Query(), &params.ThreadId)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "thread_id", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetReplies(w, r, board, params)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// CreateReply operation middleware
func (siw *ServerInterfaceWrapper) CreateReply(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "board" -------------
	var board Board

	err = runtime.BindStyledParameterWithLocation("simple", false, "board", runtime.ParamLocationPath, chi.URLParam(r, "board"), &board)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "board", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateReply(w, r, board)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetThreads operation middleware
func (siw *ServerInterfaceWrapper) GetThreads(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "board" -------------
	var board Board

	err = runtime.BindStyledParameterWithLocation("simple", false, "board", runtime.ParamLocationPath, chi.URLParam(r, "board"), &board)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "board", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetThreads(w, r, board)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// CreateThread operation middleware
func (siw *ServerInterfaceWrapper) CreateThread(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "board" -------------
	var board Board

	err = runtime.BindStyledParameterWithLocation("simple", false, "board", runtime.ParamLocationPath, chi.URLParam(r, "board"), &board)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "board", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateThread(w, r, board)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshallingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshallingParamError) Error() string {
	return fmt.Sprintf("Error unmarshalling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshallingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/api/replies/{board}", wrapper.GetReplies)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/api/replies/{board}", wrapper.CreateReply)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/api/threads/{board}", wrapper.GetThreads)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/api/threads/{board}", wrapper.CreateThread)
	})

	return r
}

type DefaultTextResponse string

type GetRepliesRequestObject struct {
	Board  Board `json:"board"`
	Params GetRepliesParams
}

type GetRepliesResponseObject interface {
	VisitGetRepliesResponse(w http.ResponseWriter) error
}

type GetReplies200JSONResponse Thread

func (response GetReplies200JSONResponse) VisitGetRepliesResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetRepliesdefaultTextResponse struct {
	Body       string
	StatusCode int
}

func (response GetRepliesdefaultTextResponse) VisitGetRepliesResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(response.StatusCode)

	_, err := w.Write([]byte(response.Body))
	return err
}

type CreateReplyRequestObject struct {
	Board        Board `json:"board"`
	JSONBody     *CreateReplyJSONRequestBody
	FormdataBody *CreateReplyFormdataRequestBody
}

type CreateReplyResponseObject interface {
	VisitCreateReplyResponse(w http.ResponseWriter) error
}

type CreateReply200TextResponse string

func (response CreateReply200TextResponse) VisitCreateReplyResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)

	_, err := w.Write([]byte(response))
	return err
}

type CreateReplydefaultTextResponse struct {
	Body       string
	StatusCode int
}

func (response CreateReplydefaultTextResponse) VisitCreateReplyResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(response.StatusCode)

	_, err := w.Write([]byte(response.Body))
	return err
}

type GetThreadsRequestObject struct {
	Board Board `json:"board"`
}

type GetThreadsResponseObject interface {
	VisitGetThreadsResponse(w http.ResponseWriter) error
}

type GetThreads200JSONResponse []Thread

func (response GetThreads200JSONResponse) VisitGetThreadsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetThreadsdefaultTextResponse struct {
	Body       string
	StatusCode int
}

func (response GetThreadsdefaultTextResponse) VisitGetThreadsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(response.StatusCode)

	_, err := w.Write([]byte(response.Body))
	return err
}

type CreateThreadRequestObject struct {
	Board        Board `json:"board"`
	JSONBody     *CreateThreadJSONRequestBody
	FormdataBody *CreateThreadFormdataRequestBody
}

type CreateThreadResponseObject interface {
	VisitCreateThreadResponse(w http.ResponseWriter) error
}

type CreateThread200TextResponse string

func (response CreateThread200TextResponse) VisitCreateThreadResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)

	_, err := w.Write([]byte(response))
	return err
}

type CreateThreaddefaultTextResponse struct {
	Body       string
	StatusCode int
}

func (response CreateThreaddefaultTextResponse) VisitCreateThreadResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(response.StatusCode)

	_, err := w.Write([]byte(response.Body))
	return err
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {

	// (GET /api/replies/{board})
	GetReplies(ctx context.Context, request GetRepliesRequestObject) (GetRepliesResponseObject, error)

	// (POST /api/replies/{board})
	CreateReply(ctx context.Context, request CreateReplyRequestObject) (CreateReplyResponseObject, error)

	// (GET /api/threads/{board})
	GetThreads(ctx context.Context, request GetThreadsRequestObject) (GetThreadsResponseObject, error)

	// (POST /api/threads/{board})
	CreateThread(ctx context.Context, request CreateThreadRequestObject) (CreateThreadResponseObject, error)
}

type StrictHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request, args interface{}) (interface{}, error)

type StrictMiddlewareFunc func(f StrictHandlerFunc, operationID string) StrictHandlerFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// GetReplies operation middleware
func (sh *strictHandler) GetReplies(w http.ResponseWriter, r *http.Request, board Board, params GetRepliesParams) {
	var request GetRepliesRequestObject

	request.Board = board
	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetReplies(ctx, request.(GetRepliesRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetReplies")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetRepliesResponseObject); ok {
		if err := validResponse.VisitGetRepliesResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// CreateReply operation middleware
func (sh *strictHandler) CreateReply(w http.ResponseWriter, r *http.Request, board Board) {
	var request CreateReplyRequestObject

	request.Board = board
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		var body CreateReplyJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
			return
		}
		request.JSONBody = &body
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
		if err := r.ParseForm(); err != nil {
			sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode formdata: %w", err))
			return
		}
		var body CreateReplyFormdataRequestBody
		if err := runtime.BindForm(&body, r.Form, nil, nil); err != nil {
			sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't bind formdata: %w", err))
			return
		}
		request.FormdataBody = &body
	}

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.CreateReply(ctx, request.(CreateReplyRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "CreateReply")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(CreateReplyResponseObject); ok {
		if err := validResponse.VisitCreateReplyResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// GetThreads operation middleware
func (sh *strictHandler) GetThreads(w http.ResponseWriter, r *http.Request, board Board) {
	var request GetThreadsRequestObject

	request.Board = board

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetThreads(ctx, request.(GetThreadsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetThreads")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetThreadsResponseObject); ok {
		if err := validResponse.VisitGetThreadsResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// CreateThread operation middleware
func (sh *strictHandler) CreateThread(w http.ResponseWriter, r *http.Request, board Board) {
	var request CreateThreadRequestObject

	request.Board = board
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		var body CreateThreadJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
			return
		}
		request.JSONBody = &body
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
		if err := r.ParseForm(); err != nil {
			sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode formdata: %w", err))
			return
		}
		var body CreateThreadFormdataRequestBody
		if err := runtime.BindForm(&body, r.Form, nil, nil); err != nil {
			sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't bind formdata: %w", err))
			return
		}
		request.FormdataBody = &body
	}

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.CreateThread(ctx, request.(CreateThreadRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "CreateThread")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(CreateThreadResponseObject); ok {
		if err := validResponse.VisitCreateThreadResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}
