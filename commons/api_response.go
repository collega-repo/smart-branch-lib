package commons

import (
	"errors"
	"fmt"
	"github.com/collega-repo/smart-branch-lib/commons/errs"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	status2 "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"net/http"
)

type status string

const (
	SuccessResponse status = "SUCCESS"
	PendingResponse status = "PENDING"
	FailedResponse  status = "FAILED"
)

type (
	Code      string
	errorCode string
)

const (
	CodeSuccess                Code      = "00"
	CodeSuccessCBS             Code      = "01"
	CodeSuccessPending         Code      = "02"
	CodeSuccessPendingApproval Code      = "03"
	CodeSuccessNoData          Code      = "04"
	CodeNotFoundCore           Code      = "84"
	CodeUnAuthenticationCore   Code      = "87"
	CodeInsufficientBalance    Code      = "90"
	CodeHtxNotFound            Code      = "91"
	CodeTrxNotFound            Code      = "92"
	CodeInternalError          Code      = "93"
	CodeNotFound               Code      = "94"
	CodeIdempotencyFailed      Code      = "95"
	CodeForbiddenAccess        Code      = "96"
	CodeUnAuthentication       Code      = "97"
	CodeInvalidRequest         Code      = "98"
	CodeFailed                 Code      = "99"
	InvalidToken               errorCode = "invalid_token"
	InvalidRequest             errorCode = "invalid_request"
	InvalidClient              errorCode = "invalid_client"
	InvalidPassword            errorCode = "invalid_password"
	InvalidUser                errorCode = "invalid_user"
	InvalidServer              errorCode = "invalid_server"
	InvalidCore                errorCode = "invalid_core"
	InvalidApplication         errorCode = "invalid_application"
	AccessDenied               errorCode = "access_denied"
)

var MapStatusCode = map[Code]int{
	CodeSuccess:                http.StatusOK,
	CodeSuccessCBS:             http.StatusInternalServerError,
	CodeSuccessPending:         http.StatusAccepted,
	CodeSuccessPendingApproval: http.StatusAccepted,
	CodeNotFoundCore:           http.StatusNotFound,
	CodeUnAuthenticationCore:   http.StatusUnauthorized,
	CodeInternalError:          http.StatusInternalServerError,
	CodeNotFound:               http.StatusNotFound,
	CodeForbiddenAccess:        http.StatusForbidden,
	CodeUnAuthentication:       http.StatusUnauthorized,
	CodeInvalidRequest:         http.StatusUnprocessableEntity,
	CodeInsufficientBalance:    http.StatusBadRequest,
	CodeFailed:                 http.StatusBadRequest,
}

var MapStatusGrpc = map[Code]codes.Code{
	CodeSuccess:              codes.OK,
	CodeNotFound:             codes.NotFound,
	CodeNotFoundCore:         codes.NotFound,
	CodeUnAuthentication:     codes.Unauthenticated,
	CodeUnAuthenticationCore: codes.Unauthenticated,
	CodeForbiddenAccess:      codes.PermissionDenied,
	CodeInternalError:        codes.Internal,
	CodeInvalidRequest:       codes.InvalidArgument,
	CodeInsufficientBalance:  codes.InvalidArgument,
	CodeFailed:               codes.InvalidArgument,
}

type ErrorResponse struct {
	Error            errorCode `json:"error"`
	ErrorDescription string    `json:"error_description,omitempty"`
	ErrorUri         string    `json:"error_uri,omitempty"`
	Details          any       `json:"details,omitempty"`
}

type ErrorCallAPi struct {
	StatusCode int    `json:"-"`
	ErrorCode  string `json:"-"`
	Errors     error  `json:"errors"`
}

func (e ErrorCallAPi) Error() string {
	return e.Errors.Error()
}

type ApiResponse[T any] struct {
	Code    Code   `json:"Code"`
	Status  status `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
	Error   error  `json:"error,omitempty"`
}

func (a *ApiResponse[T]) GetRawJSON() ([]byte, error) {
	if a.Status == FailedResponse || a.Code == CodeSuccessNoData {
		type Alias ApiResponse[T]
		apiResponse := struct {
			Data string `json:"data,omitempty"`
			*Alias
		}{Data: "", Alias: (*Alias)(a)}
		return json.Marshal(apiResponse)
	}
	switch dataByte := any(a.Data).(type) {
	case []byte:
		var data any
		if err := json.Unmarshal(dataByte, &data); err != nil {
			return nil, err
		}
		type Alias ApiResponse[T]
		apiResponse := struct {
			*Alias
			Data any `json:"data,omitempty"`
		}{Data: data, Alias: (*Alias)(a)}
		return json.Marshal(apiResponse)
	}
	return json.Marshal(a)
}

func FailedApiResponse[T any](message string, err error) ApiResponse[T] {
	return ApiResponse[T]{
		Code:    CodeFailed,
		Status:  FailedResponse,
		Message: message,
		Error:   err,
	}
}

func FailedApiResponseWithData[T any](message string, data T, err error) ApiResponse[T] {
	return ApiResponse[T]{
		Code:    CodeFailed,
		Status:  FailedResponse,
		Message: message,
		Data:    data,
		Error:   err,
	}
}

func NotFoundApiResponse[T any](message string, err error) ApiResponse[T] {
	return ApiResponse[T]{
		Code:    CodeNotFound,
		Status:  FailedResponse,
		Message: message,
		Error:   err,
	}
}

func HtxNotFoundApiResponse[T any](message string) ApiResponse[T] {
	return ApiResponse[T]{
		Code:    CodeHtxNotFound,
		Status:  FailedResponse,
		Message: message,
	}
}

func TrxNotFoundApiResponse[T any](message string) ApiResponse[T] {
	return ApiResponse[T]{
		Code:    CodeTrxNotFound,
		Status:  FailedResponse,
		Message: message,
	}
}

func BadRequestApiResponse[T any](message string) ApiResponse[T] {
	return ApiResponse[T]{
		Code:    CodeFailed,
		Status:  FailedResponse,
		Message: message,
		Error:   nil,
	}
}

func BadRequestApiResponseWithError[T any](message string, err error) ApiResponse[T] {
	return ApiResponse[T]{
		Code:    CodeFailed,
		Status:  FailedResponse,
		Message: message,
		Error:   err,
	}
}

func InvalidRequestApiResponse[T any](errMap errs.ErrMap) ApiResponse[T] {
	return ApiResponse[T]{
		Code:    CodeInvalidRequest,
		Status:  FailedResponse,
		Message: "invalid request",
		Error:   errMap,
	}
}

func InternalServerErrorApiResponse[T any](err error) ApiResponse[T] {
	return ApiResponse[T]{
		Code:    CodeInternalError,
		Status:  FailedResponse,
		Message: err.Error(),
	}
}

func InternalServerErrorApiCBS[T any](err error) ApiResponse[T] {
	return ApiResponse[T]{
		Code:    CodeSuccessCBS,
		Status:  FailedResponse,
		Message: err.Error(),
	}
}

func FailedResponseCallAPI[T any](err error) (ApiResponse[T], bool) {
	var errRes ErrorCallAPi
	ok := errors.As(err, &errRes)
	if ok {
		switch {
		case errRes.StatusCode == http.StatusUnauthorized, errRes.StatusCode >= http.StatusInternalServerError:
			return InternalServerErrorApiResponse[T](err), true
		default:
			var errMap errs.ErrMap
			if errors.As(errRes.Errors, &errMap) {
				return FailedErrResponse[T](Code(errRes.ErrorCode), errRes.Errors.Error(), errRes.Errors), true
			}
			return FailedErrResponse[T](Code(errRes.ErrorCode), errRes.Errors.Error()), true
		}
	}
	return ApiResponse[T]{}, false
}

func FailedGetTokenApiResponse[T any](message string) ApiResponse[T] {
	return ApiResponse[T]{
		Code:    CodeUnAuthentication,
		Status:  FailedResponse,
		Message: message,
	}
}

func FailedAuthenticationCore[T any](message string) ApiResponse[T] {
	return ApiResponse[T]{
		Code:    CodeUnAuthenticationCore,
		Status:  FailedResponse,
		Message: message,
	}
}

func FailedFromAnotherResponse[T any, E any](response ApiResponse[E]) ApiResponse[T] {
	return ApiResponse[T]{
		Code:    response.Code,
		Status:  response.Status,
		Message: response.Message,
		Error:   response.Error,
	}
}

func FailedErrResponse[T any](code Code, message string, errs ...error) ApiResponse[T] {
	response := ApiResponse[T]{
		Code:    code,
		Status:  FailedResponse,
		Message: message,
	}
	if len(errs) > 0 {
		response.Error = errs[0]
	}
	return response
}

func SuccessResponseApi[T any](code Code, message string, data T) ApiResponse[T] {
	return ApiResponse[T]{
		Code:    code,
		Status:  SuccessResponse,
		Message: message,
		Data:    data,
	}
}

func SuccessApiResponseWithoutData[T any](message string) ApiResponse[T] {
	return ApiResponse[T]{
		Code:    CodeSuccessNoData,
		Status:  SuccessResponse,
		Message: message,
	}
}

func SuccessApiResponse[T any](message string, data T) ApiResponse[T] {
	return ApiResponse[T]{
		Code:    CodeSuccess,
		Status:  SuccessResponse,
		Message: message,
		Data:    data,
	}
}

func PendingApiResponse[T any](message string, data T) ApiResponse[T] {
	return ApiResponse[T]{
		Code:    CodeSuccess,
		Status:  PendingResponse,
		Message: message,
		Data:    data,
	}
}

func Response[T any](c *fiber.Ctx, response ApiResponse[T]) (err error) {
	raw, err := response.GetRawJSON()
	if err != nil {
		return
	}
	c.Response().SetStatusCode(MapStatusCode[response.Code])
	c.Response().Header.SetContentType(fiber.MIMEApplicationJSON)
	c.Response().SetBodyRaw(raw)
	return
}

func ErrResponseFromGrpc(err error) *ErrResponse {
	statusResponse := status2.Convert(err)
	if len(statusResponse.Details()) > 0 {
		if errorResponse, ok := statusResponse.Details()[0].(*ErrResponse); ok {
			return errorResponse
		}
	}
	return nil
}

func ErrorFromGrpc(err error) error {
	statusResponse := status2.Convert(err)
	if len(statusResponse.Details()) > 0 {
		if errorResponse, ok := statusResponse.Details()[0].(*ErrResponse); ok {
			if errorResponse != nil {
				errorCallAPi := ErrorCallAPi{
					StatusCode: MapStatusCode[Code(errorResponse.Code)],
					ErrorCode:  errorResponse.Code,
				}
				if errorResponse.Detail != nil {
					errMap := errs.ErrMap(errorResponse.Detail.AsInterface().(map[string]any))
					errorCallAPi.Errors = errMap
				} else {
					errorCallAPi.Errors = fmt.Errorf(errorResponse.Message)
				}
				return errorCallAPi
			}
		}
	}
	return err
}

func ResponseErrorGrpc[T any](response ApiResponse[T]) error {
	errResponse := ErrResponse{
		Code:    string(response.Code),
		Message: response.Message,
	}

	var err error
	codeGrpc := MapStatusGrpc[response.Code]
	if codeGrpc == 0 {
		codeGrpc = 3
	}
	statusRes := status2.New(codeGrpc, response.Message)
	if response.Error != nil {
		var errMap errs.ErrMap
		if errors.As(response.Error, &errMap) {
			mapErr := map[string]interface{}(errMap)
			if value, err := structpb.NewValue(mapErr); err != nil {
				return err
			} else {
				errResponse.Detail = value
			}
		}
	}
	statusRes, err = statusRes.WithDetails(&errResponse)
	if err != nil {
		return err
	}
	return statusRes.Err()
}
