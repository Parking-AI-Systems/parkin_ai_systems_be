package parking_order

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"parkin-ai-system/api/parking_order/parking_order"
)

func (c *ControllerParking_order) MyParkingLotOrderList(ctx context.Context, req *parking_order.MyParkingLotOrderListReq) (res *parking_order.MyParkingLotOrderListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
