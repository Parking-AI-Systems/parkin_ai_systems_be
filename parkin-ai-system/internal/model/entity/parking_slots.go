// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ParkingSlots is the golang structure for table parking_slots.
type ParkingSlots struct {
	Id          int64       `json:"id"          orm:"id"           description:""`
	LotId       int64       `json:"lotId"       orm:"lot_id"       description:""`
	Code        string      `json:"code"        orm:"code"         description:""`
	IsAvailable bool        `json:"isAvailable" orm:"is_available" description:""`
	SlotType    string      `json:"slotType"    orm:"slot_type"    description:""`
	Floor       string      `json:"floor"       orm:"floor"        description:""`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:""`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   description:""`
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"   description:""`
}

type ParkingSlotAddReq struct {
	LotId       int64  `json:"lotId"`
	Code        string `json:"code"`
	IsAvailable bool   `json:"isAvailable"`
	SlotType    string `json:"slotType"`
	Floor       string `json:"floor"`
}

type ParkingSlotAddRes struct {
	Id int64 `json:"id"`
}

type ParkingSlotListReq struct {
	LotId       int64  `json:"lotId"`
	IsAvailable *bool  `json:"isAvailable"`
	SlotType    string `json:"slotType"`
	Page        int    `json:"page"`
	PageSize    int    `json:"pageSize"`
}

type ParkingSlotItem struct {
	Id          int64  `json:"id"`
	LotId       int64  `json:"lot_id"`
	LotName     string `json:"lot_name"`
	Code        string `json:"code"`
	IsAvailable bool   `json:"is_available"`
	SlotType    string `json:"slot_type"`
	Floor       string `json:"floor"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DeletedAt   string `json:"deleted_at"`
}

type ParkingSlotListRes struct {
	List  []ParkingSlotItem `json:"list"`
	Total int               `json:"total"`
}

type ParkingSlotGetReq struct {
	Id int64 `json:"id"`
}

type ParkingSlotGetRes struct {
	Slot ParkingSlotItem `json:"slot"`
}

type ParkingSlotUpdateReq struct {
	Id          int64  `json:"id"`
	LotId       int64  `json:"lotId"`
	Code        string `json:"code"`
	IsAvailable *bool  `json:"isAvailable"`
	SlotType    string `json:"slotType"`
	Floor       string `json:"floor"`
}

type ParkingSlotUpdateRes struct {
	Slot ParkingSlotItem `json:"slot"`
}

type ParkingSlotDeleteReq struct {
	Id int64 `json:"id"`
}

type ParkingSlotDeleteRes struct {
	Message string `json:"message"`
}