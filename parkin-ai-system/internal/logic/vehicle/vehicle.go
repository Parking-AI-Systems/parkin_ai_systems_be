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
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to add a vehicle.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}

	if !licensePlateRegex.MatchString(req.LicensePlate) {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Invalid license plate format. Use a format like XXA-12345.")
	}

	count, err := dao.Vehicles.Ctx(ctx).Where("license_plate", req.LicensePlate).Where("deleted_at IS NULL").Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to check if the license plate is available. Please try again.")
	}
	if count > 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "This license plate is already registered.")
	}

	isValidVehicleType := false
	for _, validType := range consts.ValidVehicleTypes {
		if req.Type == validType {
			isValidVehicleType = true
			break
		}
	}
	if !isValidVehicleType {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Invalid vehicle type. Please choose a valid type (e.g., car, motorcycle).")
	}

	if len(req.LicensePlate) > 20 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "License plate must be 20 characters or fewer.")
	}
	if req.Brand != "" && len(req.Brand) > 50 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Brand name must be 50 characters or fewer.")
	}
	if req.Model != "" && len(req.Model) > 50 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Model name must be 50 characters or fewer.")
	}
	if req.Color != "" && len(req.Color) > 50 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Color name must be 50 characters or fewer.")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding the vehicle. Please try again later.")
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
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding the vehicle. Please try again later.")
	}

	adminUsers, err := dao.Users.Ctx(ctx).Where("role", consts.RoleAdmin).Where("deleted_at IS NULL").All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding the vehicle. Please try again later.")
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
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding the vehicle. Please try again later.")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding the vehicle. Please try again later.")
	}

	return &entity.VehicleAddRes{Id: lastId}, nil
}

func (s *sVehicle) VehicleList(ctx context.Context, req *entity.VehicleListReq) (*entity.VehicleListRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to view vehicles.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}

	m := dao.Vehicles.Ctx(ctx).
		LeftJoin("users", "users.id = vehicles.user_id").
		Where("vehicles.deleted_at IS NULL").
		Where("users.deleted_at IS NULL")

	isAdmin := gconv.String(user.Map()["role"]) == consts.RoleAdmin
	if !isAdmin {
		m = m.Where("vehicles.user_id", userID)
	}

	if req.Type != "" {
		m = m.Where("vehicles.type", req.Type)
	}

	total, err := m.Fields("vehicles.id").Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load vehicles. Please try again later.")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	m = m.Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize)

	records, err := m.Fields(
		"vehicles.id, vehicles.user_id, vehicles.license_plate, vehicles.brand, vehicles.model, vehicles.color, vehicles.type, vehicles.created_at, users.username",
	).Order("vehicles.id DESC").All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load vehicles. Please try again later.")
	}

	g.Log().Debug(ctx, "VehicleList - Records count:", len(records))
	if len(records) > 0 {
		g.Log().Debug(ctx, "VehicleList - First record:", records[0])
		g.Log().Debug(ctx, "VehicleList - user_id value:", records[0]["user_id"])
		g.Log().Debug(ctx, "VehicleList - username value:", records[0]["username"])
	}

	list := make([]entity.VehicleItem, 0, len(records))
	for _, record := range records {
		g.Log().Debug(ctx, "Processing record:", record)

		username := ""
		if record["username"] != nil && !record["username"].IsNil() {
			username = record["username"].String()
		}

		userId := int64(0)
		if record["user_id"] != nil && !record["user_id"].IsNil() {
			userId = record["user_id"].Int64()
		}

		g.Log().Debug(ctx, "Processed - userId:", userId, "username:", username)

		item := entity.VehicleItem{
			Id:           record["id"].Int64(),
			UserId:       userId,
			Username:     username,
			LicensePlate: record["license_plate"].String(),
			Brand:        record["brand"].String(),
			Model:        record["model"].String(),
			Color:        record["color"].String(),
			Type:         record["type"].String(),
			CreatedAt:    record["created_at"].Time().Format("2006-01-02 15:04:05"),
		}

		g.Log().Debug(ctx, "Created item:", item)
		list = append(list, item)
	}

	g.Log().Debug(ctx, "Final list before return:", list)
	return &entity.VehicleListRes{
		List:  list,
		Total: total,
	}, nil
}

func (s *sVehicle) VehicleGet(ctx context.Context, req *entity.VehicleGetReq) (*entity.VehicleItem, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to view vehicle details.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}

	var vehicle struct {
		entity.Vehicles
		Username string `json:"username"`
	}
	err = dao.Vehicles.Ctx(ctx).
		Fields("vehicles.*, users.username as username").
		LeftJoin("users", "users.id = vehicles.user_id").
		Where("vehicles.id", req.Id).
		Where("vehicles.deleted_at IS NULL").
		Where("users.deleted_at IS NULL").
		Scan(&vehicle)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load vehicle details. Please try again later.")
	}
	if vehicle.Id == 0 {
		return nil, gerror.NewCode(consts.CodeNotFound, "The vehicle could not be found.")
	}

	isAdmin := gconv.String(user.Map()["role"]) == consts.RoleAdmin
	if !isAdmin && vehicle.UserId != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "You can only view your own vehicles or need admin access.")
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
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to update vehicle details.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}

	vehicle, err := dao.Vehicles.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to find the vehicle. Please try again.")
	}
	if vehicle.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "The vehicle could not be found.")
	}

	isAdmin := gconv.String(user.Map()["role"]) == consts.RoleAdmin
	if !isAdmin && gconv.Int64(vehicle.Map()["user_id"]) != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "You can only update your own vehicles or need admin access.")
	}

	if req.LicensePlate != "" {
		if !licensePlateRegex.MatchString(req.LicensePlate) {
			return nil, gerror.NewCode(consts.CodeInvalidInput, "Invalid license plate format. Use a format like XXA-12345.")
		}
		count, err := dao.Vehicles.Ctx(ctx).
			Where("license_plate", req.LicensePlate).
			Where("id != ?", req.Id).
			Where("deleted_at IS NULL").
			Count()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to check if the license plate is available. Please try again.")
		}
		if count > 0 {
			return nil, gerror.NewCode(consts.CodeInvalidInput, "This license plate is already registered.")
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
			return nil, gerror.NewCode(consts.CodeInvalidInput, "Invalid vehicle type. Please choose a valid type (e.g., car, motorcycle).")
		}
	}

	if req.LicensePlate != "" && len(req.LicensePlate) > 20 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "License plate must be 20 characters or fewer.")
	}
	if req.Brand != "" && len(req.Brand) > 50 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Brand name must be 50 characters or fewer.")
	}
	if req.Model != "" && len(req.Model) > 50 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Model name must be 50 characters or fewer.")
	}
	if req.Color != "" && len(req.Color) > 50 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Color name must be 50 characters or fewer.")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the vehicle. Please try again later.")
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
	updateData["updated_at"] = gtime.Now()
	_, err = dao.Vehicles.Ctx(ctx).TX(tx).Data(updateData).Where("id", req.Id).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the vehicle. Please try again later.")
	}

	adminUsers, err := dao.Users.Ctx(ctx).Where("role", consts.RoleAdmin).Where("deleted_at IS NULL").All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the vehicle. Please try again later.")
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
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the vehicle. Please try again later.")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the vehicle. Please try again later.")
	}

	var updatedVehicle struct {
		entity.Vehicles
		Username string `json:"username"`
	}
	err = dao.Vehicles.Ctx(ctx).
		Fields("vehicles.*, users.username as username").
		LeftJoin("users", "users.id = vehicles.user_id").
		Where("vehicles.id", req.Id).
		Where("vehicles.deleted_at IS NULL").
		Where("users.deleted_at IS NULL").
		Scan(&updatedVehicle)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load updated vehicle details. Please try again later.")
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
		UpdatedAt:    updatedVehicle.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt:    updatedVehicle.DeletedAt.Format("2006-01-02 15:04:05"),
	}

	return &item, nil
}

func (s *sVehicle) VehicleDelete(ctx context.Context, req *entity.VehicleDeleteReq) (*entity.VehicleDeleteRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to delete a vehicle.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}

	vehicle, err := dao.Vehicles.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to find the vehicle. Please try again.")
	}
	if vehicle.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "The vehicle could not be found.")
	}

	isAdmin := gconv.String(user.Map()["role"]) == consts.RoleAdmin
	if !isAdmin && gconv.Int64(vehicle.Map()["user_id"]) != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "You can only delete your own vehicles or need admin access.")
	}

	count, err := dao.ParkingOrders.Ctx(ctx).
		Where("vehicle_id", req.Id).
		Where("status NOT IN (?)", g.Slice{"completed", "cancelled"}).
		Where("deleted_at IS NULL").
		Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to check vehicle orders. Please try again later.")
	}
	if count > 0 {
		return nil, gerror.NewCode(consts.CodeInvalidOperation, "Cannot delete this vehicle because it has active parking orders.")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the vehicle. Please try again later.")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = dao.Vehicles.Ctx(ctx).TX(tx).Data(g.Map{
		"deleted_at": gtime.Now(),
	}).Where("id", req.Id).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the vehicle. Please try again later.")
	}
	adminUsers, err := dao.Users.Ctx(ctx).Where("role", consts.RoleAdmin).Where("deleted_at IS NULL").All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the vehicle. Please try again later.")
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
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the vehicle. Please try again later.")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the vehicle. Please try again later.")
	}

	return &entity.VehicleDeleteRes{Message: "Vehicle deleted successfully"}, nil
}

func (s *sVehicle) CheckVehicleSlotCompatibility(ctx context.Context, vehicleID, slotID int64) error {
	vehicle, err := dao.Vehicles.Ctx(ctx).Where("id", vehicleID).Where("deleted_at IS NULL").One()
	if err != nil {
		return gerror.NewCode(consts.CodeDatabaseError, "Unable to verify the vehicle. Please try again.")
	}
	if vehicle.IsEmpty() {
		return gerror.NewCode(consts.CodeNotFound, "The vehicle could not be found.")
	}

	slot, err := dao.ParkingSlots.Ctx(ctx).Where("id", slotID).Where("deleted_at IS NULL").One()
	if err != nil {
		return gerror.NewCode(consts.CodeDatabaseError, "Unable to verify the parking slot. Please try again.")
	}
	if slot.IsEmpty() {
		return gerror.NewCode(consts.CodeNotFound, "The parking slot could not be found.")
	}

	vehicleType := gconv.String(vehicle.Map()["type"])
	slotType := gconv.String(slot.Map()["slot_type"])

	compatibleSlots, exists := consts.VehicleSlotCompatibility[vehicleType]
	if !exists {
		return gerror.NewCode(consts.CodeInvalidInput, "Invalid vehicle type. Please choose a valid type (e.g., car, motorcycle).")
	}

	for _, compatibleSlotType := range compatibleSlots {
		if slotType == compatibleSlotType {
			return nil
		}
	}

	return gerror.NewCode(consts.CodeInvalidInput, fmt.Sprintf("This vehicle type (%s) cannot be parked in this slot type (%s).", vehicleType, slotType))
}
