package wallet_transaction

import (
	"context"

	"parkin-ai-system/api/wallet_transaction/wallet_transaction"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

func (c *ControllerWallet_transaction) WalletTransactionGet(ctx context.Context, req *wallet_transaction.WalletTransactionGetReq) (res *wallet_transaction.WalletTransactionGetRes, err error) {
	// Map API request to entity request
	input := &entity.WalletTransactionGetReq{
		Id: req.Id,
	}

	// Call service
	tx, err := service.WalletTransaction().WalletTransactionGet(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &wallet_transaction.WalletTransactionGetRes{
		Transaction: wallet_transaction.WalletTransactionItem{
			Id:             tx.Id,
			UserId:         tx.UserId,
			Amount:         tx.Amount,
			Type:           tx.Type,
			Description:    tx.Description,
			RelatedOrderId: tx.RelatedOrderId,
			CreatedAt:      tx.CreatedAt,
		},
	}
	return res, nil
}
