package vehicle

import (
	"context"
	"fmt"
	"parkin-ai-system/internal/consts"
	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/model/do"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
	"regexp"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type sVehicle struct{}

func Init() {
	service.RegisterVehicle(&sVehicle{})
}
func init() {
	Init()
}

// License plate format: e.g., 29A-12345 or 51H-67890
var licensePlateRegex = regexp.MustCompile(`^[0-9]{2}[A-Z]-[0-9]{4,5}$`)

func (s *sVehicle) VehicleAdd(ctx context.Context, req *entity.VehicleAddReq) (*entity.VehicleAddRes, error) {
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

	if !licensePlateRegex.MatchString(req.LicensePlate) {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Invalid license plate format. Must be like XXA-12345")
	}

	count, err := dao.Vehicles.Ctx(ctx).Where("license_plate", req.LicensePlate).Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking license plate uniqueness")
	}
	if count > 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "License plate already exists")
	}

	isValidVehicleType := false
	for _, validType := range consts.ValidVehicleTypes {
		if req.Type == validType {
			isValidVehicleType = true
			break
		}
	}
	if !isValidVehicleType {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Invalid vehicle type")
	}

	if len(req.LicensePlate) > 20 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "License plate must be less than 20 characters")
	}
	if req.Brand != "" && len(req.Brand) > 50 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Brand must be less than 50 characters")
	}
	if req.Model != "" && len(req.Model) > 50 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Model must be less than 50 characters")
	}
	if req.Color != "" && len(req.Color) > 50 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Color must be less than 50 characters")
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

	data := do.Vehicles{
		UserId:       gconv.Int64(userID),
		LicensePlate: req.LicensePlate,
		Brand:        req.Brand,
		Model:        req.Model,
		Color:        req.Color,
		Type:         req.Type,
		CreatedAt:    gtime.Now(),
	}
	lastId, err := dao.Vehicles.Ctx(ctx).TX(tx).Data(data).InsertAndGetId()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating vehicle")
	}

	adminUsers, err := dao.Users.Ctx(ctx).Where("role", "admin").All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving admins")
	}

	for _, admin := range adminUsers {
		notiData := do.Notifications{
			UserId:         admin["id"].Int64(),
			Type:           "vehicle_added",
			Content:        fmt.Sprintf("New vehicle #%d (%s) added by user #%d.", lastId, req.LicensePlate, gconv.Int64(userID)),
			RelatedOrderId: lastId,
			IsRead:         false,
			CreatedAt:      gtime.Now(),
		}
		_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating notification")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error committing transaction")
	}

	return &entity.VehicleAddRes{Id: lastId}, nil
}

func (s *sVehicle) VehicleList(ctx context.Context, req *entity.VehicleListReq) (*entity.VehicleListRes, error) {
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

	m := dao.Vehicles.Ctx(ctx).
		Fields("vehicles.*, users.username as username").
		LeftJoin("users", "users.id = vehicles.user_id")

	isAdmin := gconv.String(user.Map()["role"]) == "admin"
	if !isAdmin {
		m = m.Where("vehicles.user_id", userID)
	}

	if req.Type != "" {
		m = m.Where("vehicles.type", req.Type)
	}

	total, err := m.Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error counting vehicles")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	m = m.Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize)

	var vehicles []struct {
		entity.Vehicles
		Username string `json:"username"`
	}
	err = m.Order("vehicles.id DESC").Scan(&vehicles)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving vehicles")
	}

	list := make([]entity.VehicleItem, 0, len(vehicles))
	for _, v := range vehicles {
		item := entity.VehicleItem{
			Id:           v.Id,
			UserId:       v.UserId,
			Username:     v.Username,
			LicensePlate: v.LicensePlate,
			Brand:        v.Brand,
			Model:        v.Model,
			Color:        v.Color,
			Type:         v.Type,
			CreatedAt:    v.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		list = append(list, item)
	}

	return &entity.VehicleListRes{
		List:  list,
		Total: total,
	}, nil
}

func (s *sVehicle) VehicleGet(ctx context.Context, req *entity.VehicleGetReq) (*entity.VehicleItem, error) {
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

	var vehicle struct {
		entity.Vehicles
		Username string `json:"username"`
	}
	err = dao.Vehicles.Ctx(ctx).
		Fields("vehicles.*, users.username as username").
		LeftJoin("users", "users.id = vehicles.user_id").
		Where("vehicles.id", req.Id).
		Scan(&vehicle)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving vehicle")
	}
	if vehicle.Id == 0 {
		return nil, gerror.NewCode(consts.CodeNotFound, "Vehicle not found")
	}

	isAdmin := gconv.String(user.Map()["role"]) == "admin"
	if !isAdmin && vehicle.UserId != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "You can only view your own vehicles or must be an admin")
	}

	item := entity.VehicleItem{
		Id:           vehicle.Id,
		UserId:       vehicle.UserId,
		Username:     vehicle.Username,
		LicensePlate: vehicle.LicensePlate,
		Brand:        vehicle.Brand,
		Model:        vehicle.Model,
		Color:        vehicle.Color,
		Type:         vehicle.Type,
		CreatedAt:    vehicle.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return &item, nil
}

func (s *sVehicle) VehicleUpdate(ctx context.Context, req *entity.VehicleUpdateReq) (*entity.VehicleItem, error) {
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

	vehicle, err := dao.Vehicles.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking vehicle")
	}
	if vehicle.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "Vehicle not found")
	}

	isAdmin := gconv.String(user.Map()["role"]) == "admin"
	if !isAdmin && gconv.Int64(vehicle.Map()["user_id"]) != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "You can only update your own vehicles or must be an admin")
	}

	if req.LicensePlate != "" {
		if !licensePlateRegex.MatchString(req.LicensePlate) {
			return nil, gerror.NewCode(consts.CodeInvalidInput, "Invalid license plate format. Must be like XXA-12345")
		}
		count, err := dao.Vehicles.Ctx(ctx).
			Where("license_plate", req.LicensePlate).
			Where("id != ?", req.Id).
			Count()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking license plate uniqueness")
		}
		if count > 0 {
			return nil, gerror.NewCode(consts.CodeInvalidInput, "License plate already exists")
		}
	}

	if req.Type != "" {
		isValidVehicleType := false
		for _, validType := range consts.ValidVehicleTypes {
			if req.Type == validType {
				isValidVehicleType = true
				break
			}
		}
		if !isValidVehicleType {
			return nil, gerror.NewCode(consts.CodeInvalidInput, "Invalid vehicle type")
		}
	}

	if req.LicensePlate != "" && len(req.LicensePlate) > 20 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "License plate must be less than 20 characters")
	}
	if req.Brand != "" && len(req.Brand) > 50 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Brand must be less than 50 characters")
	}
	if req.Model != "" && len(req.Model) > 50 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Model must be less than 50 characters")
	}
	if req.Color != "" && len(req.Color) > 50 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Color must be less than 50 characters")
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
	if req.LicensePlate != "" {
		updateData["license_plate"] = req.LicensePlate
	}
	if req.Brand != "" {
		updateData["brand"] = req.Brand
	}
	if req.Model != "" {
		updateData["model"] = req.Model
	}
	if req.Color != "" {
		updateData["color"] = req.Color
	}
	if req.Type != "" {
		updateData["type"] = req.Type
	}

	_, err = dao.Vehicles.Ctx(ctx).TX(tx).Data(updateData).Where("id", req.Id).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error updating vehicle")
	}

	adminUsers, err := dao.Users.Ctx(ctx).Where("role", "admin").All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving admins")
	}

	for _, admin := range adminUsers {
		notiData := do.Notifications{
			UserId:         admin["id"].Int64(),
			Type:           "vehicle_updated",
			Content:        fmt.Sprintf("Vehicle #%d (%s) has been updated.", req.Id, req.LicensePlate),
			RelatedOrderId: req.Id,
			IsRead:         false,
			CreatedAt:      gtime.Now(),
		}
		_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating notification")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error committing transaction")
	}

	var updatedVehicle struct {
		entity.Vehicles
		Username string `json:"username"`
	}
	err = dao.Vehicles.Ctx(ctx).
		Fields("vehicles.*, users.username as username").
		LeftJoin("users", "users.id = vehicles.user_id").
		Where("vehicles.id", req.Id).
		Scan(&updatedVehicle)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving updated vehicle")
	}

	item := entity.VehicleItem{
		Id:           updatedVehicle.Id,
		UserId:       updatedVehicle.UserId,
		Username:     updatedVehicle.Username,
		LicensePlate: updatedVehicle.LicensePlate,
		Brand:        updatedVehicle.Brand,
		Model:        updatedVehicle.Model,
		Color:        updatedVehicle.Color,
		Type:         updatedVehicle.Type,
		CreatedAt:    updatedVehicle.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return &item, nil
}

func (s *sVehicle) VehicleDelete(ctx context.Context, req *entity.VehicleDeleteReq) (*entity.VehicleDeleteRes, error) {
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

	vehicle, err := dao.Vehicles.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking vehicle")
	}
	if vehicle.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "Vehicle not found")
	}

	isAdmin := gconv.String(user.Map()["role"]) == "admin"
	if !isAdmin && gconv.Int64(vehicle.Map()["user_id"]) != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "You can only delete your own vehicles or must be an admin")
	}

	count, err := dao.ParkingOrders.Ctx(ctx).
		Where("vehicle_id", req.Id).
		Where("status NOT IN (?)", g.Slice{"completed", "cancelled"}).
		Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking active orders")
	}
	if count > 0 {
		return nil, gerror.NewCode(consts.CodeInvalidOperation, "Cannot delete vehicle with active orders")
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

	_, err = dao.Vehicles.Ctx(ctx).TX(tx).Where("id", req.Id).Delete()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error deleting vehicle")
	}

	adminUsers, err := dao.Users.Ctx(ctx).Where("role", "admin").All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving admins")
	}

	for _, admin := range adminUsers {
		notiData := do.Notifications{
			UserId:         admin["id"].Int64(),
			Type:           "vehicle_deleted",
			Content:        fmt.Sprintf("Vehicle #%d (%s) has been deleted.", req.Id, vehicle.Map()["license_plate"]),
			RelatedOrderId: req.Id,
			IsRead:         false,
			CreatedAt:      gtime.Now(),
		}
		_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating notification")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error committing transaction")
	}

	return &entity.VehicleDeleteRes{Message: "Vehicle deleted successfully"}, nil
}

// CheckVehicleSlotCompatibility checks if a vehicle type is compatible with a parking slot type
func (s *sVehicle) CheckVehicleSlotCompatibility(ctx context.Context, vehicleID, slotID int64) error {
	vehicle, err := dao.Vehicles.Ctx(ctx).Where("id", vehicleID).One()
	if err != nil {
		return gerror.NewCode(consts.CodeDatabaseError, "Error checking vehicle")
	}
	if vehicle.IsEmpty() {
		return gerror.NewCode(consts.CodeNotFound, "Vehicle not found")
	}

	slot, err := dao.ParkingSlots.Ctx(ctx).Where("id", slotID).One()
	if err != nil {
		return gerror.NewCode(consts.CodeDatabaseError, "Error checking parking slot")
	}
	if slot.IsEmpty() {
		return gerror.NewCode(consts.CodeNotFound, "Parking slot not found")
	}

	vehicleType := gconv.String(vehicle.Map()["type"])
	slotType := gconv.String(slot.Map()["slot_type"])

	compatibleSlots, exists := consts.VehicleSlotCompatibility[vehicleType]
	if !exists {
		return gerror.NewCode(consts.CodeInvalidInput, "Invalid vehicle type")
	}

	for _, compatibleSlotType := range compatibleSlots {
		if slotType == compatibleSlotType {
			return nil
		}
	}

	return gerror.NewCode(consts.CodeInvalidInput, fmt.Sprintf("Vehicle type %s is not compatible with slot type %s", vehicleType, slotType))
}
