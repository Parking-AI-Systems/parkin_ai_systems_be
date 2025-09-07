// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ParkingLots is the golang structure of table parking_lots for DAO operations like Where/Data.
type ParkingLots struct {
	g.Meta         `orm:"table:parking_lots, do:true"`
	Id             interface{} //
	Name           interface{} //
	Address        interface{} //
	Latitude       interface{} //
	Longitude      interface{} //
	OwnerId        interface{} //
	IsVerified     interface{} //
	IsActive       interface{} //
	TotalSlots     interface{} //
	AvailableSlots interface{} //
	PricePerHour   interface{} //
	Description    interface{} //
	OpenTime       *gtime.Time //
	CloseTime      *gtime.Time //
	ImageUrl       interface{} //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
	DeletedAt	  *gtime.Time //
}
