// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ParkingOrders is the golang structure of table parking_orders for DAO operations like Where/Data.
type ParkingOrders struct {
	g.Meta        `orm:"table:parking_orders, do:true"`
	Id            interface{} //
	UserId        interface{} //
	SlotId        interface{} //
	LotId         interface{} //
	StartTime     *gtime.Time //
	EndTime       *gtime.Time //
	Status        interface{} //
	Price         interface{} //
	PaymentStatus interface{} //
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
	VehicleId     interface{} //
}
