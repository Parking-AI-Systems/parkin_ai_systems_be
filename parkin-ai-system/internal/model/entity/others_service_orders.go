// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OthersServiceOrders is the golang structure for table others_service_orders.
type OthersServiceOrders struct {
	Id            int64       `json:"id"            orm:"id"             description:""`
	UserId        int64       `json:"userId"        orm:"user_id"        description:""`
	VehicleId     int64       `json:"vehicleId"     orm:"vehicle_id"     description:""`
	ServiceId     int64       `json:"serviceId"     orm:"service_id"     description:""`
	LotId         int64       `json:"lotId"         orm:"lot_id"         description:""`
	ScheduledTime *gtime.Time `json:"scheduledTime" orm:"scheduled_time" description:""`
	Status        string      `json:"status"        orm:"status"         description:""`
	Price         float64     `json:"price"         orm:"price"          description:""`
	PaymentStatus string      `json:"paymentStatus" orm:"payment_status" description:""`
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:""`
}
