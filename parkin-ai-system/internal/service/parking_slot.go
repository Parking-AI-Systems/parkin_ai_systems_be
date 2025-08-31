package service

import (
	"parkin-ai-system/api/parking_slot"
)

type IParkingSlot interface {
	parking_slot.IParkingSlot
}

var (
	localParkingSlot IParkingSlot
)

func ParkingSlot() IParkingSlot {
	return localParkingSlot
}

func RegisterParkingSlot(i IParkingSlot) {
	localParkingSlot = i
}
