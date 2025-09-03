package wallet_transaction

import (
	"context"
	"parkin-ai-system/api/wallet_transaction"
	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/model/do"
	"parkin-ai-system/internal/model/entity"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/os/gtime"
	"parkin-ai-system/internal/service"
)

type sWalletTransaction struct{}

func init() {
	service.RegisterWalletTransaction(&sWalletTransaction{})
}

func (s *sWalletTransaction) WalletTransactionAdd(ctx context.Context, req *wallet_transaction.WalletTransactionAddReq) (*wallet_transaction.WalletTransactionAddRes, error) {
	userId := ctx.Value("user_id")
	tran := do.WalletTransactions{}
	gconv.Struct(req, &tran)
	tran.UserId = userId
	tran.CreatedAt = gtime.Now()
	lastId, err := dao.WalletTransactions.Ctx(ctx).Data(tran).InsertAndGetId()
	if err != nil {
		return nil, err
	}
	return &wallet_transaction.WalletTransactionAddRes{Id: gconv.Int64(lastId)}, nil
}

func (s *sWalletTransaction) WalletTransactionList(ctx context.Context, req *wallet_transaction.WalletTransactionListReq) (*wallet_transaction.WalletTransactionListRes, error) {
	userId := ctx.Value("user_id")
	var trans []entity.WalletTransactions
	err := dao.WalletTransactions.Ctx(ctx).Where("user_id", userId).Order("id desc").Scan(&trans)
	if err != nil {
		return nil, err
	}
	var list []wallet_transaction.WalletTransactionItem
	for _, t := range trans {
		item := wallet_transaction.WalletTransactionItem{}
		gconv.Struct(t, &item)
		item.CreatedAt = t.CreatedAt.Format("2006-01-02 15:04:05")
		list = append(list, item)
	}
	return &wallet_transaction.WalletTransactionListRes{List: list}, nil
}
