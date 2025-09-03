package other_service_order

import (
	"context"
	"parkin-ai-system/api/other_service_order"
	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/model/do"
	"parkin-ai-system/internal/model/entity"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/os/gtime"
	"parkin-ai-system/internal/service"
)

type sOtherServiceOrder struct{}

func init() {
	service.RegisterOtherServiceOrder(&sOtherServiceOrder{})
}

func (s *sOtherServiceOrder) OtherServiceOrderAdd(ctx context.Context, req *other_service_order.OtherServiceOrderAddReq) (*other_service_order.OtherServiceOrderAddRes, error) {
	userId := ctx.Value("user_id")
	order := do.OthersServiceOrders{}
	gconv.Struct(req, &order)
	order.UserId = userId
	order.CreatedAt = gtime.Now()
	lastId, err := dao.OthersServiceOrders.Ctx(ctx).Data(order).InsertAndGetId()
	if err != nil {
		return nil, err
	}
	return &other_service_order.OtherServiceOrderAddRes{OrderId: gconv.Int64(lastId)}, nil
}

func (s *sOtherServiceOrder) OtherServiceOrderUpdate(ctx context.Context, req *other_service_order.OtherServiceOrderUpdateReq) (*other_service_order.OtherServiceOrderUpdateRes, error) {
	data := do.OthersServiceOrders{}
	gconv.Struct(req, &data)
	_, err := dao.OthersServiceOrders.Ctx(ctx).Where("id", req.Id).Data(data).Update()
	if err != nil {
		return nil, err
	}
	return &other_service_order.OtherServiceOrderUpdateRes{Success: true}, nil
}

func (s *sOtherServiceOrder) OtherServiceOrderDelete(ctx context.Context, req *other_service_order.OtherServiceOrderDeleteReq) (*other_service_order.OtherServiceOrderDeleteRes, error) {
	_, err := dao.OthersServiceOrders.Ctx(ctx).Where("id", req.Id).Delete()
	if err != nil {
		return nil, err
	}
	return &other_service_order.OtherServiceOrderDeleteRes{Success: true}, nil
}

func (s *sOtherServiceOrder) OtherServiceOrderList(ctx context.Context, req *other_service_order.OtherServiceOrderListReq) (*other_service_order.OtherServiceOrderListRes, error) {
	var orders []entity.OthersServiceOrders
	m := dao.OthersServiceOrders.Ctx(ctx)
	if req.UserId != 0 {
		m = m.Where("user_id", req.UserId)
	}
	err := m.Order("id desc").Scan(&orders)
	if err != nil {
		return nil, err
	}
	var list []other_service_order.OtherServiceOrderItem
	for _, order := range orders {
		item := other_service_order.OtherServiceOrderItem{}
		gconv.Struct(order, &item)
		item.CreatedAt = order.CreatedAt.Format("2006-01-02 15:04:05")
		list = append(list, item)
	}
	return &other_service_order.OtherServiceOrderListRes{List: list}, nil
}
