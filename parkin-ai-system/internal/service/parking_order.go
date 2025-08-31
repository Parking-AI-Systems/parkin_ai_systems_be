package service

import (
	"parkin-ai-system/api/parking_order"
)

type IParkingOrder interface {
	parking_order.IParkingOrder
}

var (
	localParkingOrder IParkingOrder
)

func ParkingOrder() IParkingOrder {
	return localParkingOrder
}

func RegisterParkingOrder(i IParkingOrder) {
	localParkingOrder = i
}
