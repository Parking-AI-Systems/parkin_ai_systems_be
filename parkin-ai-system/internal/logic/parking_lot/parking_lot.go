package parking_lot

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"parkin-ai-system/internal/consts"
	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/model/do"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

type sParkingLot struct{}

func Init() {
	service.RegisterParkingLot(&sParkingLot{})
}
func init() {
	Init()
}

func (s *sParkingLot) ParkingLotAdd(ctx context.Context, req *entity.ParkingLotAddReq) (*entity.ParkingLotAddRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to add a parking lot.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}
	if gconv.String(user.Map()["role"]) != consts.RoleAdmin {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Only admins can create parking lots.")
	}

	if req.Name == "" || req.Address == "" {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Please provide both a name and address for the parking lot.")
	}
	if req.TotalSlots <= 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "The number of parking slots must be greater than zero.")
	}
	if req.PricePerHour <= 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "The price per hour must be greater than zero.")
	}
	if req.Latitude == 0 || req.Longitude == 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Please provide valid latitude and longitude coordinates.")
	}

	if req.OpenTime == nil || req.CloseTime == nil {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Please specify both open and close times.")
	}
	if req.CloseTime.Before(req.OpenTime) {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "The close time must be later than the open time.")
	}
	for _, img := range req.Images {
		if img.ImageUrl == "" {
			return nil, gerror.NewCode(consts.CodeInvalidInput, "Please provide a valid image URL for each image.")
		}
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding the parking lot. Please try again later.")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	data := do.ParkingLots{
		Name:           req.Name,
		Address:        req.Address,
		Latitude:       req.Latitude,
		Longitude:      req.Longitude,
		OwnerId:        gconv.Int64(userID),
		IsVerified:     req.IsVerified,
		IsActive:       req.IsActive,
		TotalSlots:     req.TotalSlots,
		AvailableSlots: req.TotalSlots,
		PricePerHour:   req.PricePerHour,
		Description:    req.Description,
		OpenTime:       req.OpenTime,
		CloseTime:      req.CloseTime,
		CreatedAt:      gtime.Now(),
	}
	lastId, err := dao.ParkingLots.Ctx(ctx).TX(tx).Data(data).InsertAndGetId()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding the parking lot. Please try again later.")
	}

	for _, img := range req.Images {
		imgData := do.ParkingLotImages{
			LotId:     lastId,
			ImageUrl:  img.ImageUrl,
			CreatedAt: gtime.Now(),
		}
		_, err = dao.ParkingLotImages.Ctx(ctx).TX(tx).Data(imgData).Insert()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding the parking lot. Please try again later.")
		}
	}

	notiData := do.Notifications{
		UserId:         userID,
		Type:           "parking_lot_created",
		Content:        fmt.Sprintf("Parking lot #%d (%s) with %d images has been created successfully.", lastId, req.Name, len(req.Images)),
		RelatedOrderId: lastId,
		IsRead:         false,
		CreatedAt:      gtime.Now(),
	}
	_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding the parking lot. Please try again later.")
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding the parking lot. Please try again later.")
	}

	return &entity.ParkingLotAddRes{Id: lastId}, nil
}

func (s *sParkingLot) ParkingLotList(ctx context.Context, req *entity.ParkingLotListReq) (*entity.ParkingLotListRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to view parking lots.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}

	m := dao.ParkingLots.Ctx(ctx).Where("deleted_at IS NULL")
	if req.IsActive {
		m = m.Where("is_active", true)
	}

	if req.Latitude != 0 && req.Longitude != 0 && req.Radius > 0 {
		latDelta := req.Radius / 111.0
		lonDelta := req.Radius / (111.0 * gconv.Float64(math.Cos(gconv.Float64(req.Latitude)*math.Pi/180)))
		m = m.Where("latitude BETWEEN ? AND ?", req.Latitude-latDelta, req.Latitude+latDelta).
			Where("longitude BETWEEN ? AND ?", req.Longitude-lonDelta, req.Longitude+lonDelta)
	}

	total, err := m.Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load parking lots. Please try again later.")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	m = m.Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize)

	var lots []entity.ParkingLots
	err = m.Order("id DESC").Scan(&lots)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load parking lots. Please try again later.")
	}

	list := make([]entity.ParkingLotItem, 0, len(lots))
	for _, lot := range lots {
		var images []entity.ParkingLotImages
		err = dao.ParkingLotImages.Ctx(ctx).Where("lot_id", lot.Id).Where("deleted_at IS NULL").Scan(&images)
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load images for the parking lot. Please try again later.")
		}

		imageItems := make([]entity.ParkingLotImageItem, 0, len(images))
		for _, img := range images {
			item := entity.ParkingLotImageItem{
				Id:           img.Id,
				ParkingLotId: img.LotId,
				LotName:      lot.Name,
				ImageUrl:     img.ImageUrl,
				CreatedAt:    time.Time(img.CreatedAt.Time).Format("2006-01-02 15:04:05"),
			}
			imageItems = append(imageItems, item)
		}

		item := entity.ParkingLotItem{
			Id:             lot.Id,
			Name:           lot.Name,
			Address:        lot.Address,
			Latitude:       lot.Latitude,
			Longitude:      lot.Longitude,
			OwnerId:        lot.OwnerId,
			IsVerified:     lot.IsVerified,
			IsActive:       lot.IsActive,
			TotalSlots:     lot.TotalSlots,
			AvailableSlots: lot.AvailableSlots,
			PricePerHour:   lot.PricePerHour,
			Description:    lot.Description,
			OpenTime:       lot.OpenTime.Format("15:04:05"),
			CloseTime:      lot.CloseTime.Format("15:04:05"),
			Images:         imageItems,
			CreatedAt:      time.Time(lot.CreatedAt.Time).Format("2006-01-02 15:04:05"),
		}
		if !lot.UpdatedAt.IsZero() {
			item.UpdatedAt = time.Time(lot.UpdatedAt.Time).Format("2006-01-02 15:04:05")
		}
		list = append(list, item)
	}

	return &entity.ParkingLotListRes{
		List:  list,
		Total: total,
	}, nil
}

func (s *sParkingLot) ParkingLotGet(ctx context.Context, req *entity.ParkingLotGetReq) (*entity.ParkingLotItem, error) {
	// Uncomment if user authentication is required
	// userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	// if userID == "" {
	// 	return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to view parking lot details.")
	// }
	//
	// user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	// if err != nil {
	// 	return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	// }
	// if user.IsEmpty() {
	// 	return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	// }

	var lot entity.ParkingLots
	err := dao.ParkingLots.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").Scan(&lot)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load parking lot details. Please try again later.")
	}
	if lot.Id == 0 {
		return nil, gerror.NewCode(consts.CodeNotFound, "The parking lot could not be found.")
	}

	var images []entity.ParkingLotImages
	err = dao.ParkingLotImages.Ctx(ctx).Where("lot_id", lot.Id).Where("deleted_at IS NULL").Scan(&images)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load images for the parking lot. Please try again later.")
	}

	imageItems := make([]entity.ParkingLotImageItem, 0, len(images))
	for _, img := range images {
		item := entity.ParkingLotImageItem{
			Id:           img.Id,
			ParkingLotId: img.LotId,
			LotName:      lot.Name,
			ImageUrl:     img.ImageUrl,
			CreatedAt:    time.Time(img.CreatedAt.Time).Format("2006-01-02 15:04:05"),
		}
		imageItems = append(imageItems, item)
	}

	item := entity.ParkingLotItem{
		Id:             lot.Id,
		Name:           lot.Name,
		Address:        lot.Address,
		Latitude:       lot.Latitude,
		Longitude:      lot.Longitude,
		OwnerId:        lot.OwnerId,
		IsVerified:     lot.IsVerified,
		IsActive:       lot.IsActive,
		TotalSlots:     lot.TotalSlots,
		AvailableSlots: lot.AvailableSlots,
		PricePerHour:   lot.PricePerHour,
		Description:    lot.Description,
		OpenTime:       time.Time(lot.OpenTime.Time).Format("15:04:05"),
		CloseTime:      time.Time(lot.CloseTime.Time).Format("15:04:05"),
		Images:         imageItems,
		CreatedAt:      time.Time(lot.CreatedAt.Time).Format("2006-01-02 15:04:05"),
	}
	if !lot.UpdatedAt.IsZero() {
		item.UpdatedAt = time.Time(lot.UpdatedAt.Time).Format("2006-01-02 15:04:05")
	}

	return &item, nil
}

func (s *sParkingLot) ParkingLotUpdate(ctx context.Context, req *entity.ParkingLotUpdateReq) (*entity.ParkingLotItem, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to update parking lot details.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}
	if gconv.String(user.Map()["role"]) != consts.RoleAdmin {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Only admins can update parking lots.")
	}

	lot, err := dao.ParkingLots.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to find the parking lot. Please try again.")
	}
	if lot.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "The parking lot could not be found.")
	}

	if req.CloseTime != nil && req.OpenTime != nil && req.CloseTime.Before(req.OpenTime) {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "The close time must be later than the open time.")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the parking lot. Please try again later.")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	updateData := g.Map{
		"updated_at": gtime.Now(),
	}
	if req.Name != "" {
		updateData["name"] = req.Name
	}
	if req.Address != "" {
		updateData["address"] = req.Address
	}
	if req.Latitude != 0 {
		updateData["latitude"] = req.Latitude
	}
	if req.Longitude != 0 {
		updateData["longitude"] = req.Longitude
	}
	if req.IsVerified != nil {
		updateData["is_verified"] = req.IsVerified
	}
	if req.IsActive != nil {
		updateData["is_active"] = req.IsActive
	}
	if req.TotalSlots > 0 {
		currentAvailable := gconv.Int(lot.Map()["available_slots"])
		currentTotal := gconv.Int(lot.Map()["total_slots"])
		if req.TotalSlots < currentTotal-currentAvailable {
			return nil, gerror.NewCode(consts.CodeInvalidInput, "The number of slots cannot be less than those currently occupied.")
		}
		updateData["total_slots"] = req.TotalSlots
		updateData["available_slots"] = req.TotalSlots - (currentTotal - currentAvailable)
	}
	if req.PricePerHour > 0 {
		updateData["price_per_hour"] = req.PricePerHour
	}
	if req.Description != "" {
		updateData["description"] = req.Description
	}
	if req.OpenTime != nil {
		updateData["open_time"] = req.OpenTime
	}
	if req.CloseTime != nil {
		updateData["close_time"] = req.CloseTime
	}
	if req.ImageUrl != "" {
		updateData["image_url"] = req.ImageUrl
	}

	_, err = dao.ParkingLots.Ctx(ctx).TX(tx).Data(updateData).Where("id", req.Id).Where("deleted_at IS NULL").Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the parking lot. Please try again later.")
	}

	notiData := do.Notifications{
		UserId:         userID,
		Type:           "parking_lot_updated",
		Content:        fmt.Sprintf("Parking lot #%d (%s) has been updated.", req.Id, req.Name),
		RelatedOrderId: req.Id,
		IsRead:         false,
		CreatedAt:      gtime.Now(),
	}
	_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the parking lot. Please try again later.")
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the parking lot. Please try again later.")
	}

	var updatedLot entity.ParkingLots
	err = dao.ParkingLots.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").Scan(&updatedLot)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load updated parking lot details. Please try again later.")
	}

	var images []entity.ParkingLotImages
	err = dao.ParkingLotImages.Ctx(ctx).Where("lot_id", updatedLot.Id).Where("deleted_at IS NULL").Scan(&images)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load images for the parking lot. Please try again later.")
	}

	imageItems := make([]entity.ParkingLotImageItem, 0, len(images))
	for _, img := range images {
		item := entity.ParkingLotImageItem{
			Id:           img.Id,
			ParkingLotId: img.LotId,
			LotName:      updatedLot.Name,
			ImageUrl:     img.ImageUrl,
			CreatedAt:    time.Time(img.CreatedAt.Time).Format("2006-01-02 15:04:05"),
		}
		imageItems = append(imageItems, item)
	}

	item := entity.ParkingLotItem{
		Id:             updatedLot.Id,
		Name:           updatedLot.Name,
		Address:        updatedLot.Address,
		Latitude:       updatedLot.Latitude,
		Longitude:      updatedLot.Longitude,
		OwnerId:        updatedLot.OwnerId,
		IsVerified:     updatedLot.IsVerified,
		IsActive:       updatedLot.IsActive,
		TotalSlots:     updatedLot.TotalSlots,
		AvailableSlots: updatedLot.AvailableSlots,
		PricePerHour:   updatedLot.PricePerHour,
		Description:    updatedLot.Description,
		OpenTime:       updatedLot.OpenTime.Format("15:04:05"),
		CloseTime:      updatedLot.CloseTime.Format("15:04:05"),
		Images:         imageItems,
		CreatedAt:      time.Time(updatedLot.CreatedAt.Time).Format("2006-01-02 15:04:05"),
	}
	if !updatedLot.UpdatedAt.IsZero() {
		item.UpdatedAt = time.Time(updatedLot.UpdatedAt.Time).Format("2006-01-02 15:04:05")
	}

	return &item, nil
}

func (s *sParkingLot) ParkingLotDelete(ctx context.Context, req *entity.ParkingLotDeleteReq) (*entity.ParkingLotDeleteRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to delete a parking lot.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}
	if gconv.String(user.Map()["role"]) != consts.RoleAdmin {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Only admins can delete parking lots.")
	}

	lot, err := dao.ParkingLots.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to find the parking lot. Please try again.")
	}
	if lot.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "The parking lot could not be found.")
	}

	count, err := dao.ParkingOrders.Ctx(ctx).
		Where("lot_id", req.Id).
		Where("status", "confirmed").
		Where("deleted_at IS NULL").
		Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to check parking lot orders. Please try again later.")
	}
	if count > 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Cannot delete this parking lot because it has active bookings.")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the parking lot. Please try again later.")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = dao.ParkingSlots.Ctx(ctx).TX(tx).Data(g.Map{
		"deleted_at": gtime.Now(),
		"updated_at": gtime.Now(),
	}).Where("lot_id", req.Id).Where("deleted_at IS NULL").Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the parking lot. Please try again later.")
	}

	_, err = dao.ParkingLotReviews.Ctx(ctx).TX(tx).Data(g.Map{
		"deleted_at": gtime.Now(),
		"updated_at": gtime.Now(),
	}).Where("lot_id", req.Id).Where("deleted_at IS NULL").Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the parking lot. Please try again later.")
	}

	_, err = g.Model("others_service").Ctx(ctx).TX(tx).Data(g.Map{
		"deleted_at": gtime.Now(),
		"updated_at": gtime.Now(),
	}).Where("lot_id", req.Id).Where("deleted_at IS NULL").Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the parking lot. Please try again later.")
	}

	_, err = dao.ParkingLots.Ctx(ctx).TX(tx).Data(g.Map{
		"deleted_at": gtime.Now(),
		"updated_at": gtime.Now(),
	}).Where("id", req.Id).Where("deleted_at IS NULL").Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the parking lot. Please try again later.")
	}

	_, err = dao.ParkingLotImages.Ctx(ctx).TX(tx).Data(g.Map{
		"deleted_at": gtime.Now(),
		"updated_at": gtime.Now(),
	}).Where("lot_id", req.Id).Where("deleted_at IS NULL").Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the parking lot. Please try again later.")
	}

	notiData := do.Notifications{
		UserId:         userID,
		Type:           "parking_lot_deleted",
		Content:        fmt.Sprintf("Parking lot #%d has been deleted.", req.Id),
		RelatedOrderId: req.Id,
		IsRead:         false,
		CreatedAt:      gtime.Now(),
	}
	_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the parking lot. Please try again later.")
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the parking lot. Please try again later.")
	}

	return &entity.ParkingLotDeleteRes{Message: "Parking lot deleted successfully"}, nil
}

func (s *sParkingLot) ParkingLotImageDelete(ctx context.Context, req *entity.ParkingLotImageDeleteReq) (*entity.ParkingLotImageDeleteRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to delete a parking lot image.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}
	if gconv.String(user.Map()["role"]) != consts.RoleAdmin {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Only admins can delete parking lot images.")
	}

	img, err := dao.ParkingLotImages.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to find the image. Please try again.")
	}
	if img.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "The image could not be found.")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the image. Please try again later.")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = dao.ParkingLotImages.Ctx(ctx).TX(tx).Data(g.Map{
		"deleted_at": gtime.Now(),
		"updated_at": gtime.Now(),
	}).Where("id", req.Id).Where("deleted_at IS NULL").Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the image. Please try again later.")
	}

	notiData := do.Notifications{
		UserId:         userID,
		Type:           "parking_lot_image_deleted",
		Content:        fmt.Sprintf("Image #%d for parking lot #%d has been deleted.", req.Id, img.Map()["lot_id"]),
		RelatedOrderId: req.Id,
		IsRead:         false,
		CreatedAt:      gtime.Now(),
	}
	_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the image. Please try again later.")
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the image. Please try again later.")
	}

	return &entity.ParkingLotImageDeleteRes{Message: "Parking lot image deleted successfully"}, nil
}
