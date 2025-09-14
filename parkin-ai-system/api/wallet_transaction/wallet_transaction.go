// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package wallet_transaction

import (
	"context"

	"parkin-ai-system/api/wallet_transaction/wallet_transaction"
)

type IWalletTransactionWallet_transaction interface {
	WalletTransactionAdd(ctx context.Context, req *wallet_transaction.WalletTransactionAddReq) (res *wallet_transaction.WalletTransactionAddRes, err error)
	WalletTransactionList(ctx context.Context, req *wallet_transaction.WalletTransactionListReq) (res *wallet_transaction.WalletTransactionListRes, err error)
	WalletTransactionGet(ctx context.Context, req *wallet_transaction.WalletTransactionGetReq) (res *wallet_transaction.WalletTransactionGetRes, err error)
	MyWalletTransactionGet(ctx context.Context, req *wallet_transaction.MyWalletTransactionGetReq) (res *wallet_transaction.MyWalletTransactionGetRes, err error)
}
