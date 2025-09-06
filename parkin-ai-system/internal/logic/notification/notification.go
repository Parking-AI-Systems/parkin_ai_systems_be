package notification

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
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}

	m := dao.Notifications.Ctx(ctx).
		Fields("notifications.*, parking_lots.name as lot_name, parking_orders.order_number as order_number, others_service.name as service_name").
		LeftJoin("parking_lots", "parking_lots.id = notifications.related_order_id AND notifications.type LIKE 'parking_lot%'").
		LeftJoin("parking_orders", "parking_orders.id = notifications.related_order_id AND notifications.type LIKE 'parking_order%'").
		LeftJoin("others_service", "others_service.id = notifications.related_order_id AND notifications.type LIKE 'others_service%'").
		Where("notifications.user_id", userID)

	if req.IsRead != nil {
		m = m.Where("notifications.is_read", *req.IsRead)
	}

	total, err := m.Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error counting notifications")
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
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving notifications")
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
			CreatedAt:      n.CreatedAt.Format("2006-01-02 15:04:05"),
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
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}

	var notification struct {
		entity.Notifications
		LotName     string `json:"lot_name"`
		OrderNumber string `json:"order_number"`
		ServiceName string `json:"service_name"`
	}
	err = dao.Notifications.Ctx(ctx).
		Fields("notifications.*, parking_lots.name as lot_name, parking_orders.order_number as order_number, others_service.name as service_name").
		LeftJoin("parking_lots", "parking_lots.id = notifications.related_order_id AND notifications.type LIKE 'parking_lot%'").
		LeftJoin("parking_orders", "parking_orders.id = notifications.related_order_id AND notifications.type LIKE 'parking_order%'").
		LeftJoin("others_service", "others_service.id = notifications.related_order_id AND notifications.type LIKE 'others_service%'").
		Where("notifications.id", req.Id).
		Where("notifications.user_id", userID).
		Scan(&notification)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving notification")
	}
	if notification.Id == 0 {
		return nil, gerror.NewCode(consts.CodeNotFound, "Notification not found or not authorized")
	}

	item := entity.NotificationItem{
		Id:             notification.Id,
		UserId:         notification.UserId,
		Type:           notification.Type,
		Content:        notification.Content,
		RelatedOrderId: notification.RelatedOrderId,
		IsRead:         notification.IsRead,
		CreatedAt:      notification.CreatedAt.Format("2006-01-02 15:04:05"),
		RelatedInfo:    getRelatedInfo(notification.Type, notification.LotName, notification.OrderNumber, notification.ServiceName),
	}

	return &item, nil
}

func (s *sNotification) NotificationMarkRead(ctx context.Context, req *entity.NotificationMarkReadReq) (*entity.NotificationMarkReadRes, error) {
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

	if len(req.Ids) == 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "At least one notification ID is required")
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

	count, err := dao.Notifications.Ctx(ctx).TX(tx).
		Where("id IN (?)", req.Ids).
		Where("user_id", userID).
		Where("is_read", false).
		Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking notifications")
	}
	if count == 0 {
		return nil, gerror.NewCode(consts.CodeNotFound, "No unread notifications found or not authorized")
	}

	_, err = dao.Notifications.Ctx(ctx).TX(tx).
		Data(g.Map{"is_read": true}).
		Where("id IN (?)", req.Ids).
		Where("user_id", userID).
		Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error marking notifications as read")
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error committing transaction")
	}

	return &entity.NotificationMarkReadRes{Message: fmt.Sprintf("%d notifications marked as read", count)}, nil
}

func (s *sNotification) NotificationDelete(ctx context.Context, req *entity.NotificationDeleteReq) (*entity.NotificationDeleteRes, error) {
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
	if gconv.String(user.Map()["role"]) != "admin" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Only admins can delete notifications")
	}

	notification, err := dao.Notifications.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking notification")
	}
	if notification.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "Notification not found")
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

	_, err = dao.Notifications.Ctx(ctx).TX(tx).Where("id", req.Id).Delete()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error deleting notification")
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
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating admin notification")
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error committing transaction")
	}

	return &entity.NotificationDeleteRes{Message: "Notification deleted successfully"}, nil
}
