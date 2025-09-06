// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// WalletTransactions is the golang structure for table wallet_transactions.
type WalletTransactions struct {
	Id             int64       `json:"id"             orm:"id"               description:""`
	UserId         int64       `json:"userId"         orm:"user_id"          description:""`
	Amount         float64     `json:"amount"         orm:"amount"           description:""`
	Type           string      `json:"type"           orm:"type"             description:""`
	Description    string      `json:"description"    orm:"description"      description:""`
	RelatedOrderId int64       `json:"relatedOrderId" orm:"related_order_id" description:""`
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"       description:""`
}

type WalletTransactionAddReq struct {
	UserId         int64   `json:"userId"`
	Amount         float64 `json:"amount"`
	Type           string  `json:"type"`
	Description    string  `json:"description"`
	RelatedOrderId int64   `json:"relatedOrderId"`
}

type WalletTransactionAddRes struct {
	Id int64 `json:"id"`
}

type WalletTransactionListReq struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Page        int    `json:"page"`
	PageSize    int    `json:"pageSize"`
}

type WalletTransactionItem struct {
	Id             int64   `json:"id"`
	UserId         int64   `json:"user_id"`
	Username       string  `json:"username"`
	Amount         float64 `json:"amount"`
	Type           string  `json:"type"`
	Description    string  `json:"description"`
	RelatedOrderId int64   `json:"related_order_id"`
	LicensePlate   string  `json:"license_plate"`
	ServiceType    string  `json:"service_type"`
	CreatedAt      string  `json:"created_at"`
}

type WalletTransactionListRes struct {
	List  []WalletTransactionItem `json:"list"`
	Total int                     `json:"total"`
}

type WalletTransactionGetReq struct {
	Id int64 `json:"id"`
}

type WalletTransactionGetRes struct {
	Transaction WalletTransactionItem `json:"transaction"`
}