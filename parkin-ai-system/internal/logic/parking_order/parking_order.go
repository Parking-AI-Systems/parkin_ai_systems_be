package parking_order

import (
	"context"
	"parkin-ai-system/api/parking_order"
	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/model/do"
	"parkin-ai-system/internal/model/entity"
	"github.com/gogf/gf/v2/util/gconv"
	"time"
	"parkin-ai-system/internal/service"
)

type sParkingOrder struct{}

func init() {
	service.RegisterParkingOrder(&sParkingOrder{})
}

func (s *sParkingOrder) ParkingOrderAdd(req *parking_order.ParkingOrderAddReq) (*parking_order.ParkingOrderAddRes, error) {
	data := do.ParkingOrders{}
	gconv.Struct(req, &data)
	lastId, err := dao.ParkingOrders.Ctx(context.Background()).Data(data).InsertAndGetId()
	if err != nil {
		return nil, err
	}
	return &parking_order.ParkingOrderAddRes{Id: gconv.Int64(lastId)}, nil
}

func (s *sParkingOrder) ParkingOrderList(req *parking_order.ParkingOrderListReq) (*parking_order.ParkingOrderListRes, error) {
	var orders []entity.ParkingOrders
	m := dao.ParkingOrders.Ctx(context.Background())
	if req.UserId != 0 {
		m = m.Where("user_id", req.UserId)
	}
	if req.LotId != 0 {
		m = m.Where("lot_id", req.LotId)
	}
	err := m.Order("id desc").Scan(&orders)
	if err != nil {
		return nil, err
	}
	var list []parking_order.ParkingOrderItem
	for _, order := range orders {
		item := parking_order.ParkingOrderItem{}
		gconv.Struct(order, &item)
		item.CreatedAt = order.CreatedAt.Format("2006-01-02 15:04:05")
		item.UpdatedAt = ""
		if !order.UpdatedAt.IsZero() {
			item.UpdatedAt = order.UpdatedAt.Format("2006-01-02 15:04:05")
		}
		list = append(list, item)
	}
	return &parking_order.ParkingOrderListRes{List: list}, nil
}

func (s *sParkingOrder) ParkingOrderUpdate(req *parking_order.ParkingOrderUpdateReq) (*parking_order.ParkingOrderUpdateRes, error) {
	data := do.ParkingOrders{}
	gconv.Struct(req, &data)
	data.UpdatedAt = time.Now()
	_, err := dao.ParkingOrders.Ctx(context.Background()).Where("id", req.Id).Data(data).Update()
	if err != nil {
		return nil, err
	}
	return &parking_order.ParkingOrderUpdateRes{Success: true}, nil
}

func (s *sParkingOrder) ParkingOrderDelete(req *parking_order.ParkingOrderDeleteReq) (*parking_order.ParkingOrderDeleteRes, error) {
	_, err := dao.ParkingOrders.Ctx(context.Background()).Where("id", req.Id).Delete()
	if err != nil {
		return nil, err
	}
	return &parking_order.ParkingOrderDeleteRes{Success: true}, nil
}
