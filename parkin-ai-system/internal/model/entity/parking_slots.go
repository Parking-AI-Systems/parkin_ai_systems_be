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
}
