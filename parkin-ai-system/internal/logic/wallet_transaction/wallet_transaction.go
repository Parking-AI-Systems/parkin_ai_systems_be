package wallet_transaction

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

type sWalletTransaction struct{}

func Init() {
	service.RegisterWalletTransaction(&sWalletTransaction{})
}
func init() {
	Init()
}

func (s *sWalletTransaction) WalletTransactionAdd(ctx context.Context, req *entity.WalletTransactionAddReq) (*entity.WalletTransactionAddRes, error) {
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
	if !isAdmin && req.UserId != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "You can only add transactions for yourself or must be an admin")
	}

	// Validate transaction type
	isValidType := false
	for _, validType := range consts.ValidTransactionTypes {
		if req.Type == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Invalid transaction type")
	}

	// Validate amount
	if (req.Type == consts.TransactionTypeDeposit || req.Type == consts.TransactionTypeRefund) && req.Amount <= 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Amount must be positive for deposit or refund")
	}
	if (req.Type == consts.TransactionTypePayment || req.Type == consts.TransactionTypeWithdrawal) && req.Amount >= 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Amount must be negative for payment or withdrawal")
	}

	// Check related order for payment/refund
	if req.RelatedOrderId != 0 {
		// Check parking_orders
		parkingOrder, err := dao.ParkingOrders.Ctx(ctx).Where("id", req.RelatedOrderId).One()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking parking order")
		}
		if !parkingOrder.IsEmpty() {
			// Ensure order belongs to the user (unless admin)
			if !isAdmin && gconv.Int64(parkingOrder.Map()["user_id"]) != gconv.Int64(userID) {
				return nil, gerror.NewCode(consts.CodeUnauthorized, "Parking order does not belong to you")
			}
			// For payment, validate vehicle-slot compatibility
			if req.Type == consts.TransactionTypePayment {
				vehicleId := gconv.Int64(parkingOrder.Map()["vehicle_id"])
				slotId := gconv.Int64(parkingOrder.Map()["slot_id"])
				if err := service.Vehicle().CheckVehicleSlotCompatibility(ctx, vehicleId, slotId); err != nil {
					return nil, err
				}
			}
		} else {
			// Check others_service_orders
			serviceOrder, err := dao.OthersServiceOrders.Ctx(ctx).Where("id", req.RelatedOrderId).One()
			if err != nil {
				return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking service order")
			}
			if serviceOrder.IsEmpty() {
				return nil, gerror.NewCode(consts.CodeNotFound, "Related order not found in parking or service orders")
			}
			// Ensure order belongs to the user (unless admin)
			if !isAdmin && gconv.Int64(serviceOrder.Map()["user_id"]) != gconv.Int64(userID) {
				return nil, gerror.NewCode(consts.CodeUnauthorized, "Service order does not belong to you")
			}
		}
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

	// Insert transaction
	data := do.WalletTransactions{
		UserId:         req.UserId,
		Amount:         req.Amount,
		Type:           req.Type,
		Description:    req.Description,
		RelatedOrderId: req.RelatedOrderId,
		CreatedAt:      gtime.Now(),
	}
	lastId, err := dao.WalletTransactions.Ctx(ctx).TX(tx).Data(data).InsertAndGetId()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating transaction")
	}

	// Notify admins
	adminUsers, err := dao.Users.Ctx(ctx).Where("role", "admin").All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving admins")
	}
	for _, admin := range adminUsers {
		notiData := do.Notifications{
			UserId:         gconv.Int64(admin.Map()["id"]),
			Type:           "wallet_transaction_added",
			Content:        fmt.Sprintf("New %s transaction #%d of %.2f for user #%d.", req.Type, lastId, req.Amount, req.UserId),
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

	return &entity.WalletTransactionAddRes{Id: lastId}, nil
}

func (s *sWalletTransaction) WalletTransactionList(ctx context.Context, req *entity.WalletTransactionListReq) (*entity.WalletTransactionListRes, error) {
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

	m := dao.WalletTransactions.Ctx(ctx).
		Fields("wallet_transactions.*, users.username, parking_orders.vehicle_id, others_service_orders.service_id").
		LeftJoin("users", "users.id = wallet_transactions.user_id").
		LeftJoin("parking_orders", "parking_orders.id = wallet_transactions.related_order_id").
		LeftJoin("others_service_orders", "others_service_orders.id = wallet_transactions.related_order_id")

	isAdmin := gconv.String(user.Map()["role"]) == "admin"
	if !isAdmin {
		m = m.Where("wallet_transactions.user_id", userID)
	}

	if req.Type != "" {
		m = m.Where("wallet_transactions.type", req.Type)
	}
	if req.Description != "" {
		m = m.WhereLike("wallet_transactions.description", "%"+req.Description+"%")
	}

	total, err := m.Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error counting transactions")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	m = m.Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize)

	var transactions []struct {
		entity.WalletTransactions
		Username  string `json:"username"`
		VehicleId int64  `json:"vehicle_id"`
		ServiceId int64  `json:"service_id"`
	}
	err = m.Order("wallet_transactions.id DESC").Scan(&transactions)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving transactions")
	}

	list := make([]entity.WalletTransactionItem, 0, len(transactions))
	for _, t := range transactions {
		var licensePlate string
		var serviceType string
		if t.VehicleId != 0 {
			vehicle, err := dao.Vehicles.Ctx(ctx).Where("id", t.VehicleId).One()
			if err == nil && !vehicle.IsEmpty() {
				licensePlate = gconv.String(vehicle.Map()["license_plate"])
			}
		}
		if t.ServiceId != 0 {
			service, err := dao.OthersService.Ctx(ctx).Where("id", t.ServiceId).One()
			if err == nil && !service.IsEmpty() {
				serviceType = gconv.String(service.Map()["service_type"])
			}
		}
		item := entity.WalletTransactionItem{
			Id:             t.Id,
			UserId:         t.UserId,
			Username:       t.Username,
			Amount:         t.Amount,
			Type:           t.Type,
			Description:    t.Description,
			RelatedOrderId: t.RelatedOrderId,
			LicensePlate:   licensePlate,
			ServiceType:    serviceType,
			CreatedAt:      t.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		list = append(list, item)
	}

	return &entity.WalletTransactionListRes{
		List:  list,
		Total: total,
	}, nil
}

func (s *sWalletTransaction) WalletTransactionGet(ctx context.Context, req *entity.WalletTransactionGetReq) (*entity.WalletTransactionItem, error) {
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

	var transaction struct {
		entity.WalletTransactions
		Username  string `json:"username"`
		VehicleId int64  `json:"vehicle_id"`
		ServiceId int64  `json:"service_id"`
	}
	err = dao.WalletTransactions.Ctx(ctx).
		Fields("wallet_transactions.*, users.username, parking_orders.vehicle_id, others_service_orders.service_id").
		LeftJoin("users", "users.id = wallet_transactions.user_id").
		LeftJoin("parking_orders", "parking_orders.id = wallet_transactions.related_order_id").
		LeftJoin("others_service_orders", "others_service_orders.id = wallet_transactions.related_order_id").
		Where("wallet_transactions.id", req.Id).
		Scan(&transaction)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving transaction")
	}
	if transaction.Id == 0 {
		return nil, gerror.NewCode(consts.CodeNotFound, "Transaction not found")
	}

	isAdmin := gconv.String(user.Map()["role"]) == "admin"
	if !isAdmin && transaction.UserId != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "You can only view your own transactions or must be an admin")
	}

	var licensePlate string
	var serviceType string
	if transaction.VehicleId != 0 {
		vehicle, err := dao.Vehicles.Ctx(ctx).Where("id", transaction.VehicleId).One()
		if err == nil && !vehicle.IsEmpty() {
			licensePlate = gconv.String(vehicle.Map()["license_plate"])
		}
	}
	if transaction.ServiceId != 0 {
		service, err := dao.OthersService.Ctx(ctx).Where("id", transaction.ServiceId).One()
		if err == nil && !service.IsEmpty() {
			serviceType = gconv.String(service.Map()["service_type"])
		}
	}

	item := entity.WalletTransactionItem{
		Id:             transaction.Id,
		UserId:         transaction.UserId,
		Username:       transaction.Username,
		Amount:         transaction.Amount,
		Type:           transaction.Type,
		Description:    transaction.Description,
		RelatedOrderId: transaction.RelatedOrderId,
		LicensePlate:   licensePlate,
		ServiceType:    serviceType,
		CreatedAt:      transaction.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return &item, nil
}
