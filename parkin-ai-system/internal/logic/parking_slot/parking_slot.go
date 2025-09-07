package parking_slot

import (
	"context"
	"fmt"

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

type sParkingSlot struct{}

func Init() {
	service.RegisterParkingSlot(&sParkingSlot{})
}
func init() {
	Init()
}

func (s *sParkingSlot) ParkingSlotAdd(ctx context.Context, req *entity.ParkingSlotAddReq) (*entity.ParkingSlotAddRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to add a parking slot.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}
	if gconv.String(user.Map()["role"]) != consts.RoleAdmin {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Only admins can add parking slots.")
	}

	lot, err := dao.ParkingLots.Ctx(ctx).Where("id", req.LotId).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to find the parking lot. Please try again.")
	}
	if lot.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeParkingLotNotFound, "The parking lot could not be found.")
	}

	count, err := dao.ParkingSlots.Ctx(ctx).Where("lot_id", req.LotId).Where("code", req.Code).Where("deleted_at IS NULL").Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify the slot code. Please try again.")
	}
	if count > 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "This slot code is already used in the parking lot.")
	}

	isValidSlotType := false
	for _, validType := range consts.ValidSlotTypes {
		if req.SlotType == validType {
			isValidSlotType = true
			break
		}
	}
	if !isValidSlotType {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Please select a valid slot type.")
	}

	if len(req.Code) > 20 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "The slot code must be 20 characters or fewer.")
	}
	if len(req.Floor) > 10 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "The floor name must be 10 characters or fewer.")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding the parking slot. Please try again later.")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	data := do.ParkingSlots{
		LotId:       req.LotId,
		Code:        req.Code,
		IsAvailable: req.IsAvailable,
		SlotType:    req.SlotType,
		Floor:       req.Floor,
		CreatedAt:   gtime.Now(),
	}
	lastId, err := dao.ParkingSlots.Ctx(ctx).TX(tx).Data(data).InsertAndGetId()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding the parking slot. Please try again later.")
	}

	adminUsers, err := dao.Users.Ctx(ctx).Where("role", consts.RoleAdmin).Where("deleted_at IS NULL").All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding the parking slot. Please try again later.")
	}

	for _, admin := range adminUsers {
		notiData := do.Notifications{
			UserId:         gconv.Int(admin.Map()["id"]),
			Type:           "parking_slot_added",
			Content:        fmt.Sprintf("New parking slot #%d (%s) added to parking lot #%d.", lastId, req.Code, req.LotId),
			RelatedOrderId: lastId,
			IsRead:         false,
			CreatedAt:      gtime.Now(),
		}
		_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding the parking slot. Please try again later.")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding the parking slot. Please try again later.")
	}

	return &entity.ParkingSlotAddRes{Id: lastId}, nil
}

func (s *sParkingSlot) ParkingSlotList(ctx context.Context, req *entity.ParkingSlotListReq) (*entity.ParkingSlotListRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to view parking slots.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}

	baseQuery := dao.ParkingSlots.Ctx(ctx).
		LeftJoin("parking_lots", "parking_lots.id = parking_slots.lot_id").
		Where("parking_slots.deleted_at IS NULL").
		Where("parking_lots.deleted_at IS NULL OR parking_lots.id IS NULL")

	if req.LotId != 0 {
		baseQuery = baseQuery.Where("parking_slots.lot_id", req.LotId)
	}
	if req.IsAvailable != nil {
		baseQuery = baseQuery.Where("parking_slots.is_available", *req.IsAvailable)
	}
	if req.SlotType != "" {
		baseQuery = baseQuery.Where("parking_slots.slot_type", req.SlotType)
	}

	total, err := baseQuery.Fields("parking_slots.id").Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load parking slots. Please try again later.")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	var slots []struct {
		entity.ParkingSlots
		LotName string `json:"lot_name"`
	}
	err = baseQuery.Fields("parking_slots.*, parking_lots.name as lot_name").
		Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize).
		Order("parking_slots.id DESC").Scan(&slots)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load parking slots. Please try again later.")
	}

	list := make([]entity.ParkingSlotItem, 0, len(slots))
	for _, s := range slots {
		item := entity.ParkingSlotItem{
			Id:          s.Id,
			LotId:       s.LotId,
			LotName:     s.LotName,
			Code:        s.Code,
			IsAvailable: s.IsAvailable,
			SlotType:    s.SlotType,
			Floor:       s.Floor,
			CreatedAt:   s.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		list = append(list, item)
	}

	return &entity.ParkingSlotListRes{
		List:  list,
		Total: total,
	}, nil
}

func (s *sParkingSlot) ParkingSlotGet(ctx context.Context, req *entity.ParkingSlotGetReq) (*entity.ParkingSlotItem, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to view the parking slot.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}

	var slot struct {
		entity.ParkingSlots
		LotName string `json:"lot_name"`
	}
	err = dao.ParkingSlots.Ctx(ctx).
		Fields("parking_slots.*, parking_lots.name as lot_name").
		LeftJoin("parking_lots", "parking_lots.id = parking_slots.lot_id").
		Where("parking_slots.id", req.Id).
		Where("parking_slots.deleted_at IS NULL").
		Where("parking_lots.deleted_at IS NULL OR parking_lots.id IS NULL").
		Scan(&slot)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load the parking slot. Please try again later.")
	}
	if slot.Id == 0 {
		return nil, gerror.NewCode(consts.CodeNotFound, "The parking slot could not be found.")
	}

	item := entity.ParkingSlotItem{
		Id:          slot.Id,
		LotId:       slot.LotId,
		LotName:     slot.LotName,
		Code:        slot.Code,
		IsAvailable: slot.IsAvailable,
		SlotType:    slot.SlotType,
		Floor:       slot.Floor,
		CreatedAt:   slot.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return &item, nil
}

func (s *sParkingSlot) ParkingSlotUpdate(ctx context.Context, req *entity.ParkingSlotUpdateReq) (*entity.ParkingSlotItem, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to update the parking slot.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}
	if gconv.String(user.Map()["role"]) != consts.RoleAdmin {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Only admins can update parking slots.")
	}

	slot, err := dao.ParkingSlots.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to find the parking slot. Please try again.")
	}
	if slot.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "The parking slot could not be found.")
	}

	if req.LotId != 0 {
		lot, err := dao.ParkingLots.Ctx(ctx).Where("id", req.LotId).Where("deleted_at IS NULL").One()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to find the parking lot. Please try again.")
		}
		if lot.IsEmpty() {
			return nil, gerror.NewCode(consts.CodeParkingLotNotFound, "The parking lot could not be found.")
		}
	}

	if req.Code != "" {
		count, err := dao.ParkingSlots.Ctx(ctx).
			Where("lot_id", req.LotId).
			Where("code", req.Code).
			Where("id != ?", req.Id).
			Where("deleted_at IS NULL").
			Count()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify the slot code. Please try again.")
		}
		if count > 0 {
			return nil, gerror.NewCode(consts.CodeInvalidInput, "This slot code is already used in the parking lot.")
		}
	}

	if req.SlotType != "" {
		isValidSlotType := false
		for _, validType := range consts.ValidSlotTypes {
			if req.SlotType == validType {
				isValidSlotType = true
				break
			}
		}
		if !isValidSlotType {
			return nil, gerror.NewCode(consts.CodeInvalidInput, "Please select a valid slot type.")
		}
	}

	if req.Code != "" && len(req.Code) > 20 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "The slot code must be 20 characters or fewer.")
	}
	if req.Floor != "" && len(req.Floor) > 10 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "The floor name must be 10 characters or fewer.")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the parking slot. Please try again later.")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	updateData := g.Map{
		"updated_at": gtime.Now(),
	}
	if req.LotId != 0 {
		updateData["lot_id"] = req.LotId
	}
	if req.Code != "" {
		updateData["code"] = req.Code
	}
	if req.IsAvailable != nil {
		updateData["is_available"] = *req.IsAvailable
	}
	if req.SlotType != "" {
		updateData["slot_type"] = req.SlotType
	}
	if req.Floor != "" {
		updateData["floor"] = req.Floor
	}

	_, err = dao.ParkingSlots.Ctx(ctx).TX(tx).Data(updateData).Where("id", req.Id).Where("deleted_at IS NULL").Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the parking slot. Please try again later.")
	}

	adminUsers, err := dao.Users.Ctx(ctx).Where("role", consts.RoleAdmin).Where("deleted_at IS NULL").All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the parking slot. Please try again later.")
	}

	for _, admin := range adminUsers {
		notiData := do.Notifications{
			UserId:         gconv.Int(admin.Map()["id"]),
			Type:           "parking_slot_updated",
			Content:        fmt.Sprintf("Parking slot #%d (%s) has been updated.", req.Id, req.Code),
			RelatedOrderId: req.Id,
			IsRead:         false,
			CreatedAt:      gtime.Now(),
		}
		_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the parking slot. Please try again later.")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the parking slot. Please try again later.")
	}

	var updatedSlot struct {
		entity.ParkingSlots
		LotName string `json:"lot_name"`
	}
	err = dao.ParkingSlots.Ctx(ctx).
		Fields("parking_slots.*, parking_lots.name as lot_name").
		LeftJoin("parking_lots", "parking_lots.id = parking_slots.lot_id").
		Where("parking_slots.id", req.Id).
		Where("parking_slots.deleted_at IS NULL").
		Where("parking_lots.deleted_at IS NULL OR parking_lots.id IS NULL").
		Scan(&updatedSlot)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the parking slot. Please try again later.")
	}

	item := entity.ParkingSlotItem{
		Id:          updatedSlot.Id,
		LotId:       updatedSlot.LotId,
		LotName:     updatedSlot.LotName,
		Code:        updatedSlot.Code,
		IsAvailable: updatedSlot.IsAvailable,
		SlotType:    updatedSlot.SlotType,
		Floor:       updatedSlot.Floor,
		CreatedAt:   updatedSlot.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return &item, nil
}

func (s *sParkingSlot) ParkingSlotDelete(ctx context.Context, req *entity.ParkingSlotDeleteReq) (*entity.ParkingSlotDeleteRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to delete the parking slot.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}
	if gconv.String(user.Map()["role"]) != consts.RoleAdmin {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Only admins can delete parking slots.")
	}

	slot, err := dao.ParkingSlots.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to find the parking slot. Please try again.")
	}
	if slot.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "The parking slot could not be found.")
	}

	count, err := dao.ParkingOrders.Ctx(ctx).
		Where("slot_id", req.Id).
		Where("status NOT IN (?)", g.Slice{"completed", "canceled"}).
		Where("deleted_at IS NULL").
		Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to check for active bookings. Please try again.")
	}
	if count > 0 {
		return nil, gerror.NewCode(consts.CodeInvalidOperation, "This parking slot cannot be deleted because it has active bookings.")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the parking slot. Please try again later.")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = dao.ParkingSlots.Ctx(ctx).TX(tx).Data(g.Map{
		"deleted_at": gtime.Now(),
		"updated_at": gtime.Now(),
	}).Where("id", req.Id).Where("deleted_at IS NULL").Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the parking slot. Please try again later.")
	}

	adminUsers, err := dao.Users.Ctx(ctx).Where("role", consts.RoleAdmin).Where("deleted_at IS NULL").All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the parking slot. Please try again later.")
	}

	for _, admin := range adminUsers {
		notiData := do.Notifications{
			UserId:         gconv.Int(admin.Map()["id"]),
			Type:           "parking_slot_deleted",
			Content:        fmt.Sprintf("Parking slot #%d (%s) has been deleted.", req.Id, slot.Map()["code"]),
			RelatedOrderId: req.Id,
			IsRead:         false,
			CreatedAt:      gtime.Now(),
		}
		_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the parking slot. Please try again later.")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the parking slot. Please try again later.")
	}

	return &entity.ParkingSlotDeleteRes{Message: "Parking slot deleted successfully"}, nil
}
