package wallet_transaction

import (
	"github.com/gogf/gf/v2/frame/g"
)

type WalletTransactionAddReq struct {
	g.Meta      `path:"/wallet-transactions" method:"post" tags:"WalletTransaction" summary:"Tạo giao dịch ví" security:"BearerAuth"`
	Amount      float64 `json:"amount" v:"required"`
	Type        string  `json:"type" v:"required"`
	Description string  `json:"description"`
	RelatedOrderId int64 `json:"related_order_id"`
}

type WalletTransactionAddRes struct {
	Id int64 `json:"id"`
}

type WalletTransactionListReq struct {
	g.Meta `path:"/wallet-transactions" method:"get" tags:"WalletTransaction" summary:"Danh sách giao dịch ví" security:"BearerAuth"`
}

type WalletTransactionListRes struct {
	List []WalletTransactionItem `json:"list"`
}

type WalletTransactionItem struct {
	Id            int64   `json:"id"`
	Amount        float64 `json:"amount"`
	Type          string  `json:"type"`
	Description   string  `json:"description"`
	RelatedOrderId int64  `json:"related_order_id"`
	CreatedAt     string  `json:"created_at"`
}
