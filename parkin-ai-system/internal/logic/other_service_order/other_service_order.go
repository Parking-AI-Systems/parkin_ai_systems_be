package other_service_order

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

type sOthersServiceOrder struct{}

func Init() {
	service.RegisterOthersServiceOrder(&sOthersServiceOrder{})
}
func init() {
	Init()
}

func (s *sOthersServiceOrder) OthersServiceOrderAddWithUser(ctx context.Context, req *entity.OthersServiceOrderAddReq) (*entity.OthersServiceOrderAddRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to create a service order")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at is NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while checking user details")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}

	vehicle, err := dao.Vehicles.Ctx(ctx).Where("id", req.VehicleId).Where("deleted_at is NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while checking vehicle details")
	}
	if vehicle.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeVehicleNotFound, "Vehicle not found")
	}

	lot, err := dao.ParkingLots.Ctx(ctx).Where("id", req.LotId).Where("deleted_at is NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while checking parking lot details")
	}
	if lot.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeParkingLotNotFound, "Parking lot not found")
	}

	service, err := dao.OthersService.Ctx(ctx).Where("id", req.ServiceId).Where("deleted_at is NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while checking service details")
	}
	if service.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeServiceNotFound, "Service not found")
	}
	isActive := gconv.Bool(service.Map()["is_active"])
	if !isActive {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "This service is not currently available")
	}

	if req.ScheduledTime == "" {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Please provide a scheduled time")
	}
	scheduledTime, err := gtime.StrToTime(req.ScheduledTime)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Invalid scheduled time format")
	}
	if scheduledTime.Before(gtime.Now()) {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Scheduled time must be in the future")
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

	data := do.OthersServiceOrders{
		UserId:        userID,
		LotId:         req.LotId,
		ServiceId:     req.ServiceId,
		VehicleId:     req.VehicleId,
		ScheduledTime: scheduledTime,
		Status:        "confirmed",
		Price:         gconv.Float64(service.Map()["price"]),
		PaymentStatus: "pending",
		CreatedAt:     gtime.Now(),
	}
	lastId, err := dao.OthersServiceOrders.Ctx(ctx).TX(tx).Data(data).InsertAndGetId()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while creating the service order")
	}

	notiData := do.Notifications{
		UserId:         userID,
		Type:           "service_order_created",
		Content:        fmt.Sprintf("Service order #%d has been created successfully.", lastId),
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

	return &entity.OthersServiceOrderAddRes{Id: lastId}, nil
}

func (s *sOthersServiceOrder) OthersServiceOrderList(ctx context.Context, req *entity.OthersServiceOrderListReq) (*entity.OthersServiceOrderListRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to view service orders")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at is NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while checking user details")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}
	isAdmin := gconv.String(user.Map()["role"]) == consts.RoleAdmin

	m := dao.OthersServiceOrders.Ctx(ctx).
		Fields("others_service_orders.*, parking_lots.name as lot_name, others_service.name as service_name, vehicles.license_plate as vehicle_plate").
		LeftJoin("parking_lots", "parking_lots.id = others_service_orders.lot_id").
		LeftJoin("others_service", "others_service.id = others_service_orders.service_id").
		LeftJoin("vehicles", "vehicles.id = others_service_orders.vehicle_id")

	if req.UserId != 0 {
		if !isAdmin && gconv.Int64(userID) != req.UserId {
			return nil, gerror.NewCode(consts.CodeUnauthorized, "You cannot view orders of other users")
		}
		m = m.Where("others_service_orders.user_id", req.UserId)
	} else if !isAdmin {
		m = m.Where("others_service_orders.user_id", userID)
	}
	if req.LotId != 0 {
		m = m.Where("others_service_orders.lot_id", req.LotId)
	}
	if req.Status != "" {
		m = m.Where("others_service_orders.status", req.Status)
	}

	total, err := m.Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while counting orders")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	m = m.Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize)

	var orders []struct {
		entity.OthersServiceOrders
		LotName      string `json:"lot_name"`
		ServiceName  string `json:"service_name"`
		VehiclePlate string `json:"vehicle_plate"`
	}
	err = m.Order("others_service_orders.id DESC").Scan(&orders)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while retrieving orders")
	}

	list := make([]entity.OthersServiceOrderItem, 0, len(orders))
	for _, order := range orders {
		item := entity.OthersServiceOrderItem{
			Id:            order.Id,
			UserId:        order.UserId,
			LotId:         order.LotId,
			ServiceId:     order.ServiceId,
			VehicleId:     order.VehicleId,
			LotName:       order.LotName,
			ServiceName:   order.ServiceName,
			VehiclePlate:  order.VehiclePlate,
			ScheduledTime: order.ScheduledTime.Format("2006-01-02 15:04:05"),
			Status:        order.Status,
			Price:         order.Price,
			PaymentStatus: order.PaymentStatus,
			CreatedAt:     order.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     order.UpdatedAt.Format("2006-01-02 15:04:05"),
			DeletedAt:     order.DeletedAt.Format("2006-01-02 15:04:05"),
		}
		list = append(list, item)
	}

	return &entity.OthersServiceOrderListRes{
		List:  list,
		Total: total,
	}, nil
}

func (s *sOthersServiceOrder) OthersServiceOrderGet(ctx context.Context, req *entity.OthersServiceOrderGetReq) (*entity.OthersServiceOrderItem, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to view order details")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while checking user details")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}
	isAdmin := gconv.String(user.Map()["role"]) == consts.RoleAdmin

	var order struct {
		entity.OthersServiceOrders
		LotName      string `json:"lot_name"`
		ServiceName  string `json:"service_name"`
		VehiclePlate string `json:"vehicle_plate"`
	}
	m := dao.OthersServiceOrders.Ctx(ctx).
		Fields("others_service_orders.*, parking_lots.name as lot_name, others_service.name as service_name, vehicles.license_plate as vehicle_plate").
		LeftJoin("parking_lots", "parking_lots.id = others_service_orders.lot_id").
		LeftJoin("others_service", "others_service.id = others_service_orders.service_id").
		LeftJoin("vehicles", "vehicles.id = others_service_orders.vehicle_id").
		Where("others_service_orders.id", req.Id)
	if !isAdmin {
		m = m.Where("others_service_orders.user_id", userID)
	}
	err = m.Scan(&order)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while retrieving order details")
	}
	if order.Id == 0 {
		return nil, gerror.NewCode(consts.CodeNotFound, "Service order not found")
	}

	item := entity.OthersServiceOrderItem{
		Id:            order.Id,
		UserId:        order.UserId,
		LotId:         order.LotId,
		ServiceId:     order.ServiceId,
		VehicleId:     order.VehicleId,
		LotName:       order.LotName,
		ServiceName:   order.ServiceName,
		VehiclePlate:  order.VehiclePlate,
		ScheduledTime: order.ScheduledTime.Format("2006-01-02 15:04:05"),
		Status:        order.Status,
		Price:         order.Price,
		PaymentStatus: order.PaymentStatus,
		CreatedAt:     order.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     order.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt:     order.DeletedAt.Format("2006-01-02 15:04:05"),
	}

	return &item, nil
}

func (s *sOthersServiceOrder) OthersServiceOrderUpdate(ctx context.Context, req *entity.OthersServiceOrderUpdateReq) (*entity.OthersServiceOrderItem, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to update the order")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while checking user details")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}

	isAdmin := gconv.String(user.Map()["role"]) == consts.RoleAdmin

	order, err := dao.OthersServiceOrders.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while checking order details")
	}
	if order.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "Service order not found")
	}
	if !isAdmin && gconv.Int64(order.Map()["user_id"]) != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "You cannot update orders of other users")
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

	updateData := g.Map{
		"updated_at": gtime.Now(),
	}
	if req.ScheduledTime != "" {
		scheduledTime, err := gtime.StrToTime(req.ScheduledTime)
		if err != nil {
			return nil, gerror.NewCode(consts.CodeInvalidInput, "Invalid scheduled time format")
		}
		if scheduledTime.Before(gtime.Now()) {
			return nil, gerror.NewCode(consts.CodeInvalidInput, "Scheduled time must be in the future")
		}
		updateData["scheduled_time"] = scheduledTime
	}
	if req.Status != "" {
		updateData["status"] = req.Status
	}

	_, err = dao.OthersServiceOrders.Ctx(ctx).TX(tx).Data(updateData).Where("id", req.Id).Where("deleted_at IS NULL").Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the order")
	}

	if req.Status != "" {
		notiData := do.Notifications{
			UserId:         gconv.String(order.Map()["user_id"]),
			Type:           "service_order_updated",
			Content:        fmt.Sprintf("Service order #%d status updated to %s.", req.Id, req.Status),
			RelatedOrderId: req.Id,
			IsRead:         false,
			CreatedAt:      gtime.Now(),
		}
		_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while creating the notification")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while saving the changes")
	}

	var updatedOrder struct {
		entity.OthersServiceOrders
		LotName      string `json:"lot_name"`
		ServiceName  string `json:"service_name"`
		VehiclePlate string `json:"vehicle_plate"`
	}
	err = dao.OthersServiceOrders.Ctx(ctx).
		Fields("others_service_orders.*, parking_lots.name as lot_name, others_service.name as service_name, vehicles.license_plate as vehicle_plate").
		LeftJoin("parking_lots", "parking_lots.id = others_service_orders.lot_id").
		LeftJoin("others_service", "others_service.id = others_service_orders.service_id").
		LeftJoin("vehicles", "vehicles.id = others_service_orders.vehicle_id").
		Where("others_service_orders.id", req.Id).
		Scan(&updatedOrder)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while retrieving the updated order")
	}

	item := entity.OthersServiceOrderItem{
		Id:            updatedOrder.Id,
		UserId:        updatedOrder.UserId,
		LotId:         updatedOrder.LotId,
		ServiceId:     updatedOrder.ServiceId,
		VehicleId:     updatedOrder.VehicleId,
		LotName:       updatedOrder.LotName,
		ServiceName:   updatedOrder.ServiceName,
		VehiclePlate:  updatedOrder.VehiclePlate,
		ScheduledTime: updatedOrder.ScheduledTime.Format("2006-01-02 15:04:05"),
		Status:        updatedOrder.Status,
		Price:         updatedOrder.Price,
		PaymentStatus: updatedOrder.PaymentStatus,
		CreatedAt:     updatedOrder.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     updatedOrder.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt:     updatedOrder.DeletedAt.Format("2006-01-02 15:04:05"),
	}

	return &item, nil
}

func (s *sOthersServiceOrder) OthersServiceOrderCancel(ctx context.Context, req *entity.OthersServiceOrderCancelReq) (*entity.OthersServiceOrderItem, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to cancel the order")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while checking user details")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}

	isAdmin := gconv.String(user.Map()["role"]) == consts.RoleAdmin

	order, err := dao.OthersServiceOrders.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while checking order details")
	}
	if order.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "Service order not found")
	}
	if !isAdmin && gconv.Int64(order.Map()["user_id"]) != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "You cannot cancel orders of other users")
	}
	if gconv.String(order.Map()["status"]) == "canceled" {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "This order has already been canceled")
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

	_, err = dao.OthersServiceOrders.Ctx(ctx).TX(tx).Data(g.Map{
		"status":     "canceled",
		"updated_at": gtime.Now(),
	}).Where("id", req.Id).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while canceling the order")
	}

	notiData := do.Notifications{
		UserId:         userID,
		Type:           "service_order_canceled",
		Content:        fmt.Sprintf("Service order #%d has been canceled.", req.Id),
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

	var updatedOrder struct {
		entity.OthersServiceOrders
		LotName      string `json:"lot_name"`
		ServiceName  string `json:"service_name"`
		VehiclePlate string `json:"vehicle_plate"`
	}
	err = dao.OthersServiceOrders.Ctx(ctx).
		Fields("others_service_orders.*, parking_lots.name as lot_name, others_service.name as service_name, vehicles.license_plate as vehicle_plate").
		LeftJoin("parking_lots", "parking_lots.id = others_service_orders.lot_id").
		LeftJoin("others_service", "others_service.id = others_service_orders.service_id").
		LeftJoin("vehicles", "vehicles.id = others_service_orders.vehicle_id").
		Where("others_service_orders.id", req.Id).
		Scan(&updatedOrder)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while retrieving the updated order")
	}

	item := entity.OthersServiceOrderItem{
		Id:            updatedOrder.Id,
		UserId:        updatedOrder.UserId,
		LotId:         updatedOrder.LotId,
		ServiceId:     updatedOrder.ServiceId,
		VehicleId:     updatedOrder.VehicleId,
		LotName:       updatedOrder.LotName,
		ServiceName:   updatedOrder.ServiceName,
		VehiclePlate:  updatedOrder.VehiclePlate,
		ScheduledTime: updatedOrder.ScheduledTime.Format("2006-01-02 15:04:05"),
		Status:        updatedOrder.Status,
		Price:         updatedOrder.Price,
		PaymentStatus: updatedOrder.PaymentStatus,
		CreatedAt:     updatedOrder.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     updatedOrder.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt:     updatedOrder.DeletedAt.Format("2006-01-02 15:04:05"),
	}

	return &item, nil
}

func (s *sOthersServiceOrder) OthersServiceOrderDelete(ctx context.Context, req *entity.OthersServiceOrderDeleteReq) (*entity.OthersServiceOrderDeleteRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to delete the order")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while checking user details")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}
	isAdmin := gconv.String(user.Map()["role"]) == consts.RoleAdmin

	order, err := dao.OthersServiceOrders.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while checking order details")
	}
	if order.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "Service order not found")
	}
	if !isAdmin && gconv.Int64(order.Map()["user_id"]) != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "You cannot delete orders of other users")
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

	_, err = dao.OthersServiceOrders.Ctx(ctx).TX(tx).Data(g.Map{
		"deleted_at": gtime.Now(),
	}).Where("id", req.Id).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the order")
	}

	notiData := do.Notifications{
		UserId:         userID,
		Type:           "service_order_deleted",
		Content:        fmt.Sprintf("Service order #%d has been deleted.", req.Id),
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

	return &entity.OthersServiceOrderDeleteRes{Message: "Service order deleted successfully"}, nil
}

func (s *sOthersServiceOrder) OthersServiceOrderPayment(ctx context.Context, req *entity.OthersServiceOrderPaymentReq) (*entity.OthersServiceOrderItem, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to process the payment")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while checking user details")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}

	order, err := dao.OthersServiceOrders.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while checking order details")
	}
	if order.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "Service order not found")
	}
	if gconv.Int64(order.Map()["user_id"]) != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "You cannot process payments for orders of other users")
	}
	if gconv.String(order.Map()["payment_status"]) == "paid" {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "This order has already been paid")
	}

	if gconv.Float64(user.Map()["wallet_balance"]) < gconv.Float64(order.Map()["price"]) {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Your wallet balance is insufficient for this payment")
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

	_, err = dao.OthersServiceOrders.Ctx(ctx).TX(tx).Data(g.Map{
		"payment_status": "paid",
		"updated_at":     gtime.Now(),
	}).Where("id", req.Id).Where("deleted_at IS NULL").Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the payment status")
	}

	_, err = dao.Users.Ctx(ctx).TX(tx).Data(g.Map{
		"wallet_balance": g.DB().Raw("wallet_balance - ?", order.Map()["price"]),
	}).Where("id", userID).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating your wallet balance")
	}

	txData := do.WalletTransactions{
		UserId:         userID,
		Amount:         -gconv.Float64(order.Map()["price"]),
		Type:           "debit",
		Description:    fmt.Sprintf("Payment for service order #%d", req.Id),
		RelatedOrderId: req.Id,
		CreatedAt:      gtime.Now(),
	}
	_, err = dao.WalletTransactions.Ctx(ctx).TX(tx).Data(txData).Insert()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while recording the transaction")
	}

	notiData := do.Notifications{
		UserId:         userID,
		Type:           "service_payment_confirmed",
		Content:        fmt.Sprintf("Payment of %.2f for service order #%d confirmed.", gconv.Float64(order.Map()["price"]), req.Id),
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

	var updatedOrder struct {
		entity.OthersServiceOrders
		LotName      string `json:"lot_name"`
		ServiceName  string `json:"service_name"`
		VehiclePlate string `json:"vehicle_plate"`
	}
	err = dao.OthersServiceOrders.Ctx(ctx).
		Fields("others_service_orders.*, parking_lots.name as lot_name, others_service.name as service_name, vehicles.license_plate as vehicle_plate").
		LeftJoin("parking_lots", "parking_lots.id = others_service_orders.lot_id").
		LeftJoin("others_service", "others_service.id = others_service_orders.service_id").
		LeftJoin("vehicles", "vehicles.id = others_service_orders.vehicle_id").
		Where("others_service_orders.id", req.Id).
		Scan(&updatedOrder)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while retrieving the updated order")
	}

	item := entity.OthersServiceOrderItem{
		Id:            updatedOrder.Id,
		UserId:        updatedOrder.UserId,
		LotId:         updatedOrder.LotId,
		ServiceId:     updatedOrder.ServiceId,
		VehicleId:     updatedOrder.VehicleId,
		LotName:       updatedOrder.LotName,
		ServiceName:   updatedOrder.ServiceName,
		VehiclePlate:  updatedOrder.VehiclePlate,
		ScheduledTime: updatedOrder.ScheduledTime.Format("2006-01-02 15:04:05"),
		Status:        updatedOrder.Status,
		Price:         updatedOrder.Price,
		PaymentStatus: updatedOrder.PaymentStatus,
		CreatedAt:     updatedOrder.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     updatedOrder.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt:     updatedOrder.DeletedAt.Format("2006-01-02 15:04:05"),
	}

	return &item, nil
}
