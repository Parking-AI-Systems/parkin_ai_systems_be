
package parking_lot

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/frame/g"

	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/model/do"
	"parkin-ai-system/internal/service"

	api_add "parkin-ai-system/api/parking_lot/parking_lot"
	api_detail "parkin-ai-system/api/parking_lot/parking_lot"
	api_list "parkin-ai-system/api/parking_lot/parking_lot"
	api_update "parkin-ai-system/api/parking_lot/parking_lot"
	api_delete "parkin-ai-system/api/parking_lot/parking_lot"
)


func (s *sParkingLot) ParkingLotList(ctx context.Context, req *api_list.ParkingLotListReq) (res *api_list.ParkingLotListRes, err error) {
	var lots []api_detail.ParkingLotInfo
	err = dao.ParkingLots.Ctx(ctx).Scan(&lots)
	if err != nil {
		return nil, gerror.New("Database error")
	}
	res = &api_list.ParkingLotListRes{Lots: lots}
	return
}

func (s *sParkingLot) ParkingLotUpdate(ctx context.Context, req *api_update.ParkingLotUpdateReq) (res *api_update.ParkingLotUpdateRes, err error) {
	userRole := g.RequestFromCtx(ctx).GetCtxVar("user_role").String()
	if userRole != "role_admin" {
		return nil, gerror.New("Not admin")
	}
       data := do.ParkingLots{
	       Name:         req.Name,
	       Address:      req.Address,
	       Latitude:     req.Latitude,
	       Longitude:    req.Longitude,
	       PricePerHour: req.PricePerHour,
	       Description:  req.Description,
	       OpenTime:     gtime.NewFromStr(req.OpenTime),
	       CloseTime:    gtime.NewFromStr(req.CloseTime),
	       ImageUrl:     req.ImageUrl,
	       IsActive:     req.IsActive,
	       IsVerified:   req.IsVerified,
	       TotalSlots:   req.TotalSlots,
	       AvailableSlots: req.AvailableSlots,
       }
	_, err = dao.ParkingLots.Ctx(ctx).Where("id", req.Id).Data(data).Update()
	if err != nil {
		return nil, gerror.New("Database error")
	}
	res = &api_update.ParkingLotUpdateRes{Success: true}
	return
}

func (s *sParkingLot) ParkingLotDelete(ctx context.Context, req *api_delete.ParkingLotDeleteReq) (res *api_delete.ParkingLotDeleteRes, err error) {
	userRole := g.RequestFromCtx(ctx).GetCtxVar("user_role").String()
	if userRole != "role_admin" {
		return nil, gerror.New("Not admin")
	}
	_, err = dao.ParkingLots.Ctx(ctx).Where("id", req.Id).Delete()
	if err != nil {
		return nil, gerror.New("Database error")
	}
	res = &api_delete.ParkingLotDeleteRes{Success: true}
	return
}

func (s *sParkingLot) ParkingLotAdd(ctx context.Context, req *api_add.ParkingLotAddReq) (res *api_add.ParkingLotAddRes, err error) {
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

	// Lấy owner_id từ token
	ownerID := g.RequestFromCtx(ctx).GetCtxVar("user_id").Int64()
	lotId := guid.S()
       _, err = dao.ParkingLots.Ctx(ctx).Data(do.ParkingLots{
	       Name:           req.Name,
	       Address:        req.Address,
	       Latitude:       req.Latitude,
	       Longitude:      req.Longitude,
	       OwnerId:        ownerID,
	       IsVerified:     false,
	       IsActive:       true,
	       TotalSlots:     req.TotalSlots,
	       AvailableSlots: req.TotalSlots,
	       PricePerHour:   0, // req.PricePerHour không có trong struct gốc
	       Description:    "", // req.Description không có trong struct gốc
	       OpenTime:       gtime.NewFromStr("") , // req.OpenTime không có trong struct gốc
	       CloseTime:      gtime.NewFromStr("") , // req.CloseTime không có trong struct gốc
	       ImageUrl:       "", // req.ImageUrl không có trong struct gốc
       }).Insert()
	if err != nil {
		return nil, gerror.New("Database error")
	}
	// TODO: Auto-create slots based on TotalSlots
	res = &api_add.ParkingLotAddRes{LotID: lotId}
	return
}

func (s *sParkingLot) ParkingLotDetail(ctx context.Context, req *api_detail.ParkingLotDetailReq) (res *api_detail.ParkingLotDetailRes, err error) {
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

	lotInfo := &api_detail.ParkingLotInfo{
	       Id:             lot["id"].String(),
	       Name:           lot["name"].String(),
	       Address:        lot["address"].String(),
	       Latitude:       lot["latitude"].Float64(),
	       Longitude:      lot["longitude"].Float64(),
	       TotalSlots:     lot["total_slots"].Int(),
	       AvailableSlots: lot["available_slots"].Int(),
	       Description:    lot["description"].String(),
       }
	slotList := make([]api_detail.ParkingSlotInfo, 0, len(slots))
       for _, s := range slots {
			   slotList = append(slotList, api_detail.ParkingSlotInfo{
		       Id:         s["id"].String(),
		       SlotNumber: s["slot_number"].String(),
		       Status:     s["status"].String(),
	       })
       }
	imageList := make([]api_detail.ParkingLotImage, 0, len(images))
       for _, img := range images {
			   imageList = append(imageList, api_detail.ParkingLotImage{
		       Id:       img["id"].String(),
		       ImageUrl: img["image_url"].String(),
	       })
       }
	reviewList := make([]api_detail.ParkingLotReview, 0, len(reviews))
       for _, r := range reviews {
			   reviewList = append(reviewList, api_detail.ParkingLotReview{
		       Id:      r["id"].String(),
		       Score:   r["score"].Int(),
		       Comment: r["comment"].String(),
	       })
       }

	res = &api_detail.ParkingLotDetailRes{
	       Lot:     lotInfo,
	       Slots:   slotList,
	       Images:  imageList,
	       Reviews: reviewList,
       }
       return
}

func Init() {
	service.RegisterParkingLot(&sParkingLot{})
}

func init() {
	Init()
}
type sParkingLot struct{}
