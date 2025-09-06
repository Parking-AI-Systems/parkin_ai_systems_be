package wallet_transaction

import "github.com/gogf/gf/v2/frame/g"

type WalletTransactionAddReq struct {
	g.Meta         `path:"/wallet-transactions" tags:"Wallet Transaction" method:"POST" summary:"Add a wallet transaction" description:"Creates a new wallet transaction for a user. Admins can add for any user; users can add for themselves. Related order ID can reference parking or service orders." middleware:"middleware.Auth"`
	UserId         int64   `json:"user_id" v:"required|min:1#User ID is required|User ID must be positive"`
	Amount         float64 `json:"amount" v:"required#Amount is required"`
	Type           string  `json:"type" v:"required#Transaction type is required"`
	Description    string  `json:"description" v:"required|length:1,255#Description is required|Description must be 1-255 characters"`
	RelatedOrderId int64   `json:"related_order_id" v:"min:0#Related order ID must be non-negative"`
}

type WalletTransactionAddRes struct {
	Id int64 `json:"id"`
}

type WalletTransactionListReq struct {
	g.Meta      `path:"/wallet-transactions" tags:"Wallet Transaction" method:"GET" summary:"List wallet transactions" description:"Retrieves a paginated list of wallet transactions for the authenticated user or all transactions for admins." middleware:"middleware.Auth"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Page        int    `json:"page" v:"min:1#Page must be at least 1"`
	PageSize    int    `json:"page_size" v:"min:1|max:100#Page size must be between 1 and 100"`
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
	g.Meta `path:"/wallet-transactions/:id" tags:"Wallet Transaction" method:"GET" summary:"Get wallet transaction details" description:"Retrieves details of a specific wallet transaction by ID." middleware:"middleware.Auth"`
	Id     int64 `json:"id" v:"required|min:1#Transaction ID is required|Transaction ID must be positive"`
}

type WalletTransactionGetRes struct {
	Transaction WalletTransactionItem `json:"transaction"`
}
