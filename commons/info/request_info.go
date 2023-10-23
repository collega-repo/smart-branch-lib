package info

import (
	"context"
	"github.com/collega-repo/smart-branch-lib/dto"
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
	ApplicationId  string
	OauthClient    dto.Oauth2Client
	OauthToken     dto.Oauth2Token
	CfgSys         dto.CfgSys
	ListApplId     []string
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

func GetRequestBranchId(ctx context.Context) string {
	info, _ := ctx.Value(requestInfo{}).(RequestInfo)
	return info.BranchId
}

func GetRequestApplId(ctx context.Context) string {
	info, _ := ctx.Value(requestInfo{}).(RequestInfo)
	return info.ApplicationId
}

func GetRequestOauthClient(ctx context.Context) dto.Oauth2Client {
	info, _ := ctx.Value(requestInfo{}).(RequestInfo)
	return info.OauthClient
}

func GetRequestOauthToken(ctx context.Context) dto.Oauth2Token {
	info, _ := ctx.Value(requestInfo{}).(RequestInfo)
	return info.OauthToken
}

func GetRequestCfgSys(ctx context.Context) dto.CfgSys {
	info, _ := ctx.Value(requestInfo{}).(RequestInfo)
	return info.CfgSys
}

func GetRequestCfgAppl(ctx context.Context) []string {
	info, _ := ctx.Value(requestInfo{}).(RequestInfo)
	return info.ListApplId
}

func GetRequestDeviceId(ctx context.Context) string {
	info, _ := ctx.Value(requestInfo{}).(RequestInfo)
	return info.DeviceId
}
