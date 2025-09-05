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
	VehicleId     int64       `json:"vehicleId"     orm:"vehicle_id"     description:""`
}

type ParkingOrderAddReq struct {
	VehicleId int64  `json:"vehicleId"`
	LotId     int64  `json:"lotId"`
	SlotId    int64  `json:"slotId"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}
type ParkingOrderAddRes struct {
	Id int64 `json:"id"`
}
type ParkingOrderListReq struct {
	UserId int64 `json:"userId"`
	LotId  int64 `json:"lotId"`
	Page int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Status   string `json:"status"`
}

type ParkingOrderItem struct {
	Id            int64       `json:"id"`            // Order ID
	UserId        int64       `json:"user_id"`       // User ID
	LotId         int64       `json:"lot_id"`        // Parking lot ID
	SlotId        int64       `json:"slot_id"`       // Parking slot ID
	VehicleId     int64       `json:"vehicle_id"`    // Vehicle ID
	LotName       string      `json:"lot_name"`      // Parking lot name (from join)
	SlotCode      string      `json:"slot_code"`     // Parking slot code (from join)
	VehiclePlate  string      `json:"vehicle_plate"` // Vehicle license plate (from join)
	StartTime     string      `json:"start_time"`    // Start time (formatted)
	EndTime       string      `json:"end_time"`      // End time (formatted)
	Status        string      `json:"status"`        // Order status
	Price         float64     `json:"price"`         // Order price
	PaymentStatus string      `json:"payment_status"` // Payment status
	CreatedAt     string      `json:"created_at"`    // Creation timestamp (formatted)
	UpdatedAt     string      `json:"updated_at"`    // Update timestamp (formatted, optional)
}

type ParkingOrderListRes struct {
	List  []ParkingOrderItem `json:"list"`  // List of parking orders
	Total int                `json:"total"` // Total number of matching orders
}

type ParkingOrderGetReq struct {
	Id int64 `json:"id"` // Parking order ID
}

type ParkingOrderUpdateReq struct {
	Id        int64  `json:"id" v:"required|min:1#Order ID is required|Order ID must be positive"` // Parking order ID
	StartTime string `json:"start_time" v:"date#Invalid start time format"`                        // New start time (optional)
	EndTime   string `json:"end_time" v:"date#Invalid end time format"`                           // New end time (optional)
	Status    string `json:"status" v:"in:pending,confirmed,canceled,completed#Invalid status value"` // New status (optional)
}
type ParkingOrderCancelReq struct {
	Id int64 `json:"id"` // Parking order ID
}

type ParkingOrderDeleteReq struct {
	Id int64 `json:"id"` // Parking order ID
}

type ParkingOrderDeleteRes struct {
	Message string `json:"message"` // Confirmation message
}

type ParkingOrderPaymentReq struct {
	Id int64 `json:"id"` // Parking order ID
}
