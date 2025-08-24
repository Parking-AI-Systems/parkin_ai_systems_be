// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ParkingSlots is the golang structure of table parking_slots for DAO operations like Where/Data.
type ParkingSlots struct {
	g.Meta      `orm:"table:parking_slots, do:true"`
	Id          interface{} //
	LotId       interface{} //
	Code        interface{} //
	IsAvailable interface{} //
	SlotType    interface{} //
	Floor       interface{} //
	CreatedAt   *gtime.Time //
}
