// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// WalletTransactions is the golang structure of table wallet_transactions for DAO operations like Where/Data.
type WalletTransactions struct {
	g.Meta         `orm:"table:wallet_transactions, do:true"`
	Id             interface{} //
	UserId         interface{} //
	Amount         interface{} //
	Type           interface{} //
	Description    interface{} //
	RelatedOrderId interface{} //
	CreatedAt      *gtime.Time //
}
