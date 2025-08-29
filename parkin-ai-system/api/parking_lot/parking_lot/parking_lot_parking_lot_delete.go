package parking_lot

import "github.com/gogf/gf/v2/frame/g"

type ParkingLotDeleteReq struct {
	g.Meta `path:"/parking-lots/{id}" method:"delete" tags:"ParkingLot" summary:"Xoá bãi" security:"BearerAuth"`
	Id string `json:"id" v:"required"`
}

type ParkingLotDeleteRes struct {
	Success bool `json:"success"`
}
