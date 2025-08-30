package other_service

import (
	"github.com/gogf/gf/v2/frame/g"
)

type OtherServiceAddReq struct {
	g.Meta          `path:"/other-services" method:"post" tags:"OtherService" summary:"Thêm dịch vụ khác" security:"BearerAuth"`
	LotId           int64   `json:"lot_id" v:"required"`
	Name            string  `json:"name" v:"required|length:1,64"`
	Description     string  `json:"description"`
	Price           float64 `json:"price" v:"required|min:0"`
	DurationMinutes int     `json:"duration_minutes" v:"required|min:1"`
	IsActive        bool    `json:"is_active"`
}

type OtherServiceAddRes struct {
	ServiceId int64 `json:"service_id"`
}

type OtherServiceUpdateReq struct {
	g.Meta          `path:"/other-services/{id}" method:"put" tags:"OtherService" summary:"Cập nhật dịch vụ khác" security:"BearerAuth"`
	Id              int64   `json:"id" v:"required"`
	Name            string  `json:"name" v:"length:1,64"`
	Description     string  `json:"description"`
	Price           float64 `json:"price" v:"min:0"`
	DurationMinutes int     `json:"duration_minutes" v:"min:1"`
	IsActive        bool    `json:"is_active"`
}

type OtherServiceUpdateRes struct {
	Success bool `json:"success"`
}

type OtherServiceDeleteReq struct {
	g.Meta   `path:"/other-services/{id}" method:"delete" tags:"OtherService" summary:"Xoá dịch vụ khác" security:"BearerAuth"`
	Id       int64  `json:"id" v:"required"`
}

type OtherServiceDeleteRes struct {
	Success bool `json:"success"`
}

type OtherServiceDetailReq struct {
	g.Meta   `path:"/other-services/{id}" method:"get" tags:"OtherService" summary:"Chi tiết dịch vụ khác" security:"BearerAuth"`
	Id       int64  `json:"id" v:"required"`
}

type OtherServiceDetailRes struct {
	Service *OtherServiceInfo `json:"service"`
}

type OtherServiceListReq struct {
	g.Meta   `path:"/other-services" method:"get" tags:"OtherService" summary:"Danh sách dịch vụ khác" security:"BearerAuth"`
	LotId    int64  `json:"lot_id"`
}

type OtherServiceListRes struct {
	Services []OtherServiceInfo `json:"services"`
}

type OtherServiceInfo struct {
	Id              int64   `json:"id"`
	LotId           int64   `json:"lot_id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	DurationMinutes int     `json:"duration_minutes"`
	IsActive        bool    `json:"is_active"`
	CreatedAt       string  `json:"created_at"`
}
