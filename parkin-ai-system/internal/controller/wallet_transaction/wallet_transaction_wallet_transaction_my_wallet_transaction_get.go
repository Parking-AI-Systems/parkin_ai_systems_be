package wallet_transaction

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"parkin-ai-system/api/wallet_transaction/wallet_transaction"
)

func (c *ControllerWallet_transaction) MyWalletTransactionGet(ctx context.Context, req *wallet_transaction.MyWalletTransactionGetReq) (res *wallet_transaction.MyWalletTransactionGetRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
