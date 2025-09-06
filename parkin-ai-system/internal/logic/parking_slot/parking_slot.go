package parking_slot

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
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}
	if gconv.String(user.Map()["role"]) != consts.RoleAdmin {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Only admins can add parking slots")
	}

	lot, err := dao.ParkingLots.Ctx(ctx).Where("id", req.LotId).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking parking lot")
	}
	if lot.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeParkingLotNotFound, "Parking lot not found")
	}

	count, err := dao.ParkingSlots.Ctx(ctx).Where("lot_id", req.LotId).Where("code", req.Code).Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking slot code uniqueness")
	}
	if count > 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Slot code already exists for this parking lot")
	}

	isValidSlotType := false
	for _, validType := range consts.ValidSlotTypes {
		if req.SlotType == validType {
			isValidSlotType = true
			break
		}
	}
	if !isValidSlotType {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Invalid slot type")
	}

	if len(req.Code) > 20 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Slot code must be less than 20 characters")
	}
	if len(req.Floor) > 10 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Floor must be less than 10 characters")
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
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating parking slot")
	}

	adminUsers, err := dao.Users.Ctx(ctx).Where("role", "admin").All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving admins")
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
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating notification")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error committing transaction")
	}

	return &entity.ParkingSlotAddRes{Id: lastId}, nil
}

func (s *sParkingSlot) ParkingSlotList(ctx context.Context, req *entity.ParkingSlotListReq) (*entity.ParkingSlotListRes, error) {
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

	// Build base query conditions
	baseQuery := dao.ParkingSlots.Ctx(ctx).
		LeftJoin("parking_lots", "parking_lots.id = parking_slots.lot_id")

	if req.LotId != 0 {
		baseQuery = baseQuery.Where("parking_slots.lot_id", req.LotId)
	}
	if req.IsAvailable != nil {
		baseQuery = baseQuery.Where("parking_slots.is_available", *req.IsAvailable)
	}
	if req.SlotType != "" {
		baseQuery = baseQuery.Where("parking_slots.slot_type", req.SlotType)
	}

	// Count query - use simple field
	total, err := baseQuery.Fields("parking_slots.id").Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error counting parking slots")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// Data query - use joined fields
	var slots []struct {
		entity.ParkingSlots
		LotName string `json:"lot_name"`
	}
	err = baseQuery.Fields("parking_slots.*, parking_lots.name as lot_name").
		Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize).
		Order("parking_slots.id DESC").Scan(&slots)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving parking slots")
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
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}

	var slot struct {
		entity.ParkingSlots
		LotName string `json:"lot_name"`
	}
	err = dao.ParkingSlots.Ctx(ctx).
		Fields("parking_slots.*, parking_lots.name as lot_name").
		LeftJoin("parking_lots", "parking_lots.id = parking_slots.lot_id").
		Where("parking_slots.id", req.Id).
		Scan(&slot)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving parking slot")
	}
	if slot.Id == 0 {
		return nil, gerror.NewCode(consts.CodeNotFound, "Parking slot not found")
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
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}
	if gconv.String(user.Map()["role"]) != consts.RoleAdmin {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Only admins can update parking slots")
	}

	slot, err := dao.ParkingSlots.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking parking slot")
	}
	if slot.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "Parking slot not found")
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

	if req.Code != "" {
		count, err := dao.ParkingSlots.Ctx(ctx).
			Where("lot_id", req.LotId).
			Where("code", req.Code).
			Where("id != ?", req.Id).
			Count()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking slot code uniqueness")
		}
		if count > 0 {
			return nil, gerror.NewCode(consts.CodeInvalidInput, "Slot code already exists for this parking lot")
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
			return nil, gerror.NewCode(consts.CodeInvalidInput, "Invalid slot type")
		}
	}

	if req.Code != "" && len(req.Code) > 20 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Slot code must be less than 20 characters")
	}
	if req.Floor != "" && len(req.Floor) > 10 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Floor must be less than 10 characters")
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

	_, err = dao.ParkingSlots.Ctx(ctx).TX(tx).Data(updateData).Where("id", req.Id).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error updating parking slot")
	}

	adminUsers, err := dao.Users.Ctx(ctx).Where("role", "admin").All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving admins")
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
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating notification")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error committing transaction")
	}

	var updatedSlot struct {
		entity.ParkingSlots
		LotName string `json:"lot_name"`
	}
	err = dao.ParkingSlots.Ctx(ctx).
		Fields("parking_slots.*, parking_lots.name as lot_name").
		LeftJoin("parking_lots", "parking_lots.id = parking_slots.lot_id").
		Where("parking_slots.id", req.Id).
		Scan(&updatedSlot)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving updated parking slot")
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
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}
	if gconv.String(user.Map()["role"]) != consts.RoleAdmin {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Only admins can delete parking slots")
	}

	slot, err := dao.ParkingSlots.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking parking slot")
	}
	if slot.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "Parking slot not found")
	}

	count, err := dao.ParkingOrders.Ctx(ctx).
		Where("slot_id", req.Id).
		Where("status NOT IN (?)", g.Slice{"completed", "cancelled"}).
		Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking active orders")
	}
	if count > 0 {
		return nil, gerror.NewCode(consts.CodeInvalidOperation, "Cannot delete parking slot with active orders")
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

	_, err = dao.ParkingSlots.Ctx(ctx).TX(tx).Where("id", req.Id).Delete()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error deleting parking slot")
	}

	adminUsers, err := dao.Users.Ctx(ctx).Where("role", "admin").All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving admins")
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
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating notification")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error committing transaction")
	}

	return &entity.ParkingSlotDeleteRes{Message: "Parking slot deleted successfully"}, nil
}
