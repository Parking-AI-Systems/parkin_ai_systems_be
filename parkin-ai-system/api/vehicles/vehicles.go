package vehicles

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Tạo mới vehicle
 type CreateVehicleReq struct {
	g.Meta       `path:"/vehicles" method:"post" tags:"Vehicles" summary:"Create vehicle"`
	UserId       int64  `json:"user_id" description:"User ID"`
	LicensePlate string `json:"license_plate"`
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	Color        string `json:"color"`
	Type         string `json:"type"`
}

type CreateVehicleRes struct {
	Id int64 `json:"id"`
}

// Lấy danh sách vehicles
 type ListVehicleReq struct {
	g.Meta `path:"/vehicles" method:"get" tags:"Vehicles" summary:"List vehicles"`
	UserId int64 `json:"user_id" description:"User ID"`
}

type ListVehicleRes struct {
	Vehicles []VehicleItem `json:"vehicles"`
}

type VehicleItem struct {
	Id           int64  `json:"id"`
	UserId       int64  `json:"user_id"`
	LicensePlate string `json:"license_plate"`
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	Color        string `json:"color"`
	Type         string `json:"type"`
}

// Lấy chi tiết vehicle
 type GetVehicleReq struct {
	g.Meta `path:"/vehicles/{id}" method:"get" tags:"Vehicles" summary:"Get vehicle detail"`
	Id     int64 `json:"id" in:"path"`
}

type GetVehicleRes struct {
	Vehicle VehicleItem `json:"vehicle"`
}

// Cập nhật vehicle
 type UpdateVehicleReq struct {
	g.Meta       `path:"/vehicles/{id}" method:"put" tags:"Vehicles" summary:"Update vehicle"`
	Id           int64  `json:"id" in:"path"`
	LicensePlate string `json:"license_plate"`
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	Color        string `json:"color"`
	Type         string `json:"type"`
}

type UpdateVehicleRes struct {
	Success bool `json:"success"`
}

// Xóa vehicle
 type DeleteVehicleReq struct {
	g.Meta `path:"/vehicles/{id}" method:"delete" tags:"Vehicles" summary:"Delete vehicle"`
	Id     int64 `json:"id" in:"path"`
}

type DeleteVehicleRes struct {
	Success bool `json:"success"`
}
