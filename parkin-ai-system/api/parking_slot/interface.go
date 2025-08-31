package parking_slot

type IParkingSlot interface {
	ParkingSlotAdd(req *ParkingSlotAddReq) (*ParkingSlotAddRes, error)
	ParkingSlotList(req *ParkingSlotListReq) (*ParkingSlotListRes, error)
	ParkingSlotUpdate(req *ParkingSlotUpdateReq) (*ParkingSlotUpdateRes, error)
	ParkingSlotDelete(req *ParkingSlotDeleteReq) (*ParkingSlotDeleteRes, error)
}
