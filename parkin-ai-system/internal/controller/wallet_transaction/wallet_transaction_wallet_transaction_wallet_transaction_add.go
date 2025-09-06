package wallet_transaction

import (
	"context"

	"parkin-ai-system/api/wallet_transaction/wallet_transaction"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerWallet_transaction) WalletTransactionAdd(ctx context.Context, req *wallet_transaction.WalletTransactionAddReq) (res *wallet_transaction.WalletTransactionAddRes, err error) {
	// Map API request to entity request
	input := &entity.WalletTransactionAddReq{
		UserId:         req.UserId,
		Amount:         req.Amount,
		Type:           req.Type,
		Description:    req.Description,
		RelatedOrderId: req.RelatedOrderId,
	}

	// Call service
	addRes, err := service.WalletTransaction().WalletTransactionAdd(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &wallet_transaction.WalletTransactionAddRes{
		Id: addRes.Id,
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
