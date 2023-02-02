// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package openapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/go-chi/chi/v5"
)

// Defines values for MintType.
const (
	MintTypeMint MintType = "mint"
)

// Defines values for SpendType.
const (
	SpendTypeSpend SpendType = "spend"
)

// Defines values for TransferType.
const (
	TransferTypeTransfer TransferType = "transfer"
)

// Mint defines model for Mint.
type Mint struct {
	Amount     *int       `json:"amount,omitempty"`
	InsertedAt *time.Time `json:"inserted_at,omitempty"`
	Type       *MintType  `json:"type,omitempty"`
}

// MintType defines model for Mint.Type.
type MintType string

// Spend defines model for Spend.
type Spend struct {
	Amount     *int       `json:"amount,omitempty"`
	InsertedAt *time.Time `json:"inserted_at,omitempty"`
	Type       *SpendType `json:"type,omitempty"`
}

// SpendType defines model for Spend.Type.
type SpendType string

// Transfer defines model for Transfer.
type Transfer struct {
	Amount     *int          `json:"amount,omitempty"`
	InsertedAt *time.Time    `json:"inserted_at,omitempty"`
	Recipient  *int          `json:"recipient,omitempty"`
	Type       *TransferType `json:"type,omitempty"`
}

// TransferType defines model for Transfer.Type.
type TransferType string

// AccountId defines model for AccountId.
type AccountId = int

// RegisterJSONBody defines parameters for Register.
type RegisterJSONBody struct {
	Name string `json:"name"`
}

// MintJSONBody defines parameters for Mint.
type MintJSONBody struct {
	Amount int `json:"amount"`
}

// SpendJSONBody defines parameters for Spend.
type SpendJSONBody struct {
	Amount int `json:"amount"`
}

// TransferJSONBody defines parameters for Transfer.
type TransferJSONBody struct {
	Amount int `json:"amount"`
	To     int `json:"to"`
}

// RegisterJSONRequestBody defines body for Register for application/json ContentType.
type RegisterJSONRequestBody RegisterJSONBody

// MintJSONRequestBody defines body for Mint for application/json ContentType.
type MintJSONRequestBody MintJSONBody

// SpendJSONRequestBody defines body for Spend for application/json ContentType.
type SpendJSONRequestBody SpendJSONBody

// TransferJSONRequestBody defines body for Transfer for application/json ContentType.
type TransferJSONRequestBody TransferJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /accounts/register)
	Register(w http.ResponseWriter, r *http.Request)

	// (GET /accounts/{id}/balance)
	Balance(w http.ResponseWriter, r *http.Request, id AccountId)

	// (POST /accounts/{id}/mint)
	Mint(w http.ResponseWriter, r *http.Request, id AccountId)

	// (POST /accounts/{id}/spend)
	Spend(w http.ResponseWriter, r *http.Request, id AccountId)

	// (GET /accounts/{id}/transactions)
	Transactions(w http.ResponseWriter, r *http.Request, id AccountId)

	// (POST /accounts/{id}/transfer)
	Transfer(w http.ResponseWriter, r *http.Request, id AccountId)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// Register operation middleware
func (siw *ServerInterfaceWrapper) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Register(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// Balance operation middleware
func (siw *ServerInterfaceWrapper) Balance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "id" -------------
	var id AccountId

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, chi.URLParam(r, "id"), &id)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Balance(w, r, id)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// Mint operation middleware
func (siw *ServerInterfaceWrapper) Mint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "id" -------------
	var id AccountId

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, chi.URLParam(r, "id"), &id)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Mint(w, r, id)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// Spend operation middleware
func (siw *ServerInterfaceWrapper) Spend(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "id" -------------
	var id AccountId

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, chi.URLParam(r, "id"), &id)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Spend(w, r, id)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// Transactions operation middleware
func (siw *ServerInterfaceWrapper) Transactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "id" -------------
	var id AccountId

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, chi.URLParam(r, "id"), &id)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Transactions(w, r, id)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// Transfer operation middleware
func (siw *ServerInterfaceWrapper) Transfer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "id" -------------
	var id AccountId

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, chi.URLParam(r, "id"), &id)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Transfer(w, r, id)
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
		r.Post(options.BaseURL+"/accounts/register", wrapper.Register)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/accounts/{id}/balance", wrapper.Balance)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/accounts/{id}/mint", wrapper.Mint)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/accounts/{id}/spend", wrapper.Spend)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/accounts/{id}/transactions", wrapper.Transactions)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/accounts/{id}/transfer", wrapper.Transfer)
	})

	return r
}

type RegisterRequestObject struct {
	Body *RegisterJSONRequestBody
}

type RegisterResponseObject interface {
	VisitRegisterResponse(w http.ResponseWriter) error
}

type Register200JSONResponse struct {
	AccountId *int `json:"accountId,omitempty"`
}

func (response Register200JSONResponse) VisitRegisterResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type Register400TextResponse string

func (response Register400TextResponse) VisitRegisterResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(400)

	_, err := w.Write([]byte(response))
	return err
}

type BalanceRequestObject struct {
	Id AccountId `json:"id"`
}

type BalanceResponseObject interface {
	VisitBalanceResponse(w http.ResponseWriter) error
}

type Balance200JSONResponse struct {
	Balance *int `json:"balance,omitempty"`
}

func (response Balance200JSONResponse) VisitBalanceResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type Balance404Response struct {
}

func (response Balance404Response) VisitBalanceResponse(w http.ResponseWriter) error {
	w.WriteHeader(404)
	return nil
}

type MintRequestObject struct {
	Id   AccountId `json:"id"`
	Body *MintJSONRequestBody
}

type MintResponseObject interface {
	VisitMintResponse(w http.ResponseWriter) error
}

type Mint200JSONResponse struct {
	TransactionId *int `json:"transactionId,omitempty"`
}

func (response Mint200JSONResponse) VisitMintResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type Mint400TextResponse string

func (response Mint400TextResponse) VisitMintResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(400)

	_, err := w.Write([]byte(response))
	return err
}

type Mint404Response struct {
}

func (response Mint404Response) VisitMintResponse(w http.ResponseWriter) error {
	w.WriteHeader(404)
	return nil
}

type SpendRequestObject struct {
	Id   AccountId `json:"id"`
	Body *SpendJSONRequestBody
}

type SpendResponseObject interface {
	VisitSpendResponse(w http.ResponseWriter) error
}

type Spend200JSONResponse struct {
	TransactionId *int `json:"transactionId,omitempty"`
}

func (response Spend200JSONResponse) VisitSpendResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type Spend400TextResponse string

func (response Spend400TextResponse) VisitSpendResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(400)

	_, err := w.Write([]byte(response))
	return err
}

type Spend404Response struct {
}

func (response Spend404Response) VisitSpendResponse(w http.ResponseWriter) error {
	w.WriteHeader(404)
	return nil
}

type TransactionsRequestObject struct {
	Id AccountId `json:"id"`
}

type TransactionsResponseObject interface {
	VisitTransactionsResponse(w http.ResponseWriter) error
}

type Transactions200JSONResponse []interface{}

func (response Transactions200JSONResponse) VisitTransactionsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type Transactions404Response struct {
}

func (response Transactions404Response) VisitTransactionsResponse(w http.ResponseWriter) error {
	w.WriteHeader(404)
	return nil
}

type TransferRequestObject struct {
	Id   AccountId `json:"id"`
	Body *TransferJSONRequestBody
}

type TransferResponseObject interface {
	VisitTransferResponse(w http.ResponseWriter) error
}

type Transfer200JSONResponse struct {
	TransactionId *int `json:"transactionId,omitempty"`
}

func (response Transfer200JSONResponse) VisitTransferResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type Transfer400TextResponse string

func (response Transfer400TextResponse) VisitTransferResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(400)

	_, err := w.Write([]byte(response))
	return err
}

type Transfer404Response struct {
}

func (response Transfer404Response) VisitTransferResponse(w http.ResponseWriter) error {
	w.WriteHeader(404)
	return nil
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {

	// (POST /accounts/register)
	Register(ctx context.Context, request RegisterRequestObject) (RegisterResponseObject, error)

	// (GET /accounts/{id}/balance)
	Balance(ctx context.Context, request BalanceRequestObject) (BalanceResponseObject, error)

	// (POST /accounts/{id}/mint)
	Mint(ctx context.Context, request MintRequestObject) (MintResponseObject, error)

	// (POST /accounts/{id}/spend)
	Spend(ctx context.Context, request SpendRequestObject) (SpendResponseObject, error)

	// (GET /accounts/{id}/transactions)
	Transactions(ctx context.Context, request TransactionsRequestObject) (TransactionsResponseObject, error)

	// (POST /accounts/{id}/transfer)
	Transfer(ctx context.Context, request TransferRequestObject) (TransferResponseObject, error)
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

// Register operation middleware
func (sh *strictHandler) Register(w http.ResponseWriter, r *http.Request) {
	var request RegisterRequestObject

	var body RegisterJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.Register(ctx, request.(RegisterRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Register")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(RegisterResponseObject); ok {
		if err := validResponse.VisitRegisterResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// Balance operation middleware
func (sh *strictHandler) Balance(w http.ResponseWriter, r *http.Request, id AccountId) {
	var request BalanceRequestObject

	request.Id = id

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.Balance(ctx, request.(BalanceRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Balance")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(BalanceResponseObject); ok {
		if err := validResponse.VisitBalanceResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// Mint operation middleware
func (sh *strictHandler) Mint(w http.ResponseWriter, r *http.Request, id AccountId) {
	var request MintRequestObject

	request.Id = id

	var body MintJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.Mint(ctx, request.(MintRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Mint")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(MintResponseObject); ok {
		if err := validResponse.VisitMintResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// Spend operation middleware
func (sh *strictHandler) Spend(w http.ResponseWriter, r *http.Request, id AccountId) {
	var request SpendRequestObject

	request.Id = id

	var body SpendJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.Spend(ctx, request.(SpendRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Spend")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(SpendResponseObject); ok {
		if err := validResponse.VisitSpendResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// Transactions operation middleware
func (sh *strictHandler) Transactions(w http.ResponseWriter, r *http.Request, id AccountId) {
	var request TransactionsRequestObject

	request.Id = id

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.Transactions(ctx, request.(TransactionsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Transactions")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(TransactionsResponseObject); ok {
		if err := validResponse.VisitTransactionsResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// Transfer operation middleware
func (sh *strictHandler) Transfer(w http.ResponseWriter, r *http.Request, id AccountId) {
	var request TransferRequestObject

	request.Id = id

	var body TransferJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.Transfer(ctx, request.(TransferRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Transfer")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(TransferResponseObject); ok {
		if err := validResponse.VisitTransferResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}
