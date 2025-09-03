package service

import (
	"parkin-ai-system/api/wallet_transaction"
)

type IWalletTransaction interface {
	wallet_transaction.IWalletTransaction
}

var localWalletTransaction IWalletTransaction

func WalletTransaction() IWalletTransaction {
	if localWalletTransaction == nil {
		panic("implement not found for interface IWalletTransaction, forgot register?")
	}
	return localWalletTransaction
}

func RegisterWalletTransaction(i IWalletTransaction) {
	localWalletTransaction = i
}
