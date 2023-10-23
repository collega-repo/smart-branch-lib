package commons

import (
	"errors"
	"fmt"
	"github.com/collega-repo/smart-branch-lib/commons/errs"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
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

var MapCode = map[code]map[int]string{
	CodeSuccess: {
		http.StatusOK: "",
	},
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
			return ApiResponse[T]{
				Code:    CodeInternalError,
				Status:  FailedResponse,
				Message: fmt.Sprintf(`invalid call api: %s`, err.Error()),
			}, true
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
