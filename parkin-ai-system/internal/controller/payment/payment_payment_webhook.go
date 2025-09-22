package payment

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"parkin-ai-system/api/payment/payment"
)

func (c *ControllerPayment) Webhook(ctx context.Context, req *payment.WebhookReq) (res *payment.WebhookRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
