package payment

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"

	"parkin-ai-system/api/payment/payment"
	"parkin-ai-system/internal/service"
)

func (c *ControllerPayment) CreatePaymentLink(ctx context.Context, req *payment.CreatePaymentLinkReq) (res *payment.CreatePaymentLinkRes, err error) {
	// Validate request
	if req.OrderType == "" {
		return nil, gerror.New("OrderType is required")
	}
	if req.OrderID <= 0 {
		return nil, gerror.New("OrderID must be greater than 0")
	}
	if req.OrderType != "parking" && req.OrderType != "service" {
		return nil, gerror.New("OrderType must be 'parking' or 'service'")
	}

	// Call service
	result, err := service.Payment().CreatePaymentLink(ctx, req.OrderType, req.OrderID)
	if err != nil {
		if r := g.RequestFromCtx(ctx); r != nil {
			r.Response.WriteJson(g.Map{"error": err.Error()})
			return nil, nil
		}
		return nil, err
	}

	// Convert result
	resultMap := gconv.Map(result)
	resJson := &payment.CreatePaymentLinkRes{
		PaymentLinkId: gconv.String(resultMap["paymentLinkId"]),
		CheckoutUrl:   gconv.String(resultMap["checkoutUrl"]),
		QRCode:        gconv.String(resultMap["qrCode"]),
		Amount:        gconv.Int(resultMap["amount"]),
		OrderCode:     gconv.Int64(resultMap["orderCode"]),
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(resJson)
	}
	return nil, nil
}
