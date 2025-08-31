package parking_slot

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ParkingSlotAddReq defines the request struct for adding a parking slot
// swagger:parameters ParkingSlotAddReq
// in: body
// required: true
//
type ParkingSlotAddReq struct {
	g.Meta      `path:"/parking-slots" method:"post" tags:"ParkingSlot" summary:"Add Parking Slot"`
	LotId       int64  `json:"lot_id" description:"Parking lot ID"`
	Code        string `json:"code" description:"Slot code"`
	IsAvailable bool   `json:"is_available" description:"Is available"`
	SlotType    string `json:"slot_type" description:"Slot type"`
	Floor       string `json:"floor" description:"Floor"`
}

type ParkingSlotAddRes struct {
	Id int64 `json:"id" description:"Slot ID"`
}

// ParkingSlotListReq defines the request struct for listing parking slots
// swagger:parameters ParkingSlotListReq
// in: query
// required: false
//
type ParkingSlotListReq struct {
	g.Meta `path:"/parking-slots" method:"get" tags:"ParkingSlot" summary:"List Parking Slots"`
	LotId  int64 `json:"lot_id" description:"Parking lot ID"`
}

type ParkingSlotListRes struct {
	List []ParkingSlotItem `json:"list"`
}

type ParkingSlotItem struct {
	Id          int64  `json:"id"`
	LotId       int64  `json:"lot_id"`
	Code        string `json:"code"`
	IsAvailable bool   `json:"is_available"`
	SlotType    string `json:"slot_type"`
	Floor       string `json:"floor"`
	CreatedAt   string `json:"created_at"`
}

// ParkingSlotUpdateReq defines the request struct for updating a parking slot
// swagger:parameters ParkingSlotUpdateReq
// in: body
// required: true
//
type ParkingSlotUpdateReq struct {
	g.Meta      `path:"/parking-slots/{id}" method:"put" tags:"ParkingSlot" summary:"Update Parking Slot"`
	Id          int64  `json:"id" description:"Slot ID"`
	Code        string `json:"code" description:"Slot code"`
	IsAvailable bool   `json:"is_available" description:"Is available"`
	SlotType    string `json:"slot_type" description:"Slot type"`
	Floor       string `json:"floor" description:"Floor"`
}

type ParkingSlotUpdateRes struct {
	Success bool `json:"success"`
}

// ParkingSlotDeleteReq defines the request struct for deleting a parking slot
// swagger:parameters ParkingSlotDeleteReq
// in: path
// required: true
//
type ParkingSlotDeleteReq struct {
	g.Meta `path:"/parking-slots/{id}" method:"delete" tags:"ParkingSlot" summary:"Delete Parking Slot"`
	Id     int64 `json:"id" description:"Slot ID"`
}

type ParkingSlotDeleteRes struct {
	Success bool `json:"success"`
}
