package wallet_transaction

import (
	"context"
	"parkin-ai-system/api/wallet_transaction"
	"parkin-ai-system/internal/service"
	"github.com/gogf/gf/v2/errors/gerror"
)

type ControllerWalletTransaction struct{}

func NewWalletTransaction() *ControllerWalletTransaction {
	return &ControllerWalletTransaction{}
}

func (c *ControllerWalletTransaction) WalletTransactionAdd(ctx context.Context, req *wallet_transaction.WalletTransactionAddReq) (res *wallet_transaction.WalletTransactionAddRes, err error) {
	res, err = service.WalletTransaction().WalletTransactionAdd(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}

func (c *ControllerWalletTransaction) WalletTransactionList(ctx context.Context, req *wallet_transaction.WalletTransactionListReq) (res *wallet_transaction.WalletTransactionListRes, err error) {
	res, err = service.WalletTransaction().WalletTransactionList(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}
