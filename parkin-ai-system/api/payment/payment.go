// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package payment

import (
	"context"

	"parkin-ai-system/api/payment/payment"
)

type IPaymentPayment interface {
	CheckoutAdd(ctx context.Context, req *payment.CheckoutAddReq) (res *payment.CheckoutAddRes, err error)
	PaymentLinkGet(ctx context.Context, req *payment.PaymentLinkGetReq) (res *payment.PaymentLinkGetRes, err error)
	RefundAdd(ctx context.Context, req *payment.RefundAddReq) (res *payment.RefundAddRes, err error)
	Webhook(ctx context.Context, req *payment.WebhookReq) (res *payment.WebhookRes, err error)
	CreatePaymentLink(ctx context.Context, req *payment.CreatePaymentLinkReq) (res *payment.CreatePaymentLinkRes, err error)
	PaymentStatisticsGet(ctx context.Context, req *payment.PaymentStatisticsGetReq) (res *payment.PaymentStatisticsGetRes, err error)
}
