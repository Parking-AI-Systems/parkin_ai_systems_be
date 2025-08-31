package parking_order

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ParkingOrderAddReq defines the request struct for adding a parking order
// swagger:parameters ParkingOrderAddReq
// in: body
// required: true
//
type ParkingOrderAddReq struct {
	g.Meta      `path:"/parking-orders" method:"post" tags:"ParkingOrder" summary:"Add Parking Order"`
	UserId      int64  `json:"user_id" description:"User ID"`
	SlotId      int64  `json:"slot_id" description:"Slot ID"`
	LotId       int64  `json:"lot_id" description:"Lot ID"`
	StartTime   string `json:"start_time" description:"Start time"`
	EndTime     string `json:"end_time" description:"End time"`
	Status      string `json:"status" description:"Status"`
	Price       float64 `json:"price" description:"Price"`
	PaymentStatus string `json:"payment_status" description:"Payment status"`
}

type ParkingOrderAddRes struct {
	Id int64 `json:"id" description:"Order ID"`
}

// ParkingOrderListReq defines the request struct for listing parking orders
// swagger:parameters ParkingOrderListReq
// in: query
// required: false
//
type ParkingOrderListReq struct {
	g.Meta `path:"/parking-orders" method:"get" tags:"ParkingOrder" summary:"List Parking Orders"`
	UserId int64 `json:"user_id" description:"User ID"`
	LotId  int64 `json:"lot_id" description:"Lot ID"`
}

type ParkingOrderListRes struct {
	List []ParkingOrderItem `json:"list"`
}

type ParkingOrderItem struct {
	Id            int64   `json:"id"`
	UserId        int64   `json:"user_id"`
	SlotId        int64   `json:"slot_id"`
	LotId         int64   `json:"lot_id"`
	StartTime     string  `json:"start_time"`
	EndTime       string  `json:"end_time"`
	Status        string  `json:"status"`
	Price         float64 `json:"price"`
	PaymentStatus string  `json:"payment_status"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

// ParkingOrderUpdateReq defines the request struct for updating a parking order
// swagger:parameters ParkingOrderUpdateReq
// in: body
// required: true
//
type ParkingOrderUpdateReq struct {
	g.Meta      `path:"/parking-orders/{id}" method:"put" tags:"ParkingOrder" summary:"Update Parking Order"`
	Id          int64   `json:"id" description:"Order ID"`
	Status      string  `json:"status" description:"Status"`
	EndTime     string  `json:"end_time" description:"End time"`
	Price       float64 `json:"price" description:"Price"`
	PaymentStatus string `json:"payment_status" description:"Payment status"`
}

type ParkingOrderUpdateRes struct {
	Success bool `json:"success"`
}

// ParkingOrderDeleteReq defines the request struct for deleting a parking order
// swagger:parameters ParkingOrderDeleteReq
// in: path
// required: true
//
type ParkingOrderDeleteReq struct {
	g.Meta `path:"/parking-orders/{id}" method:"delete" tags:"ParkingOrder" summary:"Delete Parking Order"`
	Id     int64 `json:"id" description:"Order ID"`
}

type ParkingOrderDeleteRes struct {
	Success bool `json:"success"`
}
