package parking_order

import (
	"context"
	"parkin-ai-system/api/parking_order"
	"parkin-ai-system/internal/consts"
	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/model/do"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type sParkingOrder struct{}

func Init() {
	service.RegisterParkingOrder(&sParkingOrder{})
}
func init() {
	Init()
}

func (s *sParkingOrder) ParkingOrderAdd(req *parking_order.ParkingOrderAddReq) (*parking_order.ParkingOrderAddRes, error) {
	return nil, nil
}
func (s *sParkingOrder) ParkingOrderAddWithUser(ctx context.Context, req *parking_order.ParkingOrderAddReq) (*parking_order.ParkingOrderAddRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	user, _ := dao.Users.Ctx(ctx).Where("id", userID).One()
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound)
	}
	vehicle, _ := dao.Vehicles.Ctx(ctx).Where("id", req.VehicleId).One()
	if vehicle.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeVehicleNotFound)
	}
	lot, _ := dao.ParkingLots.Ctx(ctx).Where("id", req.LotId).One()
	if lot.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeParkingLotNotFound)
	}
	slot, _ := dao.ParkingSlots.Ctx(ctx).Where("id", req.SlotId).One()
	if slot.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeParkingSlotNotFound)
	}
	data := do.ParkingOrders{}
	gconv.Struct(req, &data)
	data.UserId = userID
	lastId, err := dao.ParkingOrders.Ctx(ctx).Data(data).InsertAndGetId()
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
	data.UpdatedAt = gtime.Now()
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
