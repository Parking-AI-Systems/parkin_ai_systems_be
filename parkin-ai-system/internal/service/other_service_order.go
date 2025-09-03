package service

import (
	"context"
	"parkin-ai-system/api/other_service_order"
)

type IOtherServiceOrder interface {
	other_service_order.IOtherServiceOrder
}

var localOtherServiceOrder IOtherServiceOrder

func OtherServiceOrder() IOtherServiceOrder {
	if localOtherServiceOrder == nil {
		panic("implement not found for interface IOtherServiceOrder, forgot register?")
	}
	return localOtherServiceOrder
}

func RegisterOtherServiceOrder(i IOtherServiceOrder) {
	localOtherServiceOrder = i
}
