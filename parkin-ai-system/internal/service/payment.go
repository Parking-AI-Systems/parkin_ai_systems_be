// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	IPayment interface {
		CreatePaymentLink(ctx context.Context, orderType string, orderID int64) (interface{}, error)
		HandlePaymentWebhook(ctx context.Context, webhookData interface{}) error
		CheckoutAdd(ctx context.Context, req interface{}) (interface{}, error)
		PaymentLinkGet(ctx context.Context, paymentLinkId string) (interface{}, error)
		RefundAdd(ctx context.Context, paymentLinkId string, amount int, reason *string) (interface{}, error)
	}
)

var (
	localPayment IPayment
)

func Payment() IPayment {
	if localPayment == nil {
		panic("implement not found for interface IPayment, forgot register?")
	}
	return localPayment
}

func RegisterPayment(i IPayment) {
	localPayment = i
}

