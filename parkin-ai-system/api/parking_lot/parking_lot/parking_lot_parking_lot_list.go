package parking_lot

import "github.com/gogf/gf/v2/frame/g"

type ParkingLotListReq struct {
	g.Meta `path:"/parking-lots" method:"get" tags:"ParkingLot" summary:"Danh sách bãi" security:"BearerAuth"`
}

type ParkingLotListRes struct {
	Lots []ParkingLotInfo `json:"lots"`
}
