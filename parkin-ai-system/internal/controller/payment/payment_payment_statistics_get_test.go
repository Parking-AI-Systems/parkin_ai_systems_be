package payment

import (
	"testing"

	"parkin-ai-system/api/payment/payment"
)

func TestPaymentStatisticsResponse(t *testing.T) {
	res := &payment.PaymentStatisticsGetRes{
		Code: "00",
		Desc: "success",
		Data: map[string]interface{}{
			"test": "value",
		},
	}

	if res.Code == "" {
		t.Error("Code is empty")
	}
	if res.Desc == "" {
		t.Error("Desc is empty")
	}
	if len(res.Data) == 0 {
		t.Error("Data is empty")
	}

	t.Logf("Response: %+v", res)
}
