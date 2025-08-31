package parking_slot

import (
	"context"
	"parkin-ai-system/api/parking_slot"
	"parkin-ai-system/internal/service"
)

type ControllerParkingSlot struct{}

func NewParkingSlot() *ControllerParkingSlot {
	return &ControllerParkingSlot{}
}

func (c *ControllerParkingSlot) ParkingSlotAdd(ctx context.Context, req *parking_slot.ParkingSlotAddReq) (res *parking_slot.ParkingSlotAddRes, err error) {
	return service.ParkingSlot().ParkingSlotAdd(req)
}

func (c *ControllerParkingSlot) ParkingSlotList(ctx context.Context, req *parking_slot.ParkingSlotListReq) (res *parking_slot.ParkingSlotListRes, err error) {
	return service.ParkingSlot().ParkingSlotList(req)
}

func (c *ControllerParkingSlot) ParkingSlotUpdate(ctx context.Context, req *parking_slot.ParkingSlotUpdateReq) (res *parking_slot.ParkingSlotUpdateRes, err error) {
	return service.ParkingSlot().ParkingSlotUpdate(req)
}

func (c *ControllerParkingSlot) ParkingSlotDelete(ctx context.Context, req *parking_slot.ParkingSlotDeleteReq) (res *parking_slot.ParkingSlotDeleteRes, err error) {
	return service.ParkingSlot().ParkingSlotDelete(req)
}
