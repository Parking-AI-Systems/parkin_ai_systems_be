package wallet_transaction

import (
	"context"

	"parkin-ai-system/api/wallet_transaction/wallet_transaction"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerWallet_transaction) WalletTransactionList(ctx context.Context, req *wallet_transaction.WalletTransactionListReq) (res *wallet_transaction.WalletTransactionListRes, err error) {
	// Map API request to entity request
	input := &entity.WalletTransactionListReq{
		Type:        req.Type,
		Description: req.Description,
		Page:        req.Page,
		PageSize:    req.PageSize,
	}

	// Call service
	listRes, err := service.WalletTransaction().WalletTransactionList(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity list to API response
	res = &wallet_transaction.WalletTransactionListRes{
		List:  make([]wallet_transaction.WalletTransactionItem, 0, len(listRes.List)),
		Total: listRes.Total,
	}
	for _, item := range listRes.List {
		res.List = append(res.List, wallet_transaction.WalletTransactionItem{
			Id:             item.Id,
			UserId:         item.UserId,
			Amount:         item.Amount,
			Type:           item.Type,
			Description:    item.Description,
			RelatedOrderId: item.RelatedOrderId,
			CreatedAt:      item.CreatedAt,
		})
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
