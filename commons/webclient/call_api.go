package webclient

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/collega-repo/smart-branch-lib/commons"
	"github.com/collega-repo/smart-branch-lib/commons/errs"
	"github.com/collega-repo/smart-branch-lib/commons/info"
	"github.com/collega-repo/smart-branch-lib/configs"
	"github.com/collega-repo/smart-branch-lib/dto"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"regexp"
	"syscall"
	"time"
)

type logEven struct {
	idempotencyKey string
	bodyRequest    []byte
	bodyResponse   []byte
	err            error
	msg            string
	startTime      time.Time
	method         string
	url            string
}

func (l logEven) SendRequest(logger zerolog.Logger) {
	var eventLog *zerolog.Event
	if l.err != nil {
		eventLog = logger.Err(l.err)
	} else {
		eventLog = logger.Debug()
	}
	if string(l.bodyRequest) != "" {
		eventLog.RawJSON("reqBody", l.bodyRequest)
	}
	eventLog.
		Str(`idempotencyKey`, l.idempotencyKey).
		Str(`startTime`, l.startTime.Format("2006-01-02 15:04:05.000")).
		Str(`method`, l.method).
		Str(`url`, l.url)

	go func() {
		if l.msg != "" {
			eventLog.Msg(l.msg)
		} else {
			eventLog.Send()
		}
	}()
}

func (l logEven) SendResponse(logger zerolog.Logger) {
	var eventLog *zerolog.Event
	if l.err != nil {
		eventLog = logger.Err(l.err)
	} else {
		eventLog = logger.Info()
	}
	if string(l.bodyResponse) != "" {
		eventLog.RawJSON(`resBody`, l.bodyResponse)
	}
	if string(l.bodyRequest) != "" {
		eventLog.RawJSON("reqBody", l.bodyRequest)
	}
	eventLog.Str(`idempotencyKey`, l.idempotencyKey).
		Dur(`duration`, time.Since(l.startTime)).
		Str(`startTime`, l.startTime.Format("2006-01-02 15:04:05.000")).
		Str(`endTime`, time.Now().Format("2006-01-02 15:04:05.000")).
		Str(`method`, l.method).
		Str(`url`, l.url)

	go func() {
		if l.msg != "" {
			eventLog.Msg(l.msg)
		} else {
			eventLog.Send()
		}
	}()
}

var passwordRegex = regexp.MustCompile(`"password":\s*"([^"]+)"`)

func CallRestApi[T, P any](ctx context.Context, path string, method string, payLoad P) (status int, result T, err error) {
	status = 500
	var byteJson []byte
	if any(payLoad) != nil {
		byteJson, err = json.Marshal(payLoad)
		if err != nil {
			return status, result, err
		}
	}

	reqInfo := info.GetRequestInfo(ctx)

	var buffer *bytes.Buffer
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		if byteJson != nil {
			buffer = bytes.NewBuffer(byteJson)
			byteJson = []byte(passwordRegex.ReplaceAllString(string(byteJson), `"password":"********"`))
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, path, buffer)
	if err != nil {
		return status, result, err
	}
	req.Header.Set(`Content-Type`, `application/json`)
	if commons.Configs.Core.IsForward {
		req.Header.Set(`X-Forwarded-For`, reqInfo.IpAddr)
	}

	logEven := logEven{}
	if reqInfo.IdempotencyKey != "" {
		req.Header.Set(commons.Configs.App.Cache.Key, reqInfo.IdempotencyKey)
		logEven.idempotencyKey = reqInfo.IdempotencyKey
	}

	now := time.Now()
	logReq := configs.NewLoggerReqRes(`external_api_request`, now, false)
	logRes := configs.NewLoggerReqRes(`external_api_response`, now, true)

	logEven.method = method
	logEven.url = path
	logEven.bodyRequest = byteJson
	logEven.SendRequest(logReq)

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		if commons.Configs.Core.LogDebug {
			logEven.err = err
			logEven.SendResponse(logRes)
		}
		var aa syscall.Errno
		if errors.As(err, &aa) {
			return status, result, errs.ErrpApiRequest
		}
		return status, result, err
	}

	if res.Header.Get(`Content-Type`) == "text/html;charset=utf-8" {
		if commons.Configs.Core.LogDebug {
			logEven.err = err
			logEven.SendResponse(logRes)
		}
		return res.StatusCode, result, errs.ErrApiNotFound
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		if commons.Configs.Core.LogDebug {
			logEven.err = err
			logEven.SendResponse(logRes)
		}
		return status, result, err
	}
	defer func() {
		_ = res.Close
	}()

	if body != nil {
		if commons.Configs.Core.LogDebug {
			logEven.bodyResponse = body
			logEven.SendResponse(logRes)
		}

		err = json.Unmarshal(body, &result)
		if err != nil {
			configs.Loggers.Err(err).Msg(`failed read body response external API`)
			return status, result, err
		}
	}

	status = res.StatusCode
	return status, result, err
}

func CallPostServiceApi[T, P any](ctx context.Context, path string, payload P) (status int, result T, err error) {
	return CallRestApi[T](ctx, path, http.MethodPost, payload)
}

// CallGatewayApiCBS
// khusus jika core banking adalah olibs724 (collega)
func CallGatewayApiCBS[T, P any](ctx context.Context, pathApi string, req P) (dto.GatewayRes[T], error) {
	statusCode, gatewayRes, err := CallPostServiceApi[dto.GatewayRes[T]](ctx, pathApi, req)
	if err != nil {
		if errors.Is(err, errs.ErrApiNotFound) || errors.Is(err, errs.ErrpApiRequest) {
			return gatewayRes, commons.ErrorCallAPi{
				StatusCode: statusCode,
				ErrorCode:  gatewayRes.RCode,
				Errors:     fmt.Errorf(`Gagal terhubung dengan Core Banking`),
			}
		}
		return gatewayRes, commons.ErrorCallAPi{
			StatusCode: statusCode,
			ErrorCode:  gatewayRes.RCode,
			Errors:     err,
		}
	}

	if gatewayRes.StatusId == 0 && (statusCode >= 400 && statusCode < 500) {
		return gatewayRes, commons.ErrorCallAPi{
			StatusCode: statusCode,
			ErrorCode:  gatewayRes.RCode,
			Errors:     fmt.Errorf(gatewayRes.Message),
		}
	}

	if gatewayRes.StatusId == 0 && (statusCode >= 200 && statusCode < 400) {
		return gatewayRes, commons.ErrorCallAPi{
			StatusCode: http.StatusBadRequest,
			ErrorCode:  gatewayRes.RCode,
			Errors:     fmt.Errorf(gatewayRes.Message),
		}
	}

	if gatewayRes.StatusId == 0 {
		return gatewayRes, commons.ErrorCallAPi{
			StatusCode: statusCode,
			ErrorCode:  gatewayRes.RCode,
			Errors:     fmt.Errorf(gatewayRes.Message),
		}
	}
	return gatewayRes, nil
}
