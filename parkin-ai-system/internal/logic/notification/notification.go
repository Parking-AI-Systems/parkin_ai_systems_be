package notification

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

type sNotification struct{}

func Init() {
	service.RegisterNotification(&sNotification{})
}
func init() {
	Init()
}

func (s *sNotification) NotificationList(ctx context.Context, req *entity.NotificationListReq) (*entity.NotificationListRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to view your notifications.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}

	// Base query builder for joins and where conditions
	baseQuery := dao.Notifications.Ctx(ctx).
		LeftJoin("parking_lots", "parking_lots.id = notifications.related_order_id AND notifications.type LIKE 'parking_lot%'").
		LeftJoin("parking_orders", "parking_orders.id = notifications.related_order_id AND notifications.type LIKE 'parking_order%'").
		LeftJoin("others_service", "others_service.id = notifications.related_order_id AND notifications.type LIKE 'others_service%'").
		Where("notifications.user_id", userID).
		Where("notifications.deleted_at IS NULL").
		Where("parking_lots.deleted_at IS NULL OR parking_lots.id IS NULL").
		Where("parking_orders.deleted_at IS NULL OR parking_orders.id IS NULL").
		Where("others_service.deleted_at IS NULL OR others_service.id IS NULL")

	if req.IsRead != nil {
		baseQuery = baseQuery.Where("notifications.is_read", *req.IsRead)
	}

	// Count query - no fields needed for count
	total, err := baseQuery.Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load your notifications. Please try again later.")
	}

	// Data query with fields
	m := dao.Notifications.Ctx(ctx).
		Fields("notifications.*, parking_lots.name as lot_name, parking_orders.id as order_number, others_service.name as service_name").
		LeftJoin("parking_lots", "parking_lots.id = notifications.related_order_id AND notifications.type LIKE 'parking_lot%'").
		LeftJoin("parking_orders", "parking_orders.id = notifications.related_order_id AND notifications.type LIKE 'parking_order%'").
		LeftJoin("others_service", "others_service.id = notifications.related_order_id AND notifications.type LIKE 'others_service%'").
		Where("notifications.user_id", userID).
		Where("notifications.deleted_at IS NULL").
		Where("parking_lots.deleted_at IS NULL OR parking_lots.id IS NULL").
		Where("parking_orders.deleted_at IS NULL OR parking_orders.id IS NULL").
		Where("others_service.deleted_at IS NULL OR others_service.id IS NULL")

	if req.IsRead != nil {
		m = m.Where("notifications.is_read", *req.IsRead)
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	m = m.Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize)

	var notifications []struct {
		entity.Notifications
		LotName     string `json:"lot_name"`
		OrderNumber string `json:"order_number"`
		ServiceName string `json:"service_name"`
	}
	err = m.Order("notifications.id DESC").Scan(&notifications)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load your notifications. Please try again later.")
	}

	list := make([]entity.NotificationItem, 0, len(notifications))
	for _, n := range notifications {
		item := entity.NotificationItem{
			Id:             n.Id,
			UserId:         n.UserId,
			Type:           n.Type,
			Content:        n.Content,
			RelatedOrderId: n.RelatedOrderId,
			IsRead:         n.IsRead,
			CreatedAt:      time.Time(n.CreatedAt.Time).Format("2006-01-02 15:04:05"),
			RelatedInfo:    getRelatedInfo(n.Type, n.LotName, n.OrderNumber, n.ServiceName),
		}
		list = append(list, item)
	}

	return &entity.NotificationListRes{
		List:  list,
		Total: total,
	}, nil
}

func getRelatedInfo(notificationType, lotName, orderNumber, serviceName string) string {
	switch {
	case notificationType == "parking_lot_created" || notificationType == "parking_lot_updated" || notificationType == "parking_lot_deleted" || notificationType == "parking_lot_image_deleted":
		if lotName != "" {
			return fmt.Sprintf("Parking Lot: %s", lotName)
		}
	case notificationType == "parking_order_confirmed" || notificationType == "parking_order_cancelled":
		if orderNumber != "" {
			return fmt.Sprintf("Order: %s", orderNumber)
		}
	case notificationType == "others_service_added" || notificationType == "others_service_updated" || notificationType == "others_service_deleted":
		if serviceName != "" {
			return fmt.Sprintf("Service: %s", serviceName)
		}
	}
	return ""
}

func (s *sNotification) NotificationGet(ctx context.Context, req *entity.NotificationGetReq) (*entity.NotificationItem, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to view the notification.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}

	var notification struct {
		entity.Notifications
		LotName     string `json:"lot_name"`
		OrderNumber string `json:"order_number"`
		ServiceName string `json:"service_name"`
	}
	err = dao.Notifications.Ctx(ctx).
		Fields("notifications.*, parking_lots.name as lot_name, parking_orders.id as order_number, others_service.name as service_name").
		LeftJoin("parking_lots", "parking_lots.id = notifications.related_order_id AND notifications.type LIKE 'parking_lot%'").
		LeftJoin("parking_orders", "parking_orders.id = notifications.related_order_id AND notifications.type LIKE 'parking_order%'").
		LeftJoin("others_service", "others_service.id = notifications.related_order_id AND notifications.type LIKE 'others_service%'").
		Where("notifications.id", req.Id).
		Where("notifications.user_id", userID).
		Where("notifications.deleted_at IS NULL").
		Where("parking_lots.deleted_at IS NULL OR parking_lots.id IS NULL").
		Where("parking_orders.deleted_at IS NULL OR parking_orders.id IS NULL").
		Where("others_service.deleted_at IS NULL OR others_service.id IS NULL").
		Scan(&notification)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load the notification. Please try again later.")
	}
	if notification.Id == 0 {
		return nil, gerror.NewCode(consts.CodeNotFound, "The notification could not be found or you are not authorized to view it.")
	}

	item := entity.NotificationItem{
		Id:             notification.Id,
		UserId:         notification.UserId,
		Type:           notification.Type,
		Content:        notification.Content,
		RelatedOrderId: notification.RelatedOrderId,
		IsRead:         notification.IsRead,
		CreatedAt:      time.Time(notification.CreatedAt.Time).Format("2006-01-02 15:04:05"),
		RelatedInfo:    getRelatedInfo(notification.Type, notification.LotName, notification.OrderNumber, notification.ServiceName),
	}

	return &item, nil
}

func (s *sNotification) NotificationMarkRead(ctx context.Context, req *entity.NotificationMarkReadReq) (*entity.NotificationMarkReadRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to mark notifications as read.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}

	if len(req.Ids) == 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Please provide at least one notification ID.")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while marking notifications as read. Please try again later.")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	count, err := dao.Notifications.Ctx(ctx).TX(tx).
		Where("id IN (?)", req.Ids).
		Where("user_id", userID).
		Where("is_read", false).
		Where("deleted_at IS NULL").
		Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to find the notifications. Please try again.")
	}
	if count == 0 {
		return nil, gerror.NewCode(consts.CodeNotFound, "No unread notifications were found, or you are not authorized to mark them.")
	}

	_, err = dao.Notifications.Ctx(ctx).TX(tx).
		Data(g.Map{
			"is_read":    true,
			"updated_at": gtime.Now(),
		}).
		Where("id IN (?)", req.Ids).
		Where("user_id", userID).
		Where("deleted_at IS NULL").
		Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while marking notifications as read. Please try again later.")
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while marking notifications as read. Please try again later.")
	}

	return &entity.NotificationMarkReadRes{Message: fmt.Sprintf("%d notifications marked as read", count)}, nil
}

func (s *sNotification) NotificationDelete(ctx context.Context, req *entity.NotificationDeleteReq) (*entity.NotificationDeleteRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to delete a notification.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}
	if gconv.String(user.Map()["role"]) != consts.RoleAdmin {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Only admins can delete notifications.")
	}

	notification, err := dao.Notifications.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to find the notification. Please try again.")
	}
	if notification.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "The notification could not be found.")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the notification. Please try again later.")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(g.Map{
		"deleted_at": gtime.Now(),
		"updated_at": gtime.Now(),
	}).Where("id", req.Id).Where("deleted_at IS NULL").Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the notification. Please try again later.")
	}

	adminNotiData := do.Notifications{
		UserId:         userID,
		Type:           "notification_deleted",
		Content:        fmt.Sprintf("Notification #%d has been deleted.", req.Id),
		RelatedOrderId: req.Id,
		IsRead:         false,
		CreatedAt:      gtime.Now(),
	}
	_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(adminNotiData).Insert()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the notification. Please try again later.")
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the notification. Please try again later.")
	}

	return &entity.NotificationDeleteRes{Message: "Notification deleted successfully"}, nil
}
