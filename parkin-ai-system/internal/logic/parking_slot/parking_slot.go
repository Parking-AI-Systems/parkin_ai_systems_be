package parking_slot

import (
	"context"
	"parkin-ai-system/api/parking_slot"
	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/model/do"
	"parkin-ai-system/internal/model/entity"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type sParkingSlot struct{}

func init() {
	// Đăng ký logic với service
	service.RegisterParkingSlot(&sParkingSlot{})
}

func (s *sParkingSlot) ParkingSlotAdd(req *parking_slot.ParkingSlotAddReq) (*parking_slot.ParkingSlotAddRes, error) {
	data := do.ParkingSlots{}
	gconv.Struct(req, &data)
	lastId, err := dao.ParkingSlots.Ctx(context.Background()).Data(data).InsertAndGetId()
	if err != nil {
		return nil, err
	}
	return &parking_slot.ParkingSlotAddRes{Id: gconv.Int64(lastId)}, nil
}

func (s *sParkingSlot) ParkingSlotList(req *parking_slot.ParkingSlotListReq) (*parking_slot.ParkingSlotListRes, error) {
	var slots []entity.ParkingSlots
	m := dao.ParkingSlots.Ctx(context.Background())
	if req.LotId != 0 {
		m = m.Where("lot_id", req.LotId)
	}
	err := m.Order("id desc").Scan(&slots)
	if err != nil {
		return nil, err
	}
	var list []parking_slot.ParkingSlotItem
	for _, slot := range slots {
		item := parking_slot.ParkingSlotItem{}
		gconv.Struct(slot, &item)
		item.CreatedAt = slot.CreatedAt.Format("2006-01-02 15:04:05")
		list = append(list, item)
	}
	return &parking_slot.ParkingSlotListRes{List: list}, nil
}

func (s *sParkingSlot) ParkingSlotUpdate(req *parking_slot.ParkingSlotUpdateReq) (*parking_slot.ParkingSlotUpdateRes, error) {
	data := do.ParkingSlots{}
	gconv.Struct(req, &data)
	_, err := dao.ParkingSlots.Ctx(context.Background()).Where("id", req.Id).Data(data).Update()
	if err != nil {
		return nil, err
	}
	return &parking_slot.ParkingSlotUpdateRes{Success: true}, nil
}

func (s *sParkingSlot) ParkingSlotDelete(req *parking_slot.ParkingSlotDeleteReq) (*parking_slot.ParkingSlotDeleteRes, error) {
	_, err := dao.ParkingSlots.Ctx(context.Background()).Where("id", req.Id).Delete()
	if err != nil {
		return nil, err
	}
	return &parking_slot.ParkingSlotDeleteRes{Success: true}, nil
}
