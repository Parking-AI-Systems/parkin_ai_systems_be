package parking_lot

import (
	"context"
	"parkin-ai-system/api/parking_lot/parking_lot"
	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/model/do"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/guid"
)

type sParkingLot struct{}

func Init() {
	service.RegisterParkingLot(&sParkingLot{})
}

func init() {
	Init()
}

func (s *sParkingLot) ParkingLotAdd(ctx context.Context, req *parking_lot.ParkingLotAddReq) (res *parking_lot.ParkingLotAddRes, err error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	userRole := g.RequestFromCtx(ctx).GetCtxVar("user_role").String()
	if userID == "" {
		return nil, gerror.New("Unauthorized")
	}
	if userRole != "role_admin" {
		return nil, gerror.New("Not admin")
	}

	exists, err := dao.ParkingLots.Ctx(ctx).Where("latitude = ? AND longitude = ?", req.Latitude, req.Longitude).Count()
	if err != nil {
		return nil, gerror.New("Database error")
	}
	if exists > 0 {
		return nil, gerror.New("Location already exists")
	}

	lotId := guid.S()
	_, err = dao.ParkingLots.Ctx(ctx).Data(do.ParkingLots{
		Name:           req.Name,
		Address:        req.Address,
		Latitude:       req.Latitude,
		Longitude:      req.Longitude,
		TotalSlots:     req.TotalSlots,
		AvailableSlots: req.TotalSlots,
		IsActive:       true,
		IsVerified:     false,
	}).Insert()
	if err != nil {
		return nil, gerror.New("Database error")
	}
	// TODO: Auto-create slots based on TotalSlots
	res = &parking_lot.ParkingLotAddRes{LotID: lotId}
	return
}

func (s *sParkingLot) ParkingLotDetail(ctx context.Context, req *parking_lot.ParkingLotDetailReq) (res *parking_lot.ParkingLotDetailRes, err error) {
	lot, err := dao.ParkingLots.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, gerror.New("Database error")
	}
	if lot.IsEmpty() {
		return nil, gerror.New("Not found")
	}

	slots, err := dao.ParkingSlots.Ctx(ctx).Where("lot_id", req.Id).All()
	if err != nil {
		return nil, gerror.New("Database error")
	}
	images, err := dao.ParkingLotImages.Ctx(ctx).Where("lot_id", req.Id).All()
	if err != nil {
		return nil, gerror.New("Database error")
	}
	reviews, err := dao.ParkingLotReviews.Ctx(ctx).Where("lot_id", req.Id).All()
	if err != nil {
		return nil, gerror.New("Database error")
	}

	lotInfo := &parking_lot.ParkingLotInfo{
		Id:             lot["id"].String(),
		Name:           lot["name"].String(),
		Address:        lot["address"].String(),
		Latitude:       lot["latitude"].Float64(),
		Longitude:      lot["longitude"].Float64(),
		TotalSlots:     lot["total_slots"].Int(),
		AvailableSlots: lot["available_slots"].Int(),
		Description:    lot["description"].String(),
	}
	slotList := make([]parking_lot.ParkingSlotInfo, 0, len(slots))
	for _, s := range slots {
		slotList = append(slotList, parking_lot.ParkingSlotInfo{
			Id:         s["id"].String(),
			SlotNumber: s["slot_number"].String(),
			Status:     s["status"].String(),
		})
	}
	imageList := make([]parking_lot.ParkingLotImage, 0, len(images))
	for _, img := range images {
		imageList = append(imageList, parking_lot.ParkingLotImage{
			Id:       img["id"].String(),
			ImageUrl: img["image_url"].String(),
		})
	}
	reviewList := make([]parking_lot.ParkingLotReview, 0, len(reviews))
	for _, r := range reviews {
		reviewList = append(reviewList, parking_lot.ParkingLotReview{
			Id:      r["id"].String(),
			Score:   r["score"].Int(),
			Comment: r["comment"].String(),
		})
	}

	res = &parking_lot.ParkingLotDetailRes{
		Lot:     lotInfo,
		Slots:   slotList,
		Images:  imageList,
		Reviews: reviewList,
	}
	return
}
