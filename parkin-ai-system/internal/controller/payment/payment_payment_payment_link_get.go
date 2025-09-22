package payment

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"

	"parkin-ai-system/api/payment/payment"
	"parkin-ai-system/internal/service"
)

func (c *ControllerPayment) PaymentLinkGet(ctx context.Context, req *payment.PaymentLinkGetReq) (res *payment.PaymentLinkGetRes, err error) {
	// Validate request
	if req.Id == "" {
		if r := g.RequestFromCtx(ctx); r != nil {
			r.Response.WriteJson(g.Map{"error": "Payment link ID is required"})
			return nil, nil
		}
		return nil, gerror.New("Payment link ID is required")
	}

	// Call service
	result, err := service.Payment().PaymentLinkGet(ctx, req.Id)
	if err != nil {
		if r := g.RequestFromCtx(ctx); r != nil {
			r.Response.WriteJson(g.Map{"error": err.Error()})
			return nil, nil
		}
		return nil, err
	}

	// Convert result to PaymentLinkItem
	resultMap := gconv.Map(result)
	paymentLink := payment.PaymentLinkItem{
		Id:              gconv.String(resultMap["id"]),
		OrderCode:       gconv.Int64(resultMap["orderCode"]),
		Amount:          gconv.Int(resultMap["amount"]),
		AmountPaid:      gconv.Int(resultMap["amountPaid"]),
		AmountRemaining: gconv.Int(resultMap["amountRemaining"]),
		Status:          gconv.String(resultMap["status"]),
		CreatedAt:       gconv.String(resultMap["createdAt"]),
	}

	// Handle optional fields
	if resultMap["cancellationReason"] != nil {
		reason := gconv.String(resultMap["cancellationReason"])
		paymentLink.CancellationReason = &reason
	}
	if resultMap["cancelledAt"] != nil {
		cancelledAt := gconv.String(resultMap["cancelledAt"])
		paymentLink.CancelledAt = &cancelledAt
	}

	// Handle transactions
	if transactions, ok := resultMap["transactions"].([]interface{}); ok {
		for _, txInterface := range transactions {
			tx := gconv.Map(txInterface)
			txItem := payment.TransactionItem{
				Reference:           gconv.String(tx["reference"]),
				Amount:              gconv.Int(tx["amount"]),
				AccountNumber:       gconv.String(tx["accountNumber"]),
				Description:         gconv.String(tx["description"]),
				TransactionDateTime: gconv.String(tx["transactionDateTime"]),
			}

			// Handle optional transaction fields
			if tx["virtualAccountName"] != nil {
				name := gconv.String(tx["virtualAccountName"])
				txItem.VirtualAccountName = &name
			}
			if tx["virtualAccountNumber"] != nil {
				number := gconv.String(tx["virtualAccountNumber"])
				txItem.VirtualAccountNumber = &number
			}
			if tx["counterAccountBankId"] != nil {
				bankId := gconv.String(tx["counterAccountBankId"])
				txItem.CounterAccountBankId = &bankId
			}
			if tx["counterAccountBankName"] != nil {
				bankName := gconv.String(tx["counterAccountBankName"])
				txItem.CounterAccountBankName = &bankName
			}
			if tx["counterAccountName"] != nil {
				accountName := gconv.String(tx["counterAccountName"])
				txItem.CounterAccountName = &accountName
			}
			if tx["counterAccountNumber"] != nil {
				accountNumber := gconv.String(tx["counterAccountNumber"])
				txItem.CounterAccountNumber = &accountNumber
			}

			paymentLink.Transactions = append(paymentLink.Transactions, txItem)
		}
	}

	resJson := &payment.PaymentLinkGetRes{
		PaymentLink: paymentLink,
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(resJson)
	}
	return nil, nil
}
