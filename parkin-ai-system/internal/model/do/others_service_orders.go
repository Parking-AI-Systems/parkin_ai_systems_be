// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OthersServiceOrders is the golang structure of table others_service_orders for DAO operations like Where/Data.
type OthersServiceOrders struct {
	g.Meta        `orm:"table:others_service_orders, do:true"`
	Id            interface{} //
	UserId        interface{} //
	VehicleId     interface{} //
	ServiceId     interface{} //
	LotId         interface{} //
	ScheduledTime *gtime.Time //
	Status        interface{} //
	Price         interface{} //
	PaymentStatus interface{} //
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
	DeletedAt     *gtime.Time //
}
