package info

import (
	"context"
	"time"
)

type RequestInfo struct {
	UserAgent      string
	UserId         string
	BranchId       string
	DeviceId       string
	Role           string
	IpAddr         string
	Method         string
	Path           string
	ReqTime        time.Time
	IdempotencyKey string
	RequestId      string
}

type requestInfo struct{}

func NewContextWithRequestInfo(ctx context.Context, info RequestInfo) context.Context {
	return context.WithValue(ctx, requestInfo{}, info)
}

func GetRequestInfo(ctx context.Context) RequestInfo {
	info, _ := ctx.Value(requestInfo{}).(RequestInfo)
	return info
}
