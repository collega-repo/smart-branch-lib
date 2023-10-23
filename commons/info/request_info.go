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

func GetRequestUserId(ctx context.Context) string {
	info, _ := ctx.Value(requestInfo{}).(RequestInfo)
	return info.UserId
}

func GetRequestIpAddr(ctx context.Context) string {
	info, _ := ctx.Value(requestInfo{}).(RequestInfo)
	if info.IpAddr == "" {
		info.IpAddr = "0:0:0:0:0:0:0:1"
	}
	return info.IpAddr
}

func GetRequestIdempotencyKey(ctx context.Context) string {
	info, _ := ctx.Value(requestInfo{}).(RequestInfo)
	return info.IdempotencyKey
}

func GetRequestCfgSys(ctx context.Context) dto.CfgSys {
	info, _ := ctx.Value(requestInfo{}).(RequestInfo)
	return info.CfgSys
}
