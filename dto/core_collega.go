package dto

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/collega-repo/smart-branch-lib/commons"
	"github.com/collega-repo/smart-branch-lib/commons/mapper"
	"github.com/shopspring/decimal"
	"time"
)

/*
khusus untuk core banking collega
*/

type GatewayRes[T any] struct {
	RCode      string          `json:"rCode"`
	StatusId   int             `json:"statusId"`
	Message    string          `json:"message"`
	Result     T               `json:"result,omitempty"`
	SaldoAwal  decimal.Decimal `json:"saldoAwal,omitempty"`
	SaldoAkhir decimal.Decimal `json:"saldoAkhir,omitempty"`
	RolesMap   []mapper.Map    `json:"rolesMap,omitempty"`
	UserMap    mapper.Map      `json:"userMap,omitempty"`
	BranchId   string          `json:"branchId,omitempty"`
	ApplId     string          `json:"applId,omitempty"`
	AccFound   int             `json:"accFound,omitempty"`
	AccName    string          `json:"accName,omitempty"`
	AccNbr     string          `json:"accNbr,omitempty"`
	AccStatus  int             `json:"accStatus,omitempty"`
}

type GatewayReq struct {
	AuthKey   string `json:"authKey"`
	ReqId     string `json:"reqId"`
	TxDate    string `json:"txDate,omitempty"`
	TxHour    string `json:"txHour,omitempty"`
	UserGtw   string `json:"userGtw"`
	ChannelId string `json:"channelId"`
}

func NewGatewayReq(reqId, ipSource string) GatewayReq {

	now := time.Now()
	payload := fmt.Sprintf(`%s%s%s`, reqId, ipSource, now.Format(`2006-01-0215:04:05`))
	hashMac := hmac.New(sha1.New, []byte(commons.Configs.Core.SecretKey))
	hashMac.Write([]byte(payload))
	authKey := hex.EncodeToString(hashMac.Sum(nil))

	return GatewayReq{
		ReqId:     reqId,
		TxDate:    now.Format(`20060102`),
		TxHour:    now.Format(`150405`),
		UserGtw:   commons.Configs.Core.Username,
		ChannelId: commons.Configs.Core.ChannelId,
		AuthKey:   authKey,
	}
}

type GatewayReqInq struct {
	GatewayReq
	IdNbr     string `json:"idNbr,omitempty"`
	Idnbr     string `json:"idnbr,omitempty"`
	CifId     string `json:"cifId,omitempty"`
	AccNbr    string `json:"accNbr,omitempty"`
	Accnbr    string `json:"accnbr,omitempty"`
	StartDate string `json:"startDate,omitempty"`
	EndDate   string `json:"endDate,omitempty"`
	FlgSaldo  int    `json:"flgSaldo,omitempty"`
	ApplId    string `json:"applId,omitempty"`
}

type GatewayReqUser struct {
	GatewayReq
	UserId   string `json:"userId,omitempty"`
	Password string `json:"password,omitempty"`
	FlgRole  int    `json:"flgRole"`
	FlgUser  int    `json:"flgUser"`
}

type GatewayReqTrx struct {
	GatewayReq
	Date        string          `json:"date"`
	DateRk      string          `json:"date_rk"`
	CorpId      string          `json:"corpId"`
	ProdId      string          `json:"prodId"`
	BranchId    string          `json:"branchId"`
	TxCcy       string          `json:"txCcy"`
	NbrOfAcc    int             `json:"nbrOfAcc"`
	TotalAmount decimal.Decimal `json:"totalAmount"`
	ProsesId    string          `json:"prosesId"`
	UserId      string          `json:"userId"`
	SpvId       string          `json:"spvId"`
	RevSts      int             `json:"revSts"`
	TxType      string          `json:"txType"`
	RefAcc      string          `json:"refAcc"`
	Param       string          `json:"param"`
}

type TrxModel struct {
	TxId   string          `json:"txId"`
	TxMsg  string          `json:"txMsg"`
	AccNbr string          `json:"accNbr"`
	DbCr   int             `json:"dbCr"`
	TxAmt  decimal.Decimal `json:"txAmt"`
	TxCode string          `json:"txCode"`
}
