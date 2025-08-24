// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ParkingOrders is the golang structure for table parking_orders.
type ParkingOrders struct {
	Id            int64       `json:"id"            orm:"id"             description:""`
	UserId        int64       `json:"userId"        orm:"user_id"        description:""`
	SlotId        int64       `json:"slotId"        orm:"slot_id"        description:""`
	LotId         int64       `json:"lotId"         orm:"lot_id"         description:""`
	StartTime     *gtime.Time `json:"startTime"     orm:"start_time"     description:""`
	EndTime       *gtime.Time `json:"endTime"       orm:"end_time"       description:""`
	Status        string      `json:"status"        orm:"status"         description:""`
	Price         float64     `json:"price"         orm:"price"          description:""`
	PaymentStatus string      `json:"paymentStatus" orm:"payment_status" description:""`
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:""`
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     description:""`
}
