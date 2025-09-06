// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"parkin-ai-system/internal/model/entity"
)

type (
	IWalletTransaction interface {
		WalletTransactionAdd(ctx context.Context, req *entity.WalletTransactionAddReq) (*entity.WalletTransactionAddRes, error)
		WalletTransactionList(ctx context.Context, req *entity.WalletTransactionListReq) (*entity.WalletTransactionListRes, error)
		WalletTransactionGet(ctx context.Context, req *entity.WalletTransactionGetReq) (*entity.WalletTransactionItem, error)
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
