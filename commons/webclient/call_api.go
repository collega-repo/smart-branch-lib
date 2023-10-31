package webclient

import (
	"bytes"
	"context"
	"errors"
	"github.com/collega-repo/smart-branch-lib/commons"
	"github.com/collega-repo/smart-branch-lib/commons/errs"
	"github.com/collega-repo/smart-branch-lib/commons/info"
	"github.com/collega-repo/smart-branch-lib/configs"
	"github.com/goccy/go-json"
	"io"
	"net/http"
	"regexp"
	"syscall"
	"time"
)

var passwordRegex = regexp.MustCompile(`"password":\s*"([^"]+)"`)

func CallRestApi[T, P any](ctx context.Context, path string, method string, payLoad P, listReqInfo ...info.RequestInfo) (status int, result T, err error) {
	status = 500
	var byteJson []byte
	if any(payLoad) != nil {
		byteJson, err = json.Marshal(payLoad)
		if err != nil {
			return status, result, err
		}
	}

	var reqInfo info.RequestInfo
	if len(listReqInfo) == 0 {
		reqInfo = info.GetRequestInfo(ctx)
	} else {
		reqInfo = listReqInfo[0]
	}

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

	logEven := configs.LogEvent{}
	if reqInfo.IdempotencyKey != "" {
		req.Header.Set(commons.Configs.App.Cache.Key, reqInfo.IdempotencyKey)
		logEven.IdempotencyKey = reqInfo.IdempotencyKey
	}

	if reqInfo.RequestId != "" {
		req.Header.Set(`X-Request-ID`, reqInfo.RequestId)
		logEven.RequestId = reqInfo.RequestId
	}

	now := time.Now()
	logReq := configs.NewLoggerReqRes(`external_api_request`, now, false)
	logRes := configs.NewLoggerReqRes(`external_api_response`, now, true)

	logEven.Path = path
	logEven.Method = method
	logEven.BodyRequest = byteJson
	logEven.SendRequest(logReq)

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		if commons.Configs.Core.LogDebug {
			logEven.Error = err
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
			logEven.Error = err
			logEven.SendResponse(logRes)
		}
		return res.StatusCode, result, errs.ErrApiNotFound
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		if commons.Configs.Core.LogDebug {
			logEven.Error = err
			logEven.SendResponse(logRes)
		}
		return status, result, err
	}
	defer func() {
		_ = res.Close
	}()

	if body != nil {
		if commons.Configs.Core.LogDebug {
			logEven.BodyResponse = body
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

func CallPostServiceApi[T, P any](ctx context.Context, path string, payload P, reqInfo ...info.RequestInfo) (status int, result T, err error) {
	return CallRestApi[T](ctx, path, http.MethodPost, payload, reqInfo...)
}
