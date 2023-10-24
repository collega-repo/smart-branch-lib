package info

import (
	"context"
	"github.com/collega-repo/smart-branch-lib/dto"
	"time"
)

type RequestInfo struct {
	UserAgent      string
	UserId         string
	IpAddr         string
	Method         string
	Path           string
	ReqTime        time.Time
	IdempotencyKey string
	CfgSys         dto.CfgSys
}

type requestInfo struct{}

func NewContextWithRequestInfo(ctx context.Context, info RequestInfo) context.Context {
	return context.WithValue(ctx, requestInfo{}, info)
}

func GetRequestInfo(ctx context.Context) RequestInfo {
	info, _ := ctx.Value(requestInfo{}).(RequestInfo)
	return info
}
