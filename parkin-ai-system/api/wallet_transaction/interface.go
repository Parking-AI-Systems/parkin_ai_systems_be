package wallet_transaction

import "context"

type IWalletTransaction interface {
	WalletTransactionAdd(ctx context.Context, req *WalletTransactionAddReq) (*WalletTransactionAddRes, error)
	WalletTransactionList(ctx context.Context, req *WalletTransactionListReq) (*WalletTransactionListRes, error)
}
