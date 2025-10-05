package payment

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/payOSHQ/payos-lib-golang"

	"parkin-ai-system/internal/config"
	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

type sPayment struct{}

var (
	paymentInstance *sPayment
	initOnce        sync.Once
)

func Init() {
	initOnce.Do(func() {
		cfg := config.GetConfig()

		// Khởi tạo PayOS với thông tin cấu hình
		payos.Key(cfg.PayOs.ClientID, cfg.PayOs.ApiKey, cfg.PayOs.CheckSum)

		paymentInstance = &sPayment{}
		service.RegisterPayment(paymentInstance)
	})
}

func init() {
	// Register a placeholder that will initialize lazily
	service.RegisterPayment(&sPayment{})
}

// GenerateNumber tạo order code tự động theo thời gian
func GenerateNumber() int {
	millis := time.Now().UnixNano() / int64(time.Millisecond)
	millisStr := strconv.FormatInt(millis, 10)
	number, _ := strconv.Atoi(millisStr[len(millisStr)-6:])
	return number
}

// CheckoutAdd tạo payment request theo API definition
func (s *sPayment) CheckoutAdd(ctx context.Context, reqInterface interface{}) (interface{}, error) {
	// Ensure PayOS is initialized
	Init()

	// Convert interface{} to payos.CheckoutRequestType
	reqData, ok := reqInterface.(map[string]interface{})
	if !ok {
		return nil, gerror.New("invalid checkout request type")
	}

	// Tạo PayOS request
	payosReq := payos.CheckoutRequestType{
		OrderCode:   int64(gconv.Int(reqData["orderCode"])),
		Amount:      gconv.Int(reqData["amount"]),
		Description: gconv.String(reqData["description"]),
		CancelUrl:   gconv.String(reqData["cancelUrl"]),
		ReturnUrl:   gconv.String(reqData["returnUrl"]),
	}

	// Thêm items
	if itemsData, exists := reqData["items"]; exists {
		if itemsList, ok := itemsData.([]interface{}); ok {
			for _, itemData := range itemsList {
				if itemMap, ok := itemData.(map[string]interface{}); ok {
					payosReq.Items = append(payosReq.Items, payos.Item{
						Name:     gconv.String(itemMap["name"]),
						Price:    gconv.Int(itemMap["price"]),
						Quantity: 1, // Default quantity
					})
				}
			}
		}
	}

	// Thêm thông tin buyer nếu có
	if buyerName, exists := reqData["buyerName"]; exists {
		name := gconv.String(buyerName)
		payosReq.BuyerName = &name
	}
	if buyerEmail, exists := reqData["buyerEmail"]; exists {
		email := gconv.String(buyerEmail)
		payosReq.BuyerEmail = &email
	}
	if buyerPhone, exists := reqData["buyerPhone"]; exists {
		phone := gconv.String(buyerPhone)
		payosReq.BuyerPhone = &phone
	}

	// Gọi PayOS API để tạo payment link
	result, err := payos.CreatePaymentLink(payosReq)
	if err != nil {
		return nil, gerror.Wrap(err, "failed to create PayOS payment link")
	}

	// Tạo QR code link từ checkout URL
	qrCodeLink := fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=400x400&format=png&data=%s",
		url.QueryEscape(result.CheckoutUrl))

	g.Log().Info(ctx, "PayOS payment link created",
		"orderCode", payosReq.OrderCode,
		"paymentLinkId", result.PaymentLinkId,
		"amount", payosReq.Amount)

	// Trả về response với QR code link
	return map[string]interface{}{
		"paymentLinkId": result.PaymentLinkId,
		"checkoutUrl":   result.CheckoutUrl,
		"qrCode":        qrCodeLink, // QR code image URL
		"amount":        result.Amount,
		"orderCode":     result.OrderCode,
	}, nil
}

// PaymentLinkGet lấy thông tin payment link từ PayOS
func (s *sPayment) PaymentLinkGet(ctx context.Context, paymentLinkId string) (interface{}, error) {
	// Ensure PayOS is initialized
	Init()

	// Gọi PayOS API để lấy thông tin payment link
	result, err := payos.GetPaymentLinkInformation(paymentLinkId)
	if err != nil {
		return nil, gerror.Wrap(err, "failed to get payment link information")
	}

	g.Log().Info(ctx, "PayOS payment link retrieved",
		"paymentLinkId", paymentLinkId,
		"status", result.Status)

	return result, nil
}

// RefundAdd tạo refund request (tạm thời trả về thông tin mock vì PayOS library chưa hỗ trợ)
func (s *sPayment) RefundAdd(ctx context.Context, paymentLinkId string, amount int, reason *string) (interface{}, error) {
	// Ensure PayOS is initialized
	Init()

	// PayOS library hiện tại chưa có method refund, nên tạm thời return mock data
	g.Log().Info(ctx, "Refund request logged",
		"paymentLinkId", paymentLinkId,
		"amount", amount,
		"reason", reason)

	// Mock refund response
	refundResult := map[string]interface{}{
		"refundId": fmt.Sprintf("RF_%d_%s", time.Now().Unix(), paymentLinkId),
		"status":   "processing",
	}

	return refundResult, nil
}

// CreatePaymentLink implement interface method
func (s *sPayment) CreatePaymentLink(ctx context.Context, orderType string, orderID int64) (interface{}, error) {
	// Ensure PayOS is initialized
	Init()

	cfg := config.GetConfig()

	var orderCode int64
	var amount int
	var description string
	var items []entity.Item
	var buyerName, buyerEmail, buyerPhone string

	switch orderType {
	case "parking":
		// Lấy thông tin parking order
		parkingOrder, err := dao.ParkingOrders.Ctx(ctx).
			Fields("parking_orders.*, parking_lots.name as lot_name, parking_slots.code as slot_code, users.full_name, users.email, users.phone").
			LeftJoin("parking_lots", "parking_lots.id = parking_orders.lot_id").
			LeftJoin("parking_slots", "parking_slots.id = parking_orders.slot_id").
			LeftJoin("users", "users.id = parking_orders.user_id").
			Where("parking_orders.id", orderID).
			Where("parking_orders.deleted_at IS NULL").
			One()

		if err != nil {
			return nil, gerror.Wrap(err, "failed to get parking order")
		}
		if parkingOrder.IsEmpty() {
			return nil, gerror.New("parking order not found")
		}

		orderMap := parkingOrder.Map()
		orderCode = gconv.Int64(orderMap["id"])
		amount = int(gconv.Float64(orderMap["price"]) * 100) // Chuyển sang xu (VND)
		description = fmt.Sprintf("Thanh toán chỗ đậu xe %s - %s",
			gconv.String(orderMap["lot_name"]),
			gconv.String(orderMap["slot_code"]))
		if len([]rune(description)) > 25 {
			description = string([]rune(description)[:25])
		}
		items = []entity.Item{
			{
				Name:  fmt.Sprintf("Chỗ đậu xe %s", gconv.String(orderMap["slot_code"])),
				Price: amount,
			},
		}

		buyerName = gconv.String(orderMap["full_name"])
		buyerEmail = gconv.String(orderMap["email"])
		buyerPhone = gconv.String(orderMap["phone"])

	case "service":
		// Lấy thông tin service order
		serviceOrder, err := dao.OthersServiceOrders.Ctx(ctx).
			Fields("others_service_orders.*, parking_lots.name as lot_name, others_service.name as service_name, users.full_name, users.email, users.phone").
			LeftJoin("parking_lots", "parking_lots.id = others_service_orders.lot_id").
			LeftJoin("others_service", "others_service.id = others_service_orders.service_id").
			LeftJoin("users", "users.id = others_service_orders.user_id").
			Where("others_service_orders.id", orderID).
			Where("others_service_orders.deleted_at IS NULL").
			One()

		if err != nil {
			return nil, gerror.Wrap(err, "failed to get service order")
		}
		if serviceOrder.IsEmpty() {
			return nil, gerror.New("service order not found")
		}

		orderMap := serviceOrder.Map()
		orderCode = gconv.Int64(orderMap["id"]) + 1000000    // Thêm prefix để phân biệt với parking order
		amount = int(gconv.Float64(orderMap["price"]) * 100) // Chuyển sang xu (VND)
		description = fmt.Sprintf("Thanh toán dịch vụ %s tại %s",
			gconv.String(orderMap["service_name"]),
			gconv.String(orderMap["lot_name"]))
		if len([]rune(description)) > 25 {
			description = string([]rune(description)[:25])
		}
		items = []entity.Item{
			{
				Name:  gconv.String(orderMap["service_name"]),
				Price: amount,
			},
		}

		buyerName = gconv.String(orderMap["full_name"])
		buyerEmail = gconv.String(orderMap["email"])
		buyerPhone = gconv.String(orderMap["phone"])

	default:
		return nil, gerror.New("invalid order type")
	}

	// Tạo PayOS request
	payosReq := payos.CheckoutRequestType{
		OrderCode:   orderCode,
		Amount:      amount,
		Description: description,
		CancelUrl:   fmt.Sprintf("%s/payment/cancel", cfg.PayOs.ParkinDomain),
		ReturnUrl:   fmt.Sprintf("%s/payment/success", cfg.PayOs.ParkinDomain),
	}

	// Thêm items
	for _, item := range items {
		payosReq.Items = append(payosReq.Items, payos.Item{
			Name:     item.Name,
			Price:    item.Price,
			Quantity: 1,
		})
	}

	// Thêm thông tin buyer nếu có
	if buyerName != "" {
		payosReq.BuyerName = &buyerName
	}
	if buyerEmail != "" {
		payosReq.BuyerEmail = &buyerEmail
	}
	if buyerPhone != "" {
		payosReq.BuyerPhone = &buyerPhone
	}

	// Gọi PayOS API
	result, err := payos.CreatePaymentLink(payosReq)
	if err != nil {
		return nil, gerror.Wrap(err, "failed to create PayOS payment link")
	}

	// Tạo QR code link từ checkout URL sử dụng QR Server API
	qrCodeLink := fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=400x400&format=png&data=%s",
		url.QueryEscape(result.CheckoutUrl))

	g.Log().Info(ctx, "PayOS payment link created",
		"orderType", orderType,
		"orderID", orderID,
		"paymentLinkId", result.PaymentLinkId,
		"amount", amount)

	// Trả về response với QR code link thay vì raw QR data
	return map[string]interface{}{
		"paymentLinkId": result.PaymentLinkId,
		"checkoutUrl":   result.CheckoutUrl,
		"qrCode":        qrCodeLink, // QR code image URL
		"amount":        amount,
		"orderCode":     orderCode,
	}, nil
}

// HandlePaymentWebhook xử lý webhook từ PayOS
func (s *sPayment) HandlePaymentWebhook(ctx context.Context, webhookData interface{}) error {
	// Ensure PayOS is initialized
	Init()

	// Convert interface{} to payos.WebhookType
	webhook, ok := webhookData.(payos.WebhookType)
	if !ok {
		return gerror.New("invalid webhook data type")
	}

	// Verify webhook signature using PayOS library
	verifiedData, err := payos.VerifyPaymentWebhookData(webhook)
	if err != nil {
		return gerror.Wrap(err, "failed to verify webhook signature")
	}

	if verifiedData == nil {
		g.Log().Warning(ctx, "PayOS webhook verification failed")
		return gerror.New("webhook verification failed")
	}

	orderCode := verifiedData.OrderCode

	// Xác định loại order và ID thực tế
	var orderType string
	var realOrderID int64

	if orderCode >= 1000000 {
		// Service order (có prefix 1000000)
		orderType = "service"
		realOrderID = orderCode - 1000000
	} else {
		// Parking order
		orderType = "parking"
		realOrderID = orderCode
	}

	g.Log().Info(ctx, "Processing PayOS webhook",
		"orderType", orderType,
		"orderID", realOrderID,
		"paymentLinkId", verifiedData.PaymentLinkId,
		"amount", verifiedData.Amount,
		"status", "completed")

	// Bắt đầu transaction
	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return gerror.Wrap(err, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	switch orderType {
	case "parking":
		err = s.updateParkingOrderPayment(ctx, tx, realOrderID, verifiedData)
	case "service":
		err = s.updateServiceOrderPayment(ctx, tx, realOrderID, verifiedData)
	default:
		return gerror.New("unknown order type")
	}

	if err != nil {
		return err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return gerror.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// updateParkingOrderPayment cập nhật trạng thái thanh toán cho parking order
func (s *sPayment) updateParkingOrderPayment(ctx context.Context, tx gdb.TX, orderID int64, data *payos.WebhookDataType) error {
	// Cập nhật trạng thái thanh toán
	_, err := dao.ParkingOrders.Ctx(ctx).TX(tx).
		Data(g.Map{
			"payment_status": "completed",
			"updated_at":     time.Now(),
		}).
		Where("id", orderID).
		Where("deleted_at IS NULL").
		Update()

	if err != nil {
		return gerror.Wrap(err, "failed to update parking order payment status")
	}

	// Tạo thông báo
	_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(g.Map{
		"user_id":          g.Map{"SELECT user_id FROM parking_orders WHERE id = ?": orderID},
		"type":             "payment_completed",
		"content":          fmt.Sprintf("Thanh toán đơn hàng đậu xe #%d đã hoàn thành thành công. Số tiền: %d VND", orderID, data.Amount/100),
		"related_order_id": orderID,
		"is_read":          false,
		"created_at":       time.Now(),
	}).Insert()

	if err != nil {
		return gerror.Wrap(err, "failed to create payment notification")
	}

	return nil
}

// updateServiceOrderPayment cập nhật trạng thái thanh toán cho service order
func (s *sPayment) updateServiceOrderPayment(ctx context.Context, tx gdb.TX, orderID int64, data *payos.WebhookDataType) error {
	// Cập nhật trạng thái thanh toán
	_, err := dao.OthersServiceOrders.Ctx(ctx).TX(tx).
		Data(g.Map{
			"payment_status": "completed",
			"updated_at":     time.Now(),
		}).
		Where("id", orderID).
		Where("deleted_at IS NULL").
		Update()

	if err != nil {
		return gerror.Wrap(err, "failed to update service order payment status")
	}

	// Tạo thông báo
	_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(g.Map{
		"user_id":          g.Map{"SELECT user_id FROM others_service_orders WHERE id = ?": orderID},
		"type":             "payment_completed",
		"content":          fmt.Sprintf("Thanh toán đơn hàng dịch vụ #%d đã hoàn thành thành công. Số tiền: %d VND", orderID, data.Amount/100),
		"related_order_id": orderID,
		"is_read":          false,
		"created_at":       time.Now(),
	}).Insert()

	if err != nil {
		return gerror.Wrap(err, "failed to create payment notification")
	}

	return nil
}
