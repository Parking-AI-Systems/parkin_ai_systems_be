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

type OthersServiceOrderAddReq struct {
	VehicleId    int64  `json:"vehicleId"`
	LotId        int64  `json:"lotId"`
	ServiceId    int64  `json:"serviceId"`
	ScheduledTime string `json:"scheduledTime"`
}

type OthersServiceOrderAddRes struct {
	Id int64 `json:"id"`
}

type OthersServiceOrderListReq struct {
	UserId   int64  `json:"userId"`
	LotId    int64  `json:"lotId"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Status   string `json:"status"`
}

type OthersServiceOrderItem struct {
	Id            int64   `json:"id"`
	UserId        int64   `json:"user_id"`
	LotId         int64   `json:"lot_id"`
	ServiceId     int64   `json:"service_id"`
	VehicleId     int64   `json:"vehicle_id"`
	LotName       string  `json:"lot_name"`
	ServiceName   string  `json:"service_name"`
	VehiclePlate  string  `json:"vehicle_plate"`
	ScheduledTime string  `json:"scheduled_time"`
	Status        string  `json:"status"`
	Price         float64 `json:"price"`
	PaymentStatus string  `json:"payment_status"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type OthersServiceOrderListRes struct {
	List  []OthersServiceOrderItem `json:"list"`
	Total int                      `json:"total"`
}

type OthersServiceOrderGetReq struct {
	Id int64 `json:"id"`
}

type OthersServiceOrderUpdateReq struct {
	Id            int64  `json:"id" v:"required|min:1#Order ID is required|Order ID must be positive"`
	ScheduledTime string `json:"scheduled_time" v:"date#Invalid scheduled time format"`
	Status        string `json:"status" v:"in:pending,confirmed,canceled,completed#Invalid status value"`
}

type OthersServiceOrderCancelReq struct {
	Id int64 `json:"id"`
}

type OthersServiceOrderDeleteReq struct {
	Id int64 `json:"id"`
}

type OthersServiceOrderDeleteRes struct {
	Message string `json:"message"`
}

type OthersServiceOrderPaymentReq struct {
	Id int64 `json:"id"`
}