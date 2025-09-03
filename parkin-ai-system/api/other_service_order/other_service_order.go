package other_service_order

import (
	"github.com/gogf/gf/v2/frame/g"
)

type OtherServiceOrderAddReq struct {
	g.Meta         `path:"/other-service-orders" method:"post" tags:"OtherServiceOrder" summary:"Tạo đơn dịch vụ khác" security:"BearerAuth"`
	VehicleId      int64   `json:"vehicle_id" v:"required"`
	ServiceId      int64   `json:"service_id" v:"required"`
	LotId          int64   `json:"lot_id" v:"required"`
	ScheduledTime  string  `json:"scheduled_time" v:"required"`
	Status         string  `json:"status"`
	Price          float64 `json:"price"`
	PaymentStatus  string  `json:"payment_status"`
}

type OtherServiceOrderAddRes struct {
	OrderId int64 `json:"order_id"`
}

type OtherServiceOrderUpdateReq struct {
	g.Meta         `path:"/other-service-orders/{id}" method:"put" tags:"OtherServiceOrder" summary:"Cập nhật đơn dịch vụ khác" security:"BearerAuth"`
	Id             int64   `json:"id" v:"required"`
	Status         string  `json:"status"`
	Price          float64 `json:"price"`
	PaymentStatus  string  `json:"payment_status"`
	ScheduledTime  string  `json:"scheduled_time"`
}

type OtherServiceOrderUpdateRes struct {
	Success bool `json:"success"`
}

type OtherServiceOrderDeleteReq struct {
	g.Meta   `path:"/other-service-orders/{id}" method:"delete" tags:"OtherServiceOrder" summary:"Xoá đơn dịch vụ khác" security:"BearerAuth"`
	Id       int64  `json:"id" v:"required"`
}

type OtherServiceOrderDeleteRes struct {
	Success bool `json:"success"`
}

type OtherServiceOrderListReq struct {
	g.Meta `path:"/other-service-orders" method:"get" tags:"OtherServiceOrder" summary:"Danh sách đơn dịch vụ khác" security:"BearerAuth"`
	UserId int64 `json:"user_id"` // optional, for admin
}

type OtherServiceOrderListRes struct {
	List []OtherServiceOrderItem `json:"list"`
}

type OtherServiceOrderItem struct {
	Id            int64   `json:"id"`
	UserId        int64   `json:"user_id"`
	VehicleId     int64   `json:"vehicle_id"`
	ServiceId     int64   `json:"service_id"`
	LotId         int64   `json:"lot_id"`
	ScheduledTime string  `json:"scheduled_time"`
	Status        string  `json:"status"`
	Price         float64 `json:"price"`
	PaymentStatus string  `json:"payment_status"`
	CreatedAt     string  `json:"created_at"`
}
