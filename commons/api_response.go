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
	code      string
	errorCode string
)

const (
	CodeSuccess                code      = "00"
	CodeSuccessCBS             code      = "01"
	CodeSuccessPending         code      = "02"
	CodeSuccessPendingApproval code      = "03"
	CodeSuccessNoData          code      = "04"
	CodeNotFoundCore           code      = "84"
	CodeUnAuthenticationCore   code      = "87"
	CodeInsufficientBalance    code      = "90"
	CodeHtxNotFound            code      = "91"
	CodeTrxNotFound            code      = "92"
	CodeInternalError          code      = "93"
	CodeNotFound               code      = "94"
	CodeIdempotencyFailed      code      = "95"
	CodeForbiddenAccess        code      = "96"
	CodeUnAuthentication       code      = "97"
	CodeInvalidRequest         code      = "98"
	CodeFailed                 code      = "99"
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

var MapStatusCode = map[code]int{
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
	CodeFailed:                 http.StatusBadRequest,
}

var MapStatusRestToGrpc = map[code]codes.Code{
	CodeSuccess:              codes.OK,
	CodeNotFound:             codes.NotFound,
	CodeNotFoundCore:         codes.NotFound,
	CodeUnAuthentication:     codes.Unauthenticated,
	CodeUnAuthenticationCore: codes.Unauthenticated,
	CodeForbiddenAccess:      codes.PermissionDenied,
	CodeInternalError:        codes.Internal,
	CodeInvalidRequest:       codes.InvalidArgument,
	CodeFailed:               codes.InvalidArgument,
}

var MapStatusGrpcToRest = map[codes.Code]code{
	codes.OK:               CodeSuccess,
	codes.NotFound:         CodeNotFound,
	codes.Unauthenticated:  CodeUnAuthentication,
	codes.PermissionDenied: CodeForbiddenAccess,
	codes.Internal:         CodeInternalError,
	codes.InvalidArgument:  CodeFailed,
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
	Code    code   `json:"code"`
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
		case errRes.StatusCode == http.StatusNotFound:
			return NotFoundApiResponse[T](err.Error(), nil), true
		default:
			return FailedApiResponse[T](err.Error(), nil), true
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

func FailedErrResponse[T any](code code, message string) ApiResponse[T] {
	return ApiResponse[T]{
		Code:    code,
		Status:  FailedResponse,
		Message: message,
	}
}

func SuccessResponseApi[T any](code code, message string, data T) ApiResponse[T] {
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
				switch errorResponse.Code {
				case string(CodeUnAuthenticationCore), string(CodeUnAuthentication):
					return errs.ErrAuthFailed
				case string(CodeNotFoundCore), string(CodeNotFound):
					return errs.ErrRecordNotFound
				case string(CodeInsufficientBalance):
					return errs.ErrInsufficientBalance
				default:
					return ErrorCallAPi{
						StatusCode: MapStatusCode[MapStatusGrpcToRest[statusResponse.Code()]],
						ErrorCode:  errorResponse.Code,
						Errors:     fmt.Errorf(errorResponse.Message),
					}
				}
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
	statusRes := status2.New(MapStatusRestToGrpc[response.Code], response.Message)
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
