package payment

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"parkin-ai-system/api/payment/payment"
	"parkin-ai-system/internal/service"
)

func (c *ControllerPayment) PaymentStatisticsGet(ctx context.Context, req *payment.PaymentStatisticsGetReq) (res *payment.PaymentStatisticsGetRes, err error) {
	// Call service to fetch payment statistics
	result, err := service.Payment().PaymentStatisticsGet(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeInternalError, err.Error())
	}

	// Direct write to response - return PayOS-like format
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(result)
		return nil, nil
	}

	return nil, gerror.New("failed to get request context")
}
