package dto

import (
	"github.com/shopspring/decimal"
	"time"
)

type Url struct {
	Code string `json:"code" validate:"required,max=6"`
	Path string `json:"path" validate:"required"`
	Desc string `json:"desc,omitempty" validate:"required,min=3,max=100"`
}

type Menu struct {
	MenuId     string `json:"menuId,omitempty"`
	Name       string `json:"name"`
	Seq        int    `json:"seq"`
	PathFE     string `json:"pathFE"`
	MenuParent string `json:"menuParent,omitempty"`
	MenuType   int    `json:"menuType"`
	MenuChild  []Menu `json:"menuChild,omitempty"`
}

type UserMap struct {
	AccessToken string       `json:"accessToken,omitempty"`
	UserId      string       `json:"userId"`
	BranchId    string       `json:"branchId"`
	BranchNm    string       `json:"branchNm"`
	FullName    string       `json:"fullName"`
	RoleId      string       `json:"roleId,omitempty"`
	RoleNm      string       `json:"roleNm"`
	IaNbr       string       `json:"iaNbr,omitempty"`
	IpAddr      string       `json:"ipAddr"`
	Device      CfgDevice    `json:"device,omitempty"`
	Application string       `json:"application,omitempty"`
	OpenDate    time.Time    `json:"openDate"`
	UserLimit   UserLimit    `json:"userLimit,omitempty"`
	StartedAt   time.Time    `json:"startedAt,omitempty"`
	Accesses    []RoleAccess `json:"accesses,omitempty"`
	Menus       []Menu       `json:"menus,omitempty"`
}

type UserLimit struct {
	UserId   string          `json:"userId"`
	CcyId    string          `json:"ccyId"`
	BranchId string          `json:"branchId"`
	MxAccStr decimal.Decimal `json:"mxAccStr"`
	MxClAcWd decimal.Decimal `json:"mxClAcWd"`
	MxCshStr decimal.Decimal `json:"mxCshStr"`
	MaxCrOb  decimal.Decimal `json:"mxCrOb"`
	MaxDrOb  decimal.Decimal `json:"mxDbOb"`
	MxClrStr decimal.Decimal `json:"mxClrStr"`
	MxDrObAc decimal.Decimal `json:"mxDrObAc"`
	MxCsHwDr decimal.Decimal `json:"mxCsHwDr"`
}

type RoleAccess struct {
	Method string `json:"method"`
	Urls   []Url  `json:"urls"`
}
