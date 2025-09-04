// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"parkin-ai-system/api/wallet_transaction"
)

type (
	IWalletTransaction interface {
		WalletTransactionAdd(ctx context.Context, req *wallet_transaction.WalletTransactionAddReq) (*wallet_transaction.WalletTransactionAddRes, error)
		WalletTransactionList(ctx context.Context, req *wallet_transaction.WalletTransactionListReq) (*wallet_transaction.WalletTransactionListRes, error)
	}
)

var (
	localWalletTransaction IWalletTransaction
)

func WalletTransaction() IWalletTransaction {
	if localWalletTransaction == nil {
		panic("implement not found for interface IWalletTransaction, forgot register?")
	}
	return localWalletTransaction
}

func RegisterWalletTransaction(i IWalletTransaction) {
	localWalletTransaction = i
}
