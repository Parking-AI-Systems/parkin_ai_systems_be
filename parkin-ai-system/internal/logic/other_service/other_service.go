package other_service

import (
	"context"
	"fmt"
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

type sOthersService struct{}

func Init() {
	service.RegisterOthersService(&sOthersService{})
}
func init() {
	Init()
}

func (s *sOthersService) OthersServiceAdd(ctx context.Context, req *entity.OthersServiceAddReq) (*entity.OthersServiceAddRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}
	if gconv.String(user.Map()["role"]) != "admin" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Only admins can add services")
	}

	lot, err := dao.ParkingLots.Ctx(ctx).Where("id", req.LotId).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking parking lot")
	}
	if lot.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeParkingLotNotFound, "Parking lot not found")
	}

	if req.Name == "" {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Service name is required")
	}
	if req.Price <= 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Price must be positive")
	}
	if req.DurationMinutes <= 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Duration must be positive")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error starting transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	data := do.OthersService{
		LotId:           req.LotId,
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		DurationMinutes: req.DurationMinutes,
		IsActive:        req.IsActive,
		CreatedAt:       gtime.Now(),
	}
	lastId, err := dao.OthersService.Ctx(ctx).TX(tx).Data(data).InsertAndGetId()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating service")
	}

	notiData := do.Notifications{
		UserId:         userID,
		Type:           "others_service_added",
		Content:        fmt.Sprintf("Service #%d (%s) for parking lot #%d has been added.", lastId, req.Name, req.LotId),
		RelatedOrderId: lastId,
		IsRead:         false,
		CreatedAt:      gtime.Now(),
	}
	_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating notification")
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error committing transaction")
	}

	return &entity.OthersServiceAddRes{Id: lastId}, nil
}

func (s *sOthersService) OthersServiceList(ctx context.Context, req *entity.OthersServiceListReq) (*entity.OthersServiceListRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}

	m := dao.OthersService.Ctx(ctx).
		Fields("others_service.*, parking_lots.name as lot_name").
		LeftJoin("parking_lots", "parking_lots.id = others_service.lot_id").
		Where("parking_lots.deleted_at IS NULL")

	if req.LotId != 0 {
		m = m.Where("others_service.lot_id", req.LotId)
	}
	if req.IsActive {
		m = m.Where("others_service.is_active", true)
	}

	total, err := m.Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error counting services")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	m = m.Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize)

	var services []struct {
		entity.OthersService
		LotName string `json:"lot_name"`
	}
	err = m.Order("others_service.id DESC").Scan(&services)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving services")
	}

	list := make([]entity.OthersServiceItem, 0, len(services))
	for _, svc := range services {
		item := entity.OthersServiceItem{
			Id:              svc.Id,
			LotId:           svc.LotId,
			LotName:         svc.LotName,
			Name:            svc.Name,
			Description:     svc.Description,
			Price:           svc.Price,
			DurationMinutes: svc.DurationMinutes,
			IsActive:        svc.IsActive,
			CreatedAt:       svc.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		list = append(list, item)
	}

	return &entity.OthersServiceListRes{
		List:  list,
		Total: total,
	}, nil
}

func (s *sOthersService) OthersServiceGet(ctx context.Context, req *entity.OthersServiceGetReq) (*entity.OthersServiceItem, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}

	var svc struct {
		entity.OthersService
		LotName string `json:"lot_name"`
	}
	err = dao.OthersService.Ctx(ctx).
		Fields("others_service.*, parking_lots.name as lot_name").
		LeftJoin("parking_lots", "parking_lots.id = others_service.lot_id").
		Where("others_service.id", req.Id).
		Where("parking_lots.deleted_at IS NULL").
		Scan(&svc)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving service")
	}
	if svc.Id == 0 {
		return nil, gerror.NewCode(consts.CodeNotFound, "Service not found")
	}

	item := entity.OthersServiceItem{
		Id:              svc.Id,
		LotId:           svc.LotId,
		LotName:         svc.LotName,
		Name:            svc.Name,
		Description:     svc.Description,
		Price:           svc.Price,
		DurationMinutes: svc.DurationMinutes,
		IsActive:        svc.IsActive,
		CreatedAt:       svc.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return &item, nil
}

func (s *sOthersService) OthersServiceUpdate(ctx context.Context, req *entity.OthersServiceUpdateReq) (*entity.OthersServiceItem, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}
	if gconv.String(user.Map()["role"]) != "admin" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Only admins can update services")
	}

	svc, err := dao.OthersService.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking service")
	}
	if svc.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "Service not found")
	}

	if req.LotId != 0 {
		lot, err := dao.ParkingLots.Ctx(ctx).Where("id", req.LotId).One()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking parking lot")
		}
		if lot.IsEmpty() {
			return nil, gerror.NewCode(consts.CodeParkingLotNotFound, "Parking lot not found")
		}
	}

	if req.Name != "" && req.Name == "" {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Service name cannot be empty")
	}
	if req.Price < 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Price must be non-negative")
	}
	if req.DurationMinutes < 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Duration must be non-negative")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error starting transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	updateData := g.Map{}
	if req.LotId != 0 {
		updateData["lot_id"] = req.LotId
	}
	if req.Name != "" {
		updateData["name"] = req.Name
	}
	if req.Description != "" {
		updateData["description"] = req.Description
	}
	if req.Price >= 0 {
		updateData["price"] = req.Price
	}
	if req.DurationMinutes >= 0 {
		updateData["duration_minutes"] = req.DurationMinutes
	}
	if req.IsActive != nil {
		updateData["is_active"] = req.IsActive
	}

	_, err = dao.OthersService.Ctx(ctx).TX(tx).Data(updateData).Where("id", req.Id).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error updating service")
	}

	notiData := do.Notifications{
		UserId:         userID,
		Type:           "others_service_updated",
		Content:        fmt.Sprintf("Service #%d (%s) for parking lot #%d has been updated.", req.Id, req.Name, req.LotId),
		RelatedOrderId: req.Id,
		IsRead:         false,
		CreatedAt:      gtime.Now(),
	}
	_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating notification")
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error committing transaction")
	}

	var updatedSvc struct {
		entity.OthersService
		LotName string `json:"lot_name"`
	}
	err = dao.OthersService.Ctx(ctx).
		Fields("others_service.*, parking_lots.name as lot_name").
		LeftJoin("parking_lots", "parking_lots.id = others_service.lot_id").
		Where("others_service.id", req.Id).
		Scan(&updatedSvc)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving updated service")
	}

	item := entity.OthersServiceItem{
		Id:              updatedSvc.Id,
		LotId:           updatedSvc.LotId,
		LotName:         updatedSvc.LotName,
		Name:            updatedSvc.Name,
		Description:     updatedSvc.Description,
		Price:           updatedSvc.Price,
		DurationMinutes: updatedSvc.DurationMinutes,
		IsActive:        updatedSvc.IsActive,
		CreatedAt:       updatedSvc.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return &item, nil
}

func (s *sOthersService) OthersServiceDelete(ctx context.Context, req *entity.OthersServiceDeleteReq) (*entity.OthersServiceDeleteRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}
	if gconv.String(user.Map()["role"]) != "admin" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Only admins can delete services")
	}

	svc, err := dao.OthersService.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking service")
	}
	if svc.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "Service not found")
	}

	count, err := dao.OthersServiceOrders.Ctx(ctx).
		Where("service_id", req.Id).
		Where("status", "confirmed").
		Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking active orders")
	}
	if count > 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Cannot delete service with active orders")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error starting transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = dao.OthersService.Ctx(ctx).TX(tx).Where("id", req.Id).Delete()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error deleting service")
	}

	notiData := do.Notifications{
		UserId:         userID,
		Type:           "others_service_deleted",
		Content:        fmt.Sprintf("Service #%d for parking lot #%d has been deleted.", req.Id, svc.Map()["lot_id"]),
		RelatedOrderId: req.Id,
		IsRead:         false,
		CreatedAt:      gtime.Now(),
	}
	_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating notification")
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error committing transaction")
	}

	return &entity.OthersServiceDeleteRes{Message: "Service deleted successfully"}, nil
}
