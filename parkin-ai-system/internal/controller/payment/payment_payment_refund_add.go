package payment

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"

	"parkin-ai-system/api/payment/payment"
	"parkin-ai-system/internal/service"
)

func (c *ControllerPayment) RefundAdd(ctx context.Context, req *payment.RefundAddReq) (res *payment.RefundAddRes, err error) {
	// Validate request
	if req.Id == "" {
		return nil, gerror.New("Payment link ID is required")
	}
	if req.Amount <= 0 {
		return nil, gerror.New("Refund amount must be greater than 0")
	}

	// Call service
	result, err := service.Payment().RefundAdd(ctx, req.Id, req.Amount, req.Reason)
	if err != nil {
		return nil, err
	}

	// Convert result
	resultMap := gconv.Map(result)
	return &payment.RefundAddRes{
		RefundId: gconv.String(resultMap["refundId"]),
		Status:   gconv.String(resultMap["status"]),
	}, nil
}
