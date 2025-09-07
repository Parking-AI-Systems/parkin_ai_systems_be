package parking_order

import (
	"context"
	"fmt"
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

type sParkingOrder struct{}

func Init() {
	service.RegisterParkingOrder(&sParkingOrder{})
}
func init() {
	Init()
}

func (s *sParkingOrder) ParkingOrderAddWithUser(ctx context.Context, req *entity.ParkingOrderAddReq) (*entity.ParkingOrderAddRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to book a parking slot.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}

	vehicle, err := dao.Vehicles.Ctx(ctx).Where("id", req.VehicleId).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to find the vehicle. Please try again.")
	}
	if vehicle.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeVehicleNotFound, "The vehicle could not be found.")
	}

	lot, err := dao.ParkingLots.Ctx(ctx).Where("id", req.LotId).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to find the parking lot. Please try again.")
	}
	if lot.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeParkingLotNotFound, "The parking lot could not be found.")
	}

	slot, err := dao.ParkingSlots.Ctx(ctx).Where("id", req.SlotId).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to find the parking slot. Please try again.")
	}
	if slot.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeParkingSlotNotFound, "The parking slot could not be found.")
	}
	isAvailable := gconv.Bool(slot.Map()["is_available"])
	if !isAvailable {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "The selected parking slot is not available.")
	}

	if req.StartTime == "" || req.EndTime == "" {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Please provide both start and end times.")
	}
	startTime, err := gtime.StrToTime(req.StartTime)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Please use a valid format for the start time.")
	}
	endTime, err := gtime.StrToTime(req.EndTime)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Please use a valid format for the end time.")
	}
	if startTime.After(endTime) {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "The start time must be earlier than the end time.")
	}
	if startTime.Before(gtime.Now()) {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "The start time must be in the future.")
	}

	count, err := dao.ParkingOrders.Ctx(ctx).
		Where("slot_id", req.SlotId).
		Where("status", "confirmed").
		Where("start_time < ?", endTime).
		Where("end_time > ?", startTime).
		Where("deleted_at IS NULL").
		Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to check for booking conflicts. Please try again.")
	}
	if count > 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "The slot is already booked for the selected time.")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while booking your slot. Please try again later.")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	data := do.ParkingOrders{
		UserId:        userID,
		LotId:         req.LotId,
		SlotId:        req.SlotId,
		VehicleId:     req.VehicleId,
		StartTime:     startTime,
		EndTime:       endTime,
		Status:        "confirmed",
		Price:         gconv.Float64(lot.Map()["price_per_hour"]) * endTime.Sub(startTime).Hours(),
		PaymentStatus: "pending",
		CreatedAt:     gtime.Now(),
	}
	lastId, err := dao.ParkingOrders.Ctx(ctx).TX(tx).Data(data).InsertAndGetId()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while booking your slot. Please try again later.")
	}

	_, err = dao.ParkingSlots.Ctx(ctx).TX(tx).Data(g.Map{"is_available": false}).Where("id", req.SlotId).Where("deleted_at IS NULL").Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while booking your slot. Please try again later.")
	}

	_, err = tx.Exec("UPDATE parking_lots SET available_slots = available_slots - 1, updated_at = ? WHERE id = ? AND deleted_at IS NULL", gtime.Now(), req.LotId)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while booking your slot. Please try again later.")
	}

	notiData := do.Notifications{
		UserId:         userID,
		Type:           "order_created",
		Content:        fmt.Sprintf("Parking order #%d has been created successfully.", lastId),
		RelatedOrderId: lastId,
		IsRead:         false,
		CreatedAt:      gtime.Now(),
	}
	_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while booking your slot. Please try again later.")
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while booking your slot. Please try again later.")
	}

	return &entity.ParkingOrderAddRes{Id: lastId}, nil
}

func (s *sParkingOrder) ParkingOrderList(ctx context.Context, req *entity.ParkingOrderListReq) (*entity.ParkingOrderListRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to view your bookings.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}
	isAdmin := gconv.String(user.Map()["role"]) == consts.RoleAdmin

	baseQuery := dao.ParkingOrders.Ctx(ctx).
		LeftJoin("parking_lots", "parking_lots.id = parking_orders.lot_id").
		LeftJoin("parking_slots", "parking_slots.id = parking_orders.slot_id").
		LeftJoin("vehicles", "vehicles.id = parking_orders.vehicle_id").
		Where("parking_orders.deleted_at IS NULL").
		Where("parking_lots.deleted_at IS NULL OR parking_lots.id IS NULL").
		Where("parking_slots.deleted_at IS NULL OR parking_slots.id IS NULL").
		Where("vehicles.deleted_at IS NULL OR vehicles.id IS NULL")

	if req.UserId != 0 {
		if !isAdmin && gconv.Int64(userID) != req.UserId {
			return nil, gerror.NewCode(consts.CodeUnauthorized, "You can only view your own bookings or must be an admin.")
		}
		baseQuery = baseQuery.Where("parking_orders.user_id", req.UserId)
	} else if !isAdmin {
		baseQuery = baseQuery.Where("parking_orders.user_id", userID)
	}
	if req.LotId != 0 {
		baseQuery = baseQuery.Where("parking_orders.lot_id", req.LotId)
	}
	if req.Status != "" {
		baseQuery = baseQuery.Where("parking_orders.status", req.Status)
	}

	total, err := baseQuery.Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load your bookings. Please try again later.")
	}

	m := dao.ParkingOrders.Ctx(ctx).
		Fields("parking_orders.*, parking_lots.name as lot_name, parking_slots.code as slot_code, vehicles.license_plate as vehicle_plate").
		LeftJoin("parking_lots", "parking_lots.id = parking_orders.lot_id").
		LeftJoin("parking_slots", "parking_slots.id = parking_orders.slot_id").
		LeftJoin("vehicles", "vehicles.id = parking_orders.vehicle_id").
		Where("parking_orders.deleted_at IS NULL").
		Where("parking_lots.deleted_at IS NULL OR parking_lots.id IS NULL").
		Where("parking_slots.deleted_at IS NULL OR parking_slots.id IS NULL").
		Where("vehicles.deleted_at IS NULL OR vehicles.id IS NULL")

	if req.UserId != 0 {
		if !isAdmin && gconv.Int64(userID) != req.UserId {
			return nil, gerror.NewCode(consts.CodeUnauthorized, "You can only view your own bookings or must be an admin.")
		}
		m = m.Where("parking_orders.user_id", req.UserId)
	} else if !isAdmin {
		m = m.Where("parking_orders.user_id", userID)
	}
	if req.LotId != 0 {
		m = m.Where("parking_orders.lot_id", req.LotId)
	}
	if req.Status != "" {
		m = m.Where("parking_orders.status", req.Status)
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	m = m.Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize)

	var orders []struct {
		entity.ParkingOrders
		LotName      string `json:"lot_name"`
		SlotCode     string `json:"slot_code"`
		VehiclePlate string `json:"vehicle_plate"`
	}
	err = m.Order("parking_orders.id DESC").Scan(&orders)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load your bookings. Please try again later.")
	}

	list := make([]entity.ParkingOrderItem, 0, len(orders))
	for _, order := range orders {
		item := entity.ParkingOrderItem{
			Id:            order.Id,
			UserId:        order.UserId,
			LotId:         order.LotId,
			SlotId:        order.SlotId,
			VehicleId:     order.VehicleId,
			LotName:       order.LotName,
			SlotCode:      order.SlotCode,
			VehiclePlate:  order.VehiclePlate,
			StartTime:     order.StartTime.Format("2006-01-02 15:04:05"),
			EndTime:       order.EndTime.Format("2006-01-02 15:04:05"),
			Status:        order.Status,
			Price:         order.Price,
			PaymentStatus: order.PaymentStatus,
			CreatedAt:     order.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if !order.UpdatedAt.IsZero() {
			item.UpdatedAt = order.UpdatedAt.Format("2006-01-02 15:04:05")
		}
		list = append(list, item)
	}

	return &entity.ParkingOrderListRes{
		List:  list,
		Total: total,
	}, nil
}

func (s *sParkingOrder) ParkingOrderGet(ctx context.Context, req *entity.ParkingOrderGetReq) (*entity.ParkingOrderItem, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to view the booking.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}
	isAdmin := gconv.String(user.Map()["role"]) == consts.RoleAdmin

	var order struct {
		entity.ParkingOrders
		LotName      string `json:"lot_name"`
		SlotCode     string `json:"slot_code"`
		VehiclePlate string `json:"vehicle_plate"`
	}
	m := dao.ParkingOrders.Ctx(ctx).
		Fields("parking_orders.*, parking_lots.name as lot_name, parking_slots.code as slot_code, vehicles.license_plate as vehicle_plate").
		LeftJoin("parking_lots", "parking_lots.id = parking_orders.lot_id").
		LeftJoin("parking_slots", "parking_slots.id = parking_orders.slot_id").
		LeftJoin("vehicles", "vehicles.id = parking_orders.vehicle_id").
		Where("parking_orders.id", req.Id).
		Where("parking_orders.deleted_at IS NULL").
		Where("parking_lots.deleted_at IS NULL OR parking_lots.id IS NULL").
		Where("parking_slots.deleted_at IS NULL OR parking_slots.id IS NULL").
		Where("vehicles.deleted_at IS NULL OR vehicles.id IS NULL")
	if !isAdmin {
		m = m.Where("parking_orders.user_id", userID)
	}
	err = m.Scan(&order)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load the booking. Please try again later.")
	}
	if order.Id == 0 {
		return nil, gerror.NewCode(consts.CodeNotFound, "The booking could not be found.")
	}

	item := entity.ParkingOrderItem{
		Id:            order.Id,
		UserId:        order.UserId,
		LotId:         order.LotId,
		SlotId:        order.SlotId,
		VehicleId:     order.VehicleId,
		LotName:       order.LotName,
		SlotCode:      order.SlotCode,
		VehiclePlate:  order.VehiclePlate,
		StartTime:     order.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:       order.EndTime.Format("2006-01-02 15:04:05"),
		Status:        order.Status,
		Price:         order.Price,
		PaymentStatus: order.PaymentStatus,
		CreatedAt:     order.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	if !order.UpdatedAt.IsZero() {
		item.UpdatedAt = order.UpdatedAt.Format("2006-01-02 15:04:05")
	}

	return &item, nil
}

func (s *sParkingOrder) ParkingOrderUpdate(ctx context.Context, req *entity.ParkingOrderUpdateReq) (*entity.ParkingOrderItem, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to update the booking.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}

	isAdmin := gconv.String(user.Map()["role"]) == consts.RoleAdmin

	order, err := dao.ParkingOrders.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to find the booking. Please try again.")
	}
	if order.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "The booking could not be found.")
	}
	if !isAdmin && gconv.Int64(order.Map()["user_id"]) != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "You can only update your own bookings or must be an admin.")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the booking. Please try again later.")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	updateData := g.Map{
		"updated_at": gtime.Now(),
	}
	if req.StartTime != "" {
		startTime, err := gtime.StrToTime(req.StartTime)
		if err != nil {
			return nil, gerror.NewCode(consts.CodeInvalidInput, "Please use a valid format for the start time.")
		}
		if startTime.Before(gtime.Now()) {
			return nil, gerror.NewCode(consts.CodeInvalidInput, "The start time must be in the future.")
		}
		updateData["start_time"] = startTime
	}
	if req.EndTime != "" {
		endTime, err := gtime.StrToTime(req.EndTime)
		if err != nil {
			return nil, gerror.NewCode(consts.CodeInvalidInput, "Please use a valid format for the end time.")
		}
		if startTime, ok := updateData["start_time"].(*gtime.Time); ok && startTime.After(endTime) {
			return nil, gerror.NewCode(consts.CodeInvalidInput, "The start time must be earlier than the end time.")
		}
		updateData["end_time"] = endTime
	}
	if req.Status != "" {
		updateData["status"] = req.Status
	}

	if updateData["start_time"] != nil && updateData["end_time"] != nil {
		lot, err := dao.ParkingLots.Ctx(ctx).Where("id", order.Map()["lot_id"]).Where("deleted_at IS NULL").One()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to find the parking lot. Please try again.")
		}
		updateData["price"] = gconv.Float64(lot.Map()["price_per_hour"]) * updateData["end_time"].(*gtime.Time).Sub(updateData["start_time"].(*gtime.Time)).Hours()
	}

	_, err = dao.ParkingOrders.Ctx(ctx).TX(tx).Data(updateData).Where("id", req.Id).Where("deleted_at IS NULL").Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the booking. Please try again later.")
	}

	if req.Status != "" {
		notiData := do.Notifications{
			UserId:         gconv.String(order.Map()["user_id"]),
			Type:           "order_updated",
			Content:        fmt.Sprintf("Parking order #%d status updated to %s.", req.Id, req.Status),
			RelatedOrderId: req.Id,
			IsRead:         false,
			CreatedAt:      gtime.Now(),
		}
		_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the booking. Please try again later.")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the booking. Please try again later.")
	}

	var updatedOrder struct {
		entity.ParkingOrders
		LotName      string `json:"lot_name"`
		SlotCode     string `json:"slot_code"`
		VehiclePlate string `json:"vehicle_plate"`
	}
	err = dao.ParkingOrders.Ctx(ctx).
		Fields("parking_orders.*, parking_lots.name as lot_name, parking_slots.code as slot_code, vehicles.license_plate as vehicle_plate").
		LeftJoin("parking_lots", "parking_lots.id = parking_orders.lot_id").
		LeftJoin("parking_slots", "parking_slots.id = parking_orders.slot_id").
		LeftJoin("vehicles", "vehicles.id = parking_orders.vehicle_id").
		Where("parking_orders.id", req.Id).
		Where("parking_orders.deleted_at IS NULL").
		Where("parking_lots.deleted_at IS NULL OR parking_lots.id IS NULL").
		Where("parking_slots.deleted_at IS NULL OR parking_slots.id IS NULL").
		Where("vehicles.deleted_at IS NULL OR vehicles.id IS NULL").
		Scan(&updatedOrder)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the booking. Please try again later.")
	}

	item := entity.ParkingOrderItem{
		Id:            updatedOrder.Id,
		UserId:        updatedOrder.UserId,
		LotId:         updatedOrder.LotId,
		SlotId:        updatedOrder.SlotId,
		VehicleId:     updatedOrder.VehicleId,
		LotName:       updatedOrder.LotName,
		SlotCode:      updatedOrder.SlotCode,
		VehiclePlate:  updatedOrder.VehiclePlate,
		StartTime:     updatedOrder.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:       updatedOrder.EndTime.Format("2006-01-02 15:04:05"),
		Status:        updatedOrder.Status,
		Price:         updatedOrder.Price,
		PaymentStatus: updatedOrder.PaymentStatus,
		CreatedAt:     updatedOrder.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	if !updatedOrder.UpdatedAt.IsZero() {
		item.UpdatedAt = updatedOrder.UpdatedAt.Format("2006-01-02 15:04:05")
	}

	return &item, nil
}

func (s *sParkingOrder) ParkingOrderCancel(ctx context.Context, req *entity.ParkingOrderCancelReq) (*entity.ParkingOrderItem, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to cancel the booking.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}

	isAdmin := gconv.String(user.Map()["role"]) == consts.RoleAdmin

	order, err := dao.ParkingOrders.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to find the booking. Please try again.")
	}
	if order.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "The booking could not be found.")
	}
	if !isAdmin && gconv.Int64(order.Map()["user_id"]) != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "You can only cancel your own bookings or must be an admin.")
	}
	if gconv.String(order.Map()["status"]) == "canceled" {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "This booking has already been canceled.")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while canceling the booking. Please try again later.")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = dao.ParkingOrders.Ctx(ctx).TX(tx).Data(g.Map{
		"status":     "canceled",
		"updated_at": gtime.Now(),
	}).Where("id", req.Id).Where("deleted_at IS NULL").Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while canceling the booking. Please try again later.")
	}

	_, err = dao.ParkingSlots.Ctx(ctx).TX(tx).Data(g.Map{
		"is_available": true,
	}).Where("id", order.Map()["slot_id"]).Where("deleted_at IS NULL").Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while canceling the booking. Please try again later.")
	}

	_, err = tx.Exec("UPDATE parking_lots SET available_slots = available_slots + 1, updated_at = ? WHERE id = ? AND deleted_at IS NULL", gtime.Now(), order.Map()["lot_id"])
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while canceling the booking. Please try again later.")
	}

	notiData := do.Notifications{
		UserId:         userID,
		Type:           "order_canceled",
		Content:        fmt.Sprintf("Parking order #%d has been canceled.", req.Id),
		RelatedOrderId: req.Id,
		IsRead:         false,
		CreatedAt:      gtime.Now(),
	}
	_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while canceling the booking. Please try again later.")
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while canceling the booking. Please try again later.")
	}

	var updatedOrder struct {
		entity.ParkingOrders
		LotName      string `json:"lot_name"`
		SlotCode     string `json:"slot_code"`
		VehiclePlate string `json:"vehicle_plate"`
	}
	err = dao.ParkingOrders.Ctx(ctx).
		Fields("parking_orders.*, parking_lots.name as lot_name, parking_slots.code as slot_code, vehicles.license_plate as vehicle_plate").
		LeftJoin("parking_lots", "parking_lots.id = parking_orders.lot_id").
		LeftJoin("parking_slots", "parking_slots.id = parking_orders.slot_id").
		LeftJoin("vehicles", "vehicles.id = parking_orders.vehicle_id").
		Where("parking_orders.id", req.Id).
		Where("parking_orders.deleted_at IS NULL").
		Where("parking_lots.deleted_at IS NULL OR parking_lots.id IS NULL").
		Where("parking_slots.deleted_at IS NULL OR parking_slots.id IS NULL").
		Where("vehicles.deleted_at IS NULL OR vehicles.id IS NULL").
		Scan(&updatedOrder)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while canceling the booking. Please try again later.")
	}

	item := entity.ParkingOrderItem{
		Id:            updatedOrder.Id,
		UserId:        updatedOrder.UserId,
		LotId:         updatedOrder.LotId,
		SlotId:        updatedOrder.SlotId,
		VehicleId:     updatedOrder.VehicleId,
		LotName:       updatedOrder.LotName,
		SlotCode:      updatedOrder.SlotCode,
		VehiclePlate:  updatedOrder.VehiclePlate,
		StartTime:     updatedOrder.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:       updatedOrder.EndTime.Format("2006-01-02 15:04:05"),
		Status:        updatedOrder.Status,
		Price:         updatedOrder.Price,
		PaymentStatus: updatedOrder.PaymentStatus,
		CreatedAt:     updatedOrder.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	if !updatedOrder.UpdatedAt.IsZero() {
		item.UpdatedAt = updatedOrder.UpdatedAt.Format("2006-01-02 15:04:05")
	}

	return &item, nil
}

func (s *sParkingOrder) ParkingOrderDelete(ctx context.Context, req *entity.ParkingOrderDeleteReq) (*entity.ParkingOrderDeleteRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to delete the booking.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}
	isAdmin := gconv.String(user.Map()["role"]) == consts.RoleAdmin

	order, err := dao.ParkingOrders.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to find the booking. Please try again.")
	}
	if order.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "The booking could not be found.")
	}
	if !isAdmin && gconv.Int64(order.Map()["user_id"]) != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "You can only delete your own bookings or must be an admin.")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the booking. Please try again later.")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = dao.ParkingOrders.Ctx(ctx).TX(tx).Data(g.Map{
		"deleted_at": gtime.Now(),
		"updated_at": gtime.Now(),
	}).Where("id", req.Id).Where("deleted_at IS NULL").Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the booking. Please try again later.")
	}

	_, err = dao.ParkingSlots.Ctx(ctx).TX(tx).Data(g.Map{
		"is_available": true,
	}).Where("id", order.Map()["slot_id"]).Where("deleted_at IS NULL").Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the booking. Please try again later.")
	}

	_, err = tx.Exec("UPDATE parking_lots SET available_slots = available_slots + 1, updated_at = ? WHERE id = ? AND deleted_at IS NULL", gtime.Now(), order.Map()["lot_id"])
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the booking. Please try again later.")
	}

	notiData := do.Notifications{
		UserId:         userID,
		Type:           "order_deleted",
		Content:        fmt.Sprintf("Parking order #%d has been deleted.", req.Id),
		RelatedOrderId: req.Id,
		IsRead:         false,
		CreatedAt:      gtime.Now(),
	}
	_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the booking. Please try again later.")
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the booking. Please try again later.")
	}

	return &entity.ParkingOrderDeleteRes{Message: "Parking order deleted successfully"}, nil
}

func (s *sParkingOrder) ParkingOrderPayment(ctx context.Context, req *entity.ParkingOrderPaymentReq) (*entity.ParkingOrderItem, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to process the payment.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}

	order, err := dao.ParkingOrders.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to find the booking. Please try again.")
	}
	if order.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "The booking could not be found.")
	}
	if gconv.Int64(order.Map()["user_id"]) != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "You can only process payments for your own bookings.")
	}
	if gconv.String(order.Map()["payment_status"]) == "paid" {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "This booking has already been paid.")
	}

	if gconv.Float64(user.Map()["wallet_balance"]) < gconv.Float64(order.Map()["price"]) {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Your wallet balance is not sufficient for this payment.")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while processing the payment. Please try again later.")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = dao.ParkingOrders.Ctx(ctx).TX(tx).Data(g.Map{
		"payment_status": "paid",
		"updated_at":     gtime.Now(),
	}).Where("id", req.Id).Where("deleted_at IS NULL").Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while processing the payment. Please try again later.")
	}

	_, err = tx.Exec("UPDATE users SET wallet_balance = wallet_balance - ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL", order.Map()["price"], gtime.Now(), userID)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while processing the payment. Please try again later.")
	}

	txData := do.WalletTransactions{
		UserId:         userID,
		Amount:         -gconv.Float64(order.Map()["price"]),
		Type:           "debit",
		Description:    fmt.Sprintf("Payment for parking order #%d", req.Id),
		RelatedOrderId: req.Id,
		CreatedAt:      gtime.Now(),
	}
	_, err = dao.WalletTransactions.Ctx(ctx).TX(tx).Data(txData).Insert()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while processing the payment. Please try again later.")
	}

	notiData := do.Notifications{
		UserId:         userID,
		Type:           "payment_confirmed",
		Content:        fmt.Sprintf("Payment of %.2f for parking order #%d confirmed.", gconv.Float64(order.Map()["price"]), req.Id),
		RelatedOrderId: req.Id,
		IsRead:         false,
		CreatedAt:      gtime.Now(),
	}
	_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while processing the payment. Please try again later.")
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while processing the payment. Please try again later.")
	}

	var updatedOrder struct {
		entity.ParkingOrders
		LotName      string `json:"lot_name"`
		SlotCode     string `json:"slot_code"`
		VehiclePlate string `json:"vehicle_plate"`
	}
	err = dao.ParkingOrders.Ctx(ctx).
		Fields("parking_orders.*, parking_lots.name as lot_name, parking_slots.code as slot_code, vehicles.license_plate as vehicle_plate").
		LeftJoin("parking_lots", "parking_lots.id = parking_orders.lot_id").
		LeftJoin("parking_slots", "parking_slots.id = parking_orders.slot_id").
		LeftJoin("vehicles", "vehicles.id = parking_orders.vehicle_id").
		Where("parking_orders.id", req.Id).
		Where("parking_orders.deleted_at IS NULL").
		Where("parking_lots.deleted_at IS NULL OR parking_lots.id IS NULL").
		Where("parking_slots.deleted_at IS NULL OR parking_slots.id IS NULL").
		Where("vehicles.deleted_at IS NULL OR vehicles.id IS NULL").
		Scan(&updatedOrder)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while processing the payment. Please try again later.")
	}

	item := entity.ParkingOrderItem{
		Id:            updatedOrder.Id,
		UserId:        updatedOrder.UserId,
		LotId:         updatedOrder.LotId,
		SlotId:        updatedOrder.SlotId,
		VehicleId:     updatedOrder.VehicleId,
		LotName:       updatedOrder.LotName,
		SlotCode:      updatedOrder.SlotCode,
		VehiclePlate:  updatedOrder.VehiclePlate,
		StartTime:     updatedOrder.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:       updatedOrder.EndTime.Format("2006-01-02 15:04:05"),
		Status:        updatedOrder.Status,
		Price:         updatedOrder.Price,
		PaymentStatus: updatedOrder.PaymentStatus,
		CreatedAt:     updatedOrder.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	if !updatedOrder.UpdatedAt.IsZero() {
		item.UpdatedAt = updatedOrder.UpdatedAt.Format("2006-01-02 15:04:05")
	}

	return &item, nil
}

// New dashboard API implementations

func (s *sParkingOrder) ParkingOrderRevenue(ctx context.Context, req *entity.ParkingOrderRevenueReq) (*entity.ParkingOrderRevenueRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to access revenue data.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found.")
	}
	if gconv.String(user.Map()["role"]) != consts.RoleAdmin {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Only admins can access this data.")
	}

	m := dao.ParkingOrders.Ctx(ctx).Where("deleted_at IS NULL").Where("payment_status", "PAID")

	period := req.Period
	if period == "" {
		period = "1m"
	}

	var start *gtime.Time
	var end *gtime.Time = gtime.Now()

	if period == "custom" {
		start = gtime.NewFromStr(req.StartTime)
		if start == nil {
			return nil, gerror.NewCode(consts.CodeValidationFailed, "Invalid start time format.")
		}
		end = gtime.NewFromStr(req.EndTime)
		if end == nil {
			return nil, gerror.NewCode(consts.CodeValidationFailed, "Invalid end time format.")
		}
		if start.After(end) {
			return nil, gerror.NewCode(consts.CodeValidationFailed, "Start time must be before end time.")
		}
	} else {
		start = s.getStartTime(period)
	}

	m = m.WhereBetween("created_at", start, end)

	total, err := m.Sum("price")
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error calculating revenue.")
	}

	return &entity.ParkingOrderRevenueRes{TotalRevenue: gconv.Float64(total)}, nil
}

func (s *sParkingOrder) ParkingOrderTrends(ctx context.Context, req *entity.ParkingOrderTrendsReq) (*entity.ParkingOrderTrendsRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to access trends data.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found.")
	}
	if gconv.String(user.Map()["role"]) != consts.RoleAdmin {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Only admins can access this data.")
	}

	period := req.Period
	if period == "" {
		period = "1m"
	}

	var start *gtime.Time
	var end *gtime.Time = gtime.Now()

	if period == "custom" {
		start = gtime.NewFromStr(req.StartTime)
		if start == nil {
			return nil, gerror.NewCode(consts.CodeValidationFailed, "Invalid start time format.")
		}
		end = gtime.NewFromStr(req.EndTime)
		if end == nil {
			return nil, gerror.NewCode(consts.CodeValidationFailed, "Invalid end time format.")
		}
		if start.After(end) {
			return nil, gerror.NewCode(consts.CodeValidationFailed, "Start time must be before end time.")
		}
	} else {
		start = s.getStartTime(period)
	}

	groupField, dateFormat, step := s.getTrendsConfig(period, start, end)

	m := dao.ParkingOrders.Ctx(ctx).
		Fields(groupField+" as date, COUNT(*) as count").
		Where("deleted_at IS NULL").
		WhereBetween("created_at", start, end).
		Group("date").Order("date ASC")

	var dataList []struct {
		Date  *gtime.Time `json:"date"`
		Count int64       `json:"count"`
	}
	err = m.Scan(&dataList)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving trends data.")
	}

	dataMap := make(map[string]int64)
	for _, item := range dataList {
		key := item.Date.Format(dateFormat)
		dataMap[key] = item.Count
	}

	var orders []entity.ParkingOrderTrendsItem
	total := int64(0)
	current := start.Clone()
	for !current.After(end) {
		key := current.Format(dateFormat)
		count := dataMap[key]
		orders = append(orders, entity.ParkingOrderTrendsItem{Date: key, Count: count})
		total += count
		current = current.Add(step)
	}

	return &entity.ParkingOrderTrendsRes{Orders: orders, Total: total}, nil
}

func (s *sParkingOrder) ParkingOrderStatusBreakdown(ctx context.Context, req *entity.ParkingOrderStatusBreakdownReq) (*entity.ParkingOrderStatusBreakdownRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to access status breakdown data.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found.")
	}
	if gconv.String(user.Map()["role"]) != consts.RoleAdmin {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Only admins can access this data.")
	}

	m := dao.ParkingOrders.Ctx(ctx).Fields("status, COUNT(*) as count").Where("deleted_at IS NULL")

	period := req.Period
	if period == "" {
		period = "1m"
	}

	var start *gtime.Time
	var end *gtime.Time = gtime.Now()

	if period == "custom" {
		start = gtime.NewFromStr(req.StartTime)
		if start == nil {
			return nil, gerror.NewCode(consts.CodeValidationFailed, "Invalid start time format.")
		}
		end = gtime.NewFromStr(req.EndTime)
		if end == nil {
			return nil, gerror.NewCode(consts.CodeValidationFailed, "Invalid end time format.")
		}
		if start.After(end) {
			return nil, gerror.NewCode(consts.CodeValidationFailed, "Start time must be before end time.")
		}
	} else {
		start = s.getStartTime(period)
	}

	m = m.WhereBetween("created_at", start, end)

	var list []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	err = m.Group("status").Scan(&list)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving status breakdown.")
	}

	total := int64(0)
	statuses := []entity.ParkingOrderStatusItem{
		{Status: "pending", Count: 0},
		{Status: "confirmed", Count: 0},
		{Status: "canceled", Count: 0},
		{Status: "completed", Count: 0},
	}
	for _, item := range list {
		for i, status := range statuses {
			if status.Status == item.Status {
				statuses[i].Count = item.Count
				total += item.Count
			}
		}
	}

	return &entity.ParkingOrderStatusBreakdownRes{Statuses: statuses, Total: total}, nil
}

func (s *sParkingOrder) getStartTime(period string) *gtime.Time {
	now := gtime.Now()
	switch period {
	case "1h":
		return now.Add(-time.Hour)
	case "1d":
		return now.Add(-24 * time.Hour)
	case "1w":
		return now.Add(-7 * 24 * time.Hour)
	case "1m":
		return now.AddDate(0, -1, 0)
	default:
		return now.AddDate(0, -1, 0)
	}
}

func (s *sParkingOrder) getTrendsConfig(period string, start, end *gtime.Time) (groupField, dateFormat string, step time.Duration) {
	if period != "custom" {
		switch period {
		case "1h":
			return "DATE_TRUNC('minute', created_at)", "Y-m-d H:i", time.Minute
		case "1d":
			return "DATE_TRUNC('hour', created_at)", "Y-m-d H", time.Hour
		case "1w", "1m":
			return "DATE(created_at)", "Y-m-d", 24 * time.Hour
		default:
			return "DATE(created_at)", "Y-m-d", 24 * time.Hour
		}
	} else {
		diff := end.Sub(start)
		if diff <= time.Hour {
			return "DATE_TRUNC('minute', created_at)", "Y-m-d H:i", time.Minute
		} else if diff <= 24*time.Hour {
			return "DATE_TRUNC('hour', created_at)", "Y-m-d H", time.Hour
		} else {
			return "DATE(created_at)", "Y-m-d", 24 * time.Hour
		}
	}
}
