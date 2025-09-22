package payment

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/payOSHQ/payos-lib-golang"

	"parkin-ai-system/api/payment/payment"
	"parkin-ai-system/internal/service"
)

func (c *ControllerPayment) Webhook(ctx context.Context, req *payment.WebhookReq) (res *payment.WebhookRes, err error) {
	// Validate request
	if req.Signature == "" {
		if r := g.RequestFromCtx(ctx); r != nil {
			r.Response.WriteJson(g.Map{"error": "Webhook signature is required"})
			return nil, nil
		}
		return nil, gerror.New("Webhook signature is required")
	}

	// Convert to PayOS webhook structure
	var payosData *payos.WebhookDataType
	if req.Data != nil {
		payosData = &payos.WebhookDataType{
			OrderCode:              req.Data.OrderCode,
			Amount:                 req.Data.Amount,
			Description:            req.Data.Description,
			AccountNumber:          req.Data.AccountNumber,
			Reference:              req.Data.Reference,
			TransactionDateTime:    req.Data.TransactionDateTime,
			Currency:               req.Data.Currency,
			PaymentLinkId:          req.Data.PaymentLinkId,
			Code:                   req.Data.Code,
			Desc:                   req.Data.Desc,
			CounterAccountBankId:   req.Data.CounterAccountBankId,
			CounterAccountBankName: req.Data.CounterAccountBankName,
			CounterAccountName:     req.Data.CounterAccountName,
			CounterAccountNumber:   req.Data.CounterAccountNumber,
			VirtualAccountName:     req.Data.VirtualAccountName,
			VirtualAccountNumber:   req.Data.VirtualAccountNumber,
		}
	}

	webhookData := payos.WebhookType{
		Code:      req.Code,
		Desc:      req.Desc,
		Success:   req.Success,
		Data:      payosData,
		Signature: req.Signature,
	}

	// Call service to handle webhook
	err = service.Payment().HandlePaymentWebhook(ctx, webhookData)
	if err != nil {
		if r := g.RequestFromCtx(ctx); r != nil {
			r.Response.WriteJson(g.Map{"error": err.Error()})
			return nil, nil
		}
		return nil, err
	}

	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(&payment.WebhookRes{
			Message: "Webhook processed successfully",
		})
	}
	return nil, nil
}
