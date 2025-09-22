package payment

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"

	"parkin-ai-system/api/payment/payment"
	"parkin-ai-system/internal/service"
)

func (c *ControllerPayment) CheckoutAdd(ctx context.Context, req *payment.CheckoutAddReq) (res *payment.CheckoutAddRes, err error) {
	// Validate request
	if req.Amount < 1000 {
		return nil, gerror.New("Amount must be at least 1000 VND")
	}

	// Prepare request data for service
	reqData := map[string]interface{}{
		"orderCode":    req.OrderCode,
		"amount":       req.Amount,
		"description":  req.Description,
		"cancelUrl":    req.CancelUrl,
		"returnUrl":    req.ReturnUrl,
		"items":        req.Items,
		"buyerName":    req.BuyerName,
		"buyerEmail":   req.BuyerEmail,
		"buyerPhone":   req.BuyerPhone,
		"buyerAddress": req.BuyerAddress,
		"expiredAt":    req.ExpiredAt,
	}

	// Call service
	result, err := service.Payment().CheckoutAdd(ctx, reqData)
	if err != nil {
		return nil, err
	}

	// Convert result
	resultMap := gconv.Map(result)
	return &payment.CheckoutAddRes{
		PaymentLinkId: gconv.String(resultMap["paymentLinkId"]),
		CheckoutUrl:   gconv.String(resultMap["checkoutUrl"]),
		QRCode:        gconv.String(resultMap["qrCode"]),
	}, nil
}
