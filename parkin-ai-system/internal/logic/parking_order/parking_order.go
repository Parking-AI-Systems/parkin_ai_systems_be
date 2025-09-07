package parking_order

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
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}

	vehicle, err := dao.Vehicles.Ctx(ctx).Where("id", req.VehicleId).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking vehicle")
	}
	if vehicle.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeVehicleNotFound, "Vehicle not found")
	}

	lot, err := dao.ParkingLots.Ctx(ctx).Where("id", req.LotId).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking parking lot")
	}
	if lot.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeParkingLotNotFound, "Parking lot not found")
	}

	slot, err := dao.ParkingSlots.Ctx(ctx).Where("id", req.SlotId).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking parking slot")
	}
	if slot.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeParkingSlotNotFound, "Parking slot not found")
	}
	isAvailable := gconv.Bool(slot.Map()["is_available"])
	if !isAvailable {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Parking slot is not available")
	}

	if req.StartTime == "" || req.EndTime == "" {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Start time and end time are required")
	}
	startTime, err := gtime.StrToTime(req.StartTime)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Invalid start time format")
	}
	endTime, err := gtime.StrToTime(req.EndTime)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Invalid end time format")
	}
	if startTime.After(endTime) {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Start time must be before end time")
	}
	if startTime.Before(gtime.Now()) {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Start time must be in the future")
	}

	count, err := dao.ParkingOrders.Ctx(ctx).
		Where("slot_id", req.SlotId).
		Where("status", "confirmed").
		Where("start_time < ?", endTime).
		Where("end_time > ?", startTime).
		Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking overlapping orders")
	}
	if count > 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Slot is already booked for this time period")
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
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating parking order")
	}

	_, err = dao.ParkingSlots.Ctx(ctx).TX(tx).Data(g.Map{"is_available": false}).Where("id", req.SlotId).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error updating parking slot")
	}

	_, err = tx.Exec("UPDATE parking_lots SET available_slots = available_slots - 1, updated_at = ? WHERE id = ?", gtime.Now(), req.LotId)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error updating parking lot")
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
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating notification")
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error committing transaction")
	}

	return &entity.ParkingOrderAddRes{Id: lastId}, nil
}
func (s *sParkingOrder) ParkingOrderList(ctx context.Context, req *entity.ParkingOrderListReq) (*entity.ParkingOrderListRes, error) {
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
	isAdmin := gconv.String(user.Map()["role"]) == "admin"

	// Build base query for joins and where conditions
	baseQuery := dao.ParkingOrders.Ctx(ctx).
		LeftJoin("parking_lots", "parking_lots.id = parking_orders.lot_id").
		LeftJoin("parking_slots", "parking_slots.id = parking_orders.slot_id").
		LeftJoin("vehicles", "vehicles.id = parking_orders.vehicle_id")

	// Apply filters
	if req.UserId != 0 {
		if !isAdmin && gconv.Int64(userID) != req.UserId {
			return nil, gerror.NewCode(consts.CodeUnauthorized, "Cannot access orders of other users")
		}
		baseQuery = baseQuery.Where("parking_orders.user_id", req.UserId)
	} else if !isAdmin {
		// Non-admin users can only see their own orders
		baseQuery = baseQuery.Where("parking_orders.user_id", userID)
	}
	if req.LotId != 0 {
		baseQuery = baseQuery.Where("parking_orders.lot_id", req.LotId)
	}
	if req.Status != "" {
		baseQuery = baseQuery.Where("parking_orders.status", req.Status)
	}

	// Get total count for pagination
	total, err := baseQuery.Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error counting orders")
	}

	// Build data query with fields
	m := dao.ParkingOrders.Ctx(ctx).
		Fields("parking_orders.*, parking_lots.name as lot_name, parking_slots.code as slot_code, vehicles.license_plate as vehicle_plate").
		LeftJoin("parking_lots", "parking_lots.id = parking_orders.lot_id").
		LeftJoin("parking_slots", "parking_slots.id = parking_orders.slot_id").
		LeftJoin("vehicles", "vehicles.id = parking_orders.vehicle_id")

	// Apply same filters to data query
	if req.UserId != 0 {
		if !isAdmin && gconv.Int64(userID) != req.UserId {
			return nil, gerror.NewCode(consts.CodeUnauthorized, "Cannot access orders of other users")
		}
		m = m.Where("parking_orders.user_id", req.UserId)
	} else if !isAdmin {
		// Non-admin users can only see their own orders
		m = m.Where("parking_orders.user_id", userID)
	}
	if req.LotId != 0 {
		m = m.Where("parking_orders.lot_id", req.LotId)
	}
	if req.Status != "" {
		m = m.Where("parking_orders.status", req.Status)
	}

	// Apply pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	m = m.Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize)

	// Execute query
	var orders []struct {
		entity.ParkingOrders
		LotName      string `json:"lot_name"`
		SlotCode     string `json:"slot_code"`
		VehiclePlate string `json:"vehicle_plate"`
	}
	err = m.Order("parking_orders.id DESC").Scan(&orders)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving orders")
	}

	// Convert to response format
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
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	// Check if user exists and has admin role
	user, err := dao.Users.Ctx(ctx).Where("id", userID).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}
	isAdmin := gconv.String(user.Map()["role"]) == "admin"

	// Build query with joins
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
		Where("parking_orders.id", req.Id)
	if !isAdmin {
		m = m.Where("parking_orders.user_id", userID)
	}
	err = m.Scan(&order)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving order")
	}
	if order.Id == 0 {
		return nil, gerror.NewCode(consts.CodeNotFound, "Parking order not found")
	}

	// Convert to response format
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
	// Get user_id from context
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	// Check user and role
	user, err := dao.Users.Ctx(ctx).Where("id", userID).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}

	isAdmin := gconv.String(user.Map()["role"]) == "admin"

	// Check if order exists and user has access
	order, err := dao.ParkingOrders.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking order")
	}
	if order.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "Parking order not found")
	}
	if !isAdmin && gconv.Int64(order.Map()["user_id"]) != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Cannot update orders of other users")
	}

	// Start transaction
	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error starting transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Prepare update data
	updateData := g.Map{
		"updated_at": gtime.Now(),
	}
	if req.StartTime != "" {
		startTime, err := gtime.StrToTime(req.StartTime)
		if err != nil {
			return nil, gerror.NewCode(consts.CodeInvalidInput, "Invalid start time format")
		}
		if startTime.Before(gtime.Now()) {
			return nil, gerror.NewCode(consts.CodeInvalidInput, "Start time must be in the future")
		}
		updateData["start_time"] = startTime
	}
	if req.EndTime != "" {
		endTime, err := gtime.StrToTime(req.EndTime)
		if err != nil {
			return nil, gerror.NewCode(consts.CodeInvalidInput, "Invalid end time format")
		}
		if startTime, ok := updateData["start_time"].(*gtime.Time); ok && startTime.After(endTime) {
			return nil, gerror.NewCode(consts.CodeInvalidInput, "Start time must be before end time")
		}
		updateData["end_time"] = endTime
	}
	if req.Status != "" {
		updateData["status"] = req.Status
	}

	// Recalculate price if times change
	if updateData["start_time"] != nil && updateData["end_time"] != nil {
		lot, err := dao.ParkingLots.Ctx(ctx).Where("id", order.Map()["lot_id"]).One()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking parking lot")
		}
		updateData["price"] = gconv.Float64(lot.Map()["price_per_hour"]) * updateData["end_time"].(*gtime.Time).Sub(updateData["start_time"].(*gtime.Time)).Hours()
	}

	// Update order
	_, err = dao.ParkingOrders.Ctx(ctx).TX(tx).Data(updateData).Where("id", req.Id).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error updating order")
	}

	// Create notification
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
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating notification")
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error committing transaction")
	}

	// Retrieve updated order
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
		Scan(&updatedOrder)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving updated order")
	}

	// Convert to response format
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
	// Get user_id from context
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	// Check user and role
	user, err := dao.Users.Ctx(ctx).Where("id", userID).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}

	isAdmin := gconv.String(user.Map()["role"]) == "admin"

	// Check if order exists and user has access
	order, err := dao.ParkingOrders.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking order")
	}
	if order.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "Parking order not found")
	}
	if !isAdmin && gconv.Int64(order.Map()["user_id"]) != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Cannot cancel orders of other users")
	}
	if gconv.String(order.Map()["status"]) == "canceled" {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Order is already canceled")
	}

	// Start transaction
	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error starting transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Update order status
	_, err = dao.ParkingOrders.Ctx(ctx).TX(tx).Data(g.Map{
		"status":     "canceled",
		"updated_at": gtime.Now(),
	}).Where("id", req.Id).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error canceling order")
	}

	// Update slot availability
	_, err = dao.ParkingSlots.Ctx(ctx).TX(tx).Data(g.Map{
		"is_available": true,
	}).Where("id", order.Map()["slot_id"]).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error updating parking slot")
	}

	// Update lot available slots
	_, err = tx.Exec("UPDATE parking_lots SET available_slots = available_slots + 1, updated_at = ? WHERE id = ?", gtime.Now(), order.Map()["lot_id"])
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error updating parking lot")
	}

	// Create notification
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
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating notification")
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error committing transaction")
	}

	// Retrieve updated order
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
		Scan(&updatedOrder)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving updated order")
	}

	// Convert to response format
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
	// Get user_id from context
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	// Check user and role
	user, err := dao.Users.Ctx(ctx).Where("id", userID).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}
	isAdmin := gconv.String(user.Map()["role"]) == "admin"

	// Check if order exists and user has access
	order, err := dao.ParkingOrders.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking order")
	}
	if order.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "Parking order not found")
	}
	if !isAdmin && gconv.Int64(order.Map()["user_id"]) != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Cannot delete orders of other users")
	}

	// Start transaction
	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error starting transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Soft delete order
	_, err = dao.ParkingOrders.Ctx(ctx).TX(tx).Data(g.Map{
		"deleted_at": gtime.Now(),
	}).Where("id", req.Id).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error deleting order")
	}

	// Update slot availability
	_, err = dao.ParkingSlots.Ctx(ctx).TX(tx).Data(g.Map{
		"is_available": true,
	}).Where("id", order.Map()["slot_id"]).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error updating parking slot")
	}

	// Update lot available slots
	_, err = tx.Exec("UPDATE parking_lots SET available_slots = available_slots + 1, updated_at = ? WHERE id = ?", gtime.Now(), order.Map()["lot_id"])
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error updating parking lot")
	}

	// Create notification
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
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating notification")
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error committing transaction")
	}

	return &entity.ParkingOrderDeleteRes{Message: "Parking order deleted successfully"}, nil
}

func (s *sParkingOrder) ParkingOrderPayment(ctx context.Context, req *entity.ParkingOrderPaymentReq) (*entity.ParkingOrderItem, error) {
	// Get user_id from context
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	// Check user
	user, err := dao.Users.Ctx(ctx).Where("id", userID).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}

	// Check if order exists and user has access
	order, err := dao.ParkingOrders.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking order")
	}
	if order.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "Parking order not found")
	}
	if gconv.Int64(order.Map()["user_id"]) != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Cannot process payment for orders of other users")
	}
	if gconv.String(order.Map()["payment_status"]) == "paid" {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Order is already paid")
	}

	// Check wallet balance
	if gconv.Float64(user.Map()["wallet_balance"]) < gconv.Float64(order.Map()["price"]) {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Insufficient wallet balance")
	}

	// Start transaction
	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error starting transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Update payment status
	_, err = dao.ParkingOrders.Ctx(ctx).TX(tx).Data(g.Map{
		"payment_status": "paid",
		"updated_at":     gtime.Now(),
	}).Where("id", req.Id).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error updating payment status")
	}

	// Deduct wallet balance
	_, err = tx.Exec("UPDATE users SET wallet_balance = wallet_balance - ?, updated_at = ? WHERE id = ?", order.Map()["price"], gtime.Now(), userID)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error updating wallet balance")
	}

	// Create wallet transaction
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
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating wallet transaction")
	}

	// Create notification
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
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating notification")
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error committing transaction")
	}

	// Retrieve updated order
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
		Scan(&updatedOrder)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving updated order")
	}

	// Convert to response format
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
