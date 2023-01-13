// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
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
	"time"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/iden3/go-merkletree-sql"
)

// CreateClaimRequest defines model for CreateClaimRequest.
type CreateClaimRequest struct {
	CredentialSchema      string          `json:"credentialSchema"`
	CredentialSubject     json.RawMessage `json:"credentialSubject"`
	Expiration            *int64          `json:"expiration,omitempty"`
	MerklizedRootPosition *string         `json:"merklizedRootPosition,omitempty"`
	RevNonce              *uint64         `json:"revNonce,omitempty"`
	SubjectPosition       *string         `json:"subjectPosition,omitempty"`
	Type                  string          `json:"type"`
	Version               *uint32         `json:"version,omitempty"`
}

// CreateClaimResponse defines model for CreateClaimResponse.
type CreateClaimResponse struct {
	Id string `json:"id"`
}

// CreateIdentityResponse defines model for CreateIdentityResponse.
type CreateIdentityResponse struct {
	Identifier *string        `json:"identifier,omitempty"`
	Immutable  bool           `json:"immutable"`
	Relay      string         `json:"relay"`
	State      *IdentityState `json:"state,omitempty"`
}

// GenericErrorMessage defines model for GenericErrorMessage.
type GenericErrorMessage struct {
	Message string `json:"message"`
}

// Health defines model for Health.
type Health struct {
	Cache bool `json:"cache"`
	Db    bool `json:"db"`
}

// IdentityState defines model for IdentityState.
type IdentityState struct {
	BlockNumber        *int      `json:"blockNumber,omitempty"`
	BlockTimestamp     *int      `json:"blockTimestamp,omitempty"`
	ClaimsTreeRoot     *string   `json:"claimsTreeRoot,omitempty"`
	CreatedAt          time.Time `json:"createdAt"`
	Identifier         string    `json:"-"`
	ModifiedAt         time.Time `json:"modifiedAt"`
	PreviousState      *string   `json:"previousState,omitempty"`
	RevocationTreeRoot *string   `json:"revocationTreeRoot,omitempty"`
	RootOfRoots        *string   `json:"rootOfRoots,omitempty"`
	State              *string   `json:"state,omitempty"`
	StateID            int64     `json:"-"`
	Status             string    `json:"status"`
	TxID               *string   `json:"txID,omitempty"`
}

// PublishStateResponse defines model for PublishStateResponse.
type PublishStateResponse struct {
	Hex *string `json:"hex,omitempty"`
}

// RevocationStatusResponse defines model for RevocationStatusResponse.
type RevocationStatusResponse struct {
	Issuer struct {
		ClaimsTreeRoot     *string `json:"claimsTreeRoot,omitempty"`
		RevocationTreeRoot *string `json:"revocationTreeRoot,omitempty"`
		RootOfRoots        *string `json:"rootOfRoots,omitempty"`
		State              *string `json:"state,omitempty"`
	} `json:"issuer"`
	Mtp struct {
		Existence bool `json:"existence"`
		NodeAux   *struct {
			Key   *merkletree.Hash `json:"key,omitempty"`
			Value *merkletree.Hash `json:"value,omitempty"`
		} `json:"nodeAux,omitempty"`
	} `json:"mtp"`
}

// RevokeClaimResponse defines model for RevokeClaimResponse.
type RevokeClaimResponse struct {
	Status string `json:"status"`
}

// PathIdentifier defines model for pathIdentifier.
type PathIdentifier = string

// PathNonce defines model for pathNonce.
type PathNonce = int64

// N400 defines model for 400.
type N400 = GenericErrorMessage

// N404 defines model for 404.
type N404 = GenericErrorMessage

// N500 defines model for 500.
type N500 = GenericErrorMessage

// N500CreateIdentity defines model for 500-CreateIdentity.
type N500CreateIdentity struct {
	Code      *int    `json:"code,omitempty"`
	Error     *string `json:"error,omitempty"`
	RequestID *string `json:"requestID,omitempty"`
}

// CreateClaimJSONRequestBody defines body for CreateClaim for application/json ContentType.
type CreateClaimJSONRequestBody = CreateClaimRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get the documentation
	// (GET /)
	GetDocumentation(w http.ResponseWriter, r *http.Request)
	// Get the documentation yaml file
	// (GET /static/docs/api/api.yaml)
	GetYaml(w http.ResponseWriter, r *http.Request)
	// Healthcheck
	// (GET /status)
	Health(w http.ResponseWriter, r *http.Request)
	// Create Identity
	// (POST /v1/identities)
	CreateIdentity(w http.ResponseWriter, r *http.Request)
	// Publish State On-Chain
	// (POST /v1/identities/state)
	PublishState(w http.ResponseWriter, r *http.Request)
	// Create Claim
	// (POST /v1/{identifier}/claims)
	CreateClaim(w http.ResponseWriter, r *http.Request, identifier PathIdentifier)
	// Get Revocation Status
	// (GET /v1/{identifier}/claims/revocation/status/{nonce})
	GetRevocationStatus(w http.ResponseWriter, r *http.Request, identifier PathIdentifier, nonce PathNonce)
	// Revoke Claim
	// (POST /v1/{identifier}/claims/revoke/{nonce})
	RevokeClaim(w http.ResponseWriter, r *http.Request, identifier PathIdentifier, nonce PathNonce)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetDocumentation operation middleware
func (siw *ServerInterfaceWrapper) GetDocumentation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetDocumentation(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetYaml operation middleware
func (siw *ServerInterfaceWrapper) GetYaml(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetYaml(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// Health operation middleware
func (siw *ServerInterfaceWrapper) Health(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Health(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// CreateIdentity operation middleware
func (siw *ServerInterfaceWrapper) CreateIdentity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateIdentity(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PublishState operation middleware
func (siw *ServerInterfaceWrapper) PublishState(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PublishState(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// CreateClaim operation middleware
func (siw *ServerInterfaceWrapper) CreateClaim(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "identifier" -------------
	var identifier PathIdentifier

	err = runtime.BindStyledParameterWithLocation("simple", false, "identifier", runtime.ParamLocationPath, chi.URLParam(r, "identifier"), &identifier)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "identifier", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateClaim(w, r, identifier)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetRevocationStatus operation middleware
func (siw *ServerInterfaceWrapper) GetRevocationStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "identifier" -------------
	var identifier PathIdentifier

	err = runtime.BindStyledParameterWithLocation("simple", false, "identifier", runtime.ParamLocationPath, chi.URLParam(r, "identifier"), &identifier)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "identifier", Err: err})
		return
	}

	// ------------- Path parameter "nonce" -------------
	var nonce PathNonce

	err = runtime.BindStyledParameterWithLocation("simple", false, "nonce", runtime.ParamLocationPath, chi.URLParam(r, "nonce"), &nonce)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "nonce", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetRevocationStatus(w, r, identifier, nonce)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// RevokeClaim operation middleware
func (siw *ServerInterfaceWrapper) RevokeClaim(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "identifier" -------------
	var identifier PathIdentifier

	err = runtime.BindStyledParameterWithLocation("simple", false, "identifier", runtime.ParamLocationPath, chi.URLParam(r, "identifier"), &identifier)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "identifier", Err: err})
		return
	}

	// ------------- Path parameter "nonce" -------------
	var nonce PathNonce

	err = runtime.BindStyledParameterWithLocation("simple", false, "nonce", runtime.ParamLocationPath, chi.URLParam(r, "nonce"), &nonce)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "nonce", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.RevokeClaim(w, r, identifier, nonce)
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
		r.Get(options.BaseURL+"/", wrapper.GetDocumentation)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/static/docs/api/api.yaml", wrapper.GetYaml)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/status", wrapper.Health)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1/identities", wrapper.CreateIdentity)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1/identities/state", wrapper.PublishState)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1/{identifier}/claims", wrapper.CreateClaim)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/v1/{identifier}/claims/revocation/status/{nonce}", wrapper.GetRevocationStatus)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1/{identifier}/claims/revoke/{nonce}", wrapper.RevokeClaim)
	})

	return r
}

type N400JSONResponse GenericErrorMessage

type N404JSONResponse GenericErrorMessage

type N500JSONResponse GenericErrorMessage

type N500CreateIdentityJSONResponse struct {
	Code      *int    `json:"code,omitempty"`
	Error     *string `json:"error,omitempty"`
	RequestID *string `json:"requestID,omitempty"`
}

type GetDocumentationRequestObject struct {
}

type GetDocumentationResponseObject interface {
	VisitGetDocumentationResponse(w http.ResponseWriter) error
}

type GetDocumentation200Response struct {
}

func (response GetDocumentation200Response) VisitGetDocumentationResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type GetYamlRequestObject struct {
}

type GetYamlResponseObject interface {
	VisitGetYamlResponse(w http.ResponseWriter) error
}

type GetYaml200Response struct {
}

func (response GetYaml200Response) VisitGetYamlResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type HealthRequestObject struct {
}

type HealthResponseObject interface {
	VisitHealthResponse(w http.ResponseWriter) error
}

type Health200JSONResponse Health

func (response Health200JSONResponse) VisitHealthResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type Health500JSONResponse struct{ N500JSONResponse }

func (response Health500JSONResponse) VisitHealthResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type CreateIdentityRequestObject struct {
}

type CreateIdentityResponseObject interface {
	VisitCreateIdentityResponse(w http.ResponseWriter) error
}

type CreateIdentity201JSONResponse CreateIdentityResponse

func (response CreateIdentity201JSONResponse) VisitCreateIdentityResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	return json.NewEncoder(w).Encode(response)
}

type CreateIdentity500JSONResponse struct{ N500CreateIdentityJSONResponse }

func (response CreateIdentity500JSONResponse) VisitCreateIdentityResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type PublishStateRequestObject struct {
}

type PublishStateResponseObject interface {
	VisitPublishStateResponse(w http.ResponseWriter) error
}

type PublishState200JSONResponse PublishStateResponse

func (response PublishState200JSONResponse) VisitPublishStateResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type PublishState500JSONResponse struct{ N500JSONResponse }

func (response PublishState500JSONResponse) VisitPublishStateResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type CreateClaimRequestObject struct {
	Identifier PathIdentifier `json:"identifier"`
	Body       *CreateClaimJSONRequestBody
}

type CreateClaimResponseObject interface {
	VisitCreateClaimResponse(w http.ResponseWriter) error
}

type CreateClaim201JSONResponse CreateClaimResponse

func (response CreateClaim201JSONResponse) VisitCreateClaimResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	return json.NewEncoder(w).Encode(response)
}

type CreateClaim400JSONResponse struct{ N400JSONResponse }

func (response CreateClaim400JSONResponse) VisitCreateClaimResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type CreateClaim500JSONResponse struct{ N500JSONResponse }

func (response CreateClaim500JSONResponse) VisitCreateClaimResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type GetRevocationStatusRequestObject struct {
	Identifier PathIdentifier `json:"identifier"`
	Nonce      PathNonce      `json:"nonce"`
}

type GetRevocationStatusResponseObject interface {
	VisitGetRevocationStatusResponse(w http.ResponseWriter) error
}

type GetRevocationStatus200JSONResponse RevocationStatusResponse

func (response GetRevocationStatus200JSONResponse) VisitGetRevocationStatusResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetRevocationStatus400JSONResponse struct{ N400JSONResponse }

func (response GetRevocationStatus400JSONResponse) VisitGetRevocationStatusResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type GetRevocationStatus500JSONResponse struct{ N500JSONResponse }

func (response GetRevocationStatus500JSONResponse) VisitGetRevocationStatusResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type RevokeClaimRequestObject struct {
	Identifier PathIdentifier `json:"identifier"`
	Nonce      PathNonce      `json:"nonce"`
}

type RevokeClaimResponseObject interface {
	VisitRevokeClaimResponse(w http.ResponseWriter) error
}

type RevokeClaim202JSONResponse RevokeClaimResponse

func (response RevokeClaim202JSONResponse) VisitRevokeClaimResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(202)

	return json.NewEncoder(w).Encode(response)
}

type RevokeClaim404JSONResponse struct{ N404JSONResponse }

func (response RevokeClaim404JSONResponse) VisitRevokeClaimResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type RevokeClaim500JSONResponse struct{ N500JSONResponse }

func (response RevokeClaim500JSONResponse) VisitRevokeClaimResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Get the documentation
	// (GET /)
	GetDocumentation(ctx context.Context, request GetDocumentationRequestObject) (GetDocumentationResponseObject, error)
	// Get the documentation yaml file
	// (GET /static/docs/api/api.yaml)
	GetYaml(ctx context.Context, request GetYamlRequestObject) (GetYamlResponseObject, error)
	// Healthcheck
	// (GET /status)
	Health(ctx context.Context, request HealthRequestObject) (HealthResponseObject, error)
	// Create Identity
	// (POST /v1/identities)
	CreateIdentity(ctx context.Context, request CreateIdentityRequestObject) (CreateIdentityResponseObject, error)
	// Publish State On-Chain
	// (POST /v1/identities/state)
	PublishState(ctx context.Context, request PublishStateRequestObject) (PublishStateResponseObject, error)
	// Create Claim
	// (POST /v1/{identifier}/claims)
	CreateClaim(ctx context.Context, request CreateClaimRequestObject) (CreateClaimResponseObject, error)
	// Get Revocation Status
	// (GET /v1/{identifier}/claims/revocation/status/{nonce})
	GetRevocationStatus(ctx context.Context, request GetRevocationStatusRequestObject) (GetRevocationStatusResponseObject, error)
	// Revoke Claim
	// (POST /v1/{identifier}/claims/revoke/{nonce})
	RevokeClaim(ctx context.Context, request RevokeClaimRequestObject) (RevokeClaimResponseObject, error)
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

// GetDocumentation operation middleware
func (sh *strictHandler) GetDocumentation(w http.ResponseWriter, r *http.Request) {
	var request GetDocumentationRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetDocumentation(ctx, request.(GetDocumentationRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetDocumentation")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetDocumentationResponseObject); ok {
		if err := validResponse.VisitGetDocumentationResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// GetYaml operation middleware
func (sh *strictHandler) GetYaml(w http.ResponseWriter, r *http.Request) {
	var request GetYamlRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetYaml(ctx, request.(GetYamlRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetYaml")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetYamlResponseObject); ok {
		if err := validResponse.VisitGetYamlResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// Health operation middleware
func (sh *strictHandler) Health(w http.ResponseWriter, r *http.Request) {
	var request HealthRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.Health(ctx, request.(HealthRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Health")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(HealthResponseObject); ok {
		if err := validResponse.VisitHealthResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// CreateIdentity operation middleware
func (sh *strictHandler) CreateIdentity(w http.ResponseWriter, r *http.Request) {
	var request CreateIdentityRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.CreateIdentity(ctx, request.(CreateIdentityRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "CreateIdentity")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(CreateIdentityResponseObject); ok {
		if err := validResponse.VisitCreateIdentityResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// PublishState operation middleware
func (sh *strictHandler) PublishState(w http.ResponseWriter, r *http.Request) {
	var request PublishStateRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.PublishState(ctx, request.(PublishStateRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PublishState")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(PublishStateResponseObject); ok {
		if err := validResponse.VisitPublishStateResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// CreateClaim operation middleware
func (sh *strictHandler) CreateClaim(w http.ResponseWriter, r *http.Request, identifier PathIdentifier) {
	var request CreateClaimRequestObject

	request.Identifier = identifier

	var body CreateClaimJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.CreateClaim(ctx, request.(CreateClaimRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "CreateClaim")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(CreateClaimResponseObject); ok {
		if err := validResponse.VisitCreateClaimResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// GetRevocationStatus operation middleware
func (sh *strictHandler) GetRevocationStatus(w http.ResponseWriter, r *http.Request, identifier PathIdentifier, nonce PathNonce) {
	var request GetRevocationStatusRequestObject

	request.Identifier = identifier
	request.Nonce = nonce

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetRevocationStatus(ctx, request.(GetRevocationStatusRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetRevocationStatus")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetRevocationStatusResponseObject); ok {
		if err := validResponse.VisitGetRevocationStatusResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// RevokeClaim operation middleware
func (sh *strictHandler) RevokeClaim(w http.ResponseWriter, r *http.Request, identifier PathIdentifier, nonce PathNonce) {
	var request RevokeClaimRequestObject

	request.Identifier = identifier
	request.Nonce = nonce

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.RevokeClaim(ctx, request.(RevokeClaimRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "RevokeClaim")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(RevokeClaimResponseObject); ok {
		if err := validResponse.VisitRevokeClaimResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8RZ2W7bOBd+FYL/fylbTuqZAXyXpkUToEtQ92bQyQVNHVusJVIlKcOewO8+4KLNory0",
	"SXvRJeLhWT5+ZyHzhKnIC8GBa4VnT7ggkuSgQfqfdHqfANdsyUCaLwkoKlmhmeB4ht2a3iHWCEWYmSWz",
	"FUeYkxzwDHfWJXwvmYQEz7QsIcKKppATo13vCiOttGR8hff7yKr5KDiFvvHbjLAccbsYNFotDdtbCpkT",
	"bfzj+s8pjioHGNewAon3xgUJqhBcgUVkOpmYf6jgGrg2/yVFkTFKjFPxN2U8e2pZ+L+EJZ7h/8UNzLFb",
	"VfE74CAZfSulkB9AKbICZ7Eb52uSoM/wvQSl8T7C08n0V3vwUWi0FCVPjP0/fj0C91yD5CRDc5AbkAiM",
	"vPdldCuBaKiYeJFrhRQFSM3c0VKRQIuENQci7Oz1+emoBUrfvwmz138Ri29A9QWR7SuSWsdchJbuFQ/6",
	"zkuwCJBsPpBNEd6ORM405IWBaUkyBfuovbF0brbzwmIWnacJtgWTxAV2RmpFOAe5zti/kHwWQj8Ixaq9",
	"AZg3dRGoFZfDmpUL5ahO9+FMlDYg1WFgxv6r61DRaJecr/2j8TtC2D/2OBN1j9+Vov75s+TMWA6cY8kR",
	"m1VSHTPbbg49jFmel5ossjbQCyEyINwda0Z2wY1KEw2nqkfl3twK9yPzFaGy03bHWShVMPhQVepFnjcL",
	"sCV5YWLEc5GDThlfoZQUBfRT59DJSkvIjTsgmU4DqU5o2rXr2lof32RxhtiBR8nCENOaCDnVxbzn2yIT",
	"dP2xzBcdQrQy0wp8YTkoTfIiLEMN1dUXCWAKQ5Af1PIzuelWq4RoGGmWQx/36ChVTaasxMhUuxFbcSE9",
	"WqZIicRsushUIWHDRKlqkEIFTbi2dDRMKYT+tDTL6niahFdcWzpViY8E75MkWD63Q02vzabKjag7AXq9",
	"HXTbhxpi3kO5yJhKLabDBSmF7VmtOMKf6yOYW2+OFDmlSseag0Q8zdMXPudQXLku+q7ClikNvoP2CwUX",
	"CdyU2/6+NewGs8V/tF0ctAQY3xGVtldHLC+EtCH7kdxL2Dl9hldMp+ViTEUeG3q8ildi1Kgbqe+ZDXFD",
	"shJ+txuBaa5L9Qbix1OinlDurB4HuLk+1fKb5GwaUAE8cdhcPgoMdkQjyPhS9K9gbwQtc+DaMhwthUQ6",
	"BTSHbInuhNKQoPs36CEj2tSgf2w7ZNp1yrAMbo1aeDK+Gk8MHKIATgqGZ/iV/eROzoYem79WYA/XgGM9",
	"uU/wDL8D3XEPH9zkrt09phuQKikFpRDhCZKgS8mVDSnpBMo4uvvy4T3yldWO62WeE7lzdvtb7AEwP+z7",
	"8rqPcGwwZzROBFUxKZj5M96RPDsW1d9m/VmDMRovCAbtrDyzc9RgWI6awSD8ZBOO4Vkult5C4MZ1k2VI",
	"gdwwCgoRCUiWnPsy66+2IcW1p7ER6oLkjNEU6NquxJur2E+fPlUL4e5sXVfe8qQQjGukBXK9DxGOWnNr",
	"F7WDi24PvatnQ29g+g/dX6sHIN+7L0Hx8ObeBdUtola4mqyUqVP1p8c+2HHdKU9DXriBAtktSPARTQnj",
	"Pdjbc8dLUjY43wQgtwKV85chfgCxt4icxk98dOsBGEb6qRnj9rGbfy7jN7JNbYDb1Vr7EfJrOK5GJD54",
	"pNw/1s8yr0Wye+ak6LzC7Lst1A/NL5yW3akgQBD3LtpKyOk59DBCP04ln63VAVYEcj8fY0/cTMi+a8RP",
	"9t123+oew7Ra+f7UaEH13aLXOA8n/p9mWnTWDvdy5Wj5QrVj8C4ToMeDFGL5i2hhxofGN1TDfiE/1tAm",
	"xeli4/Yggmiw2LSm699NgetnpcD6dHG4oRSKui5MzyHA9CcI4LwargtG2L58O/C7vr4XlJhJt5QZnuFU",
	"62IWx1fXf40n48n4ysLpFR7unAO1hLtC4Gmhml8L1Yv29MIbr49tvA5svBVZ5pfFstmMJGSmDBtWtkYZ",
	"r7AZfX5E361rv7U2B+r+cf9fAAAA//9c8yLN1BsAAA==",
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
