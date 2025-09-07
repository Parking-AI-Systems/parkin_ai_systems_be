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
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to add a service")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while checking user details")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}
	if gconv.String(user.Map()["role"]) != consts.RoleAdmin {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Only admins can add services")
	}

	lot, err := dao.ParkingLots.Ctx(ctx).Where("id", req.LotId).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while checking the parking lot")
	}
	if lot.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeParkingLotNotFound, "Parking lot not found")
	}

	if req.Name == "" {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Please provide a service name")
	}
	if req.Price <= 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Service price must be greater than 0")
	}
	if req.DurationMinutes <= 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Service duration must be greater than 0")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while starting the transaction")
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
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding the service")
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
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while creating the notification")
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while saving the changes")
	}

	return &entity.OthersServiceAddRes{Id: lastId}, nil
}

func (s *sOthersService) OthersServiceList(ctx context.Context, req *entity.OthersServiceListReq) (*entity.OthersServiceListRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to view the service list")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while checking user details")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}

	// Build base query conditions
	baseQuery := dao.OthersService.Ctx(ctx).
		LeftJoin("parking_lots", "parking_lots.id = others_service.lot_id").
		Where("others_service.deleted_at IS NULL").
		Where("parking_lots.deleted_at IS NULL")

	if req.LotId != 0 {
		baseQuery = baseQuery.Where("others_service.lot_id", req.LotId)
	}
	if req.IsActive {
		baseQuery = baseQuery.Where("others_service.is_active", true)
	}

	// Count query - use simple field
	total, err := baseQuery.Fields("others_service.id").Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while counting services")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// Data query - use joined fields
	var services []struct {
		entity.OthersService
		LotName string `json:"lot_name"`
	}
	err = baseQuery.Fields("others_service.*, parking_lots.name as lot_name").
		Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize).
		Order("others_service.id DESC").Scan(&services)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while retrieving the service list")
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
			UpdatedAt:       svc.UpdatedAt.Format("2006-01-02 15:04:05"),
			DeletedAt:       svc.DeletedAt.Format("2006-01-02 15:04:05"),
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
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to view service details")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while checking user details")
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
		Where("others_service.deleted_at IS NULL").
		Where("parking_lots.deleted_at IS NULL").
		Scan(&svc)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while retrieving service details")
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
		UpdatedAt:       svc.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt:       svc.DeletedAt.Format("2006-01-02 15:04:05"),
	}

	return &item, nil
}

func (s *sOthersService) OthersServiceUpdate(ctx context.Context, req *entity.OthersServiceUpdateReq) (*entity.OthersServiceItem, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to update the service")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while checking user details")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}
	if gconv.String(user.Map()["role"]) != consts.RoleAdmin {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Only admins can update services")
	}

	svc, err := dao.OthersService.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while checking the service")
	}
	if svc.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "Service not found")
	}

	if req.LotId != 0 {
		lot, err := dao.ParkingLots.Ctx(ctx).Where("id", req.LotId).Where("deleted_at IS NULL").One()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while checking the parking lot")
		}
		if lot.IsEmpty() {
			return nil, gerror.NewCode(consts.CodeParkingLotNotFound, "Parking lot not found")
		}
	}

	if req.Name != "" && req.Name == "" {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Service name cannot be empty")
	}
	if req.Price < 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Service price cannot be negative")
	}
	if req.DurationMinutes < 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Service duration cannot be negative")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while starting the transaction")
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
	updateData["updated_at"] = gtime.Now()
	_, err = dao.OthersService.Ctx(ctx).TX(tx).Data(updateData).Where("id", req.Id).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the service")
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
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while creating the notification")
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while saving the changes")
	}

	var updatedSvc struct {
		entity.OthersService
		LotName string `json:"lot_name"`
	}
	err = dao.OthersService.Ctx(ctx).
		Fields("others_service.*, parking_lots.name as lot_name").
		LeftJoin("parking_lots", "parking_lots.id = others_service.lot_id").
		Where("others_service.id", req.Id).
		Where("others_service.deleted_at IS NULL").
		Where("parking_lots.deleted_at IS NULL").
		Scan(&updatedSvc)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while retrieving the updated service")
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
		UpdatedAt:       updatedSvc.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt:       updatedSvc.DeletedAt.Format("2006-01-02 15:04:05"),
	}

	return &item, nil
}

func (s *sOthersService) OthersServiceDelete(ctx context.Context, req *entity.OthersServiceDeleteReq) (*entity.OthersServiceDeleteRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to delete the service")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while checking user details")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}
	if gconv.String(user.Map()["role"]) != consts.RoleAdmin {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Only admins can delete services")
	}

	svc, err := dao.OthersService.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while checking the service")
	}
	if svc.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "Service not found")
	}

	count, err := dao.OthersServiceOrders.Ctx(ctx).
		Where("service_id", req.Id).
		Where("status", "confirmed").
		Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while checking active orders")
	}
	if count > 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Cannot delete the service because it has active orders")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while starting the transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = dao.OthersService.Ctx(ctx).TX(tx).Data(g.Map{
		"deleted_at": gtime.Now(),
	}).Where("id", req.Id).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the service")
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
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while creating the notification")
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while saving the changes")
	}

	return &entity.OthersServiceDeleteRes{Message: "Service deleted successfully"}, nil
}
