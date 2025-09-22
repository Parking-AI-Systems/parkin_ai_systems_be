package payment

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"parkin-ai-system/api/payment/payment"
)

func (c *ControllerPayment) CreatePaymentLink(ctx context.Context, req *payment.CreatePaymentLinkReq) (res *payment.CreatePaymentLinkRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
