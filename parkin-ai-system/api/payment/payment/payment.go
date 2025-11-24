package payment

import "github.com/gogf/gf/v2/frame/g"

// CheckoutAddReq represents the request payload for creating a payment request.
type CheckoutAddReq struct {
	g.Meta       `path:"/payment-requests" tags:"Payment" method:"POST" summary:"Create a payment request" description:"Creates a new payment request for processing payments. Requires authentication." middleware:"middleware.Auth"`
	OrderCode    int64   `json:"orderCode" v:"required|min:1#Order code is required|Order code must be positive"`
	Amount       int     `json:"amount" v:"required|min:1000#Amount is required|Amount must be at least 1000 VND"`
	Description  string  `json:"description" v:"required|length:1,255#Description is required|Description must be between 1 and 255 characters"`
	CancelUrl    string  `json:"cancelUrl" v:"required|url#Cancel URL is required|Cancel URL must be a valid URL"`
	ReturnUrl    string  `json:"returnUrl" v:"required|url#Return URL is required|Return URL must be a valid URL"`
	Items        []Item  `json:"items" v:"required#Items list is required"`
	BuyerName    *string `json:"buyerName" v:"length:0,100#Buyer name must be less than 100 characters"`
	BuyerEmail   *string `json:"buyerEmail" v:"email#Buyer email must be a valid email"`
	BuyerPhone   *string `json:"buyerPhone" v:"phone#Buyer phone must be a valid phone number"`
	BuyerAddress *string `json:"buyerAddress" v:"length:0,255#Buyer address must be less than 255 characters"`
	ExpiredAt    *int    `json:"expiredAt" v:"min:0#Expiration time must be non-negative"`
	Signature    *string `json:"signature" v:"length:0,64#Signature must be less than 64 characters"`
}

// Item represents a single item in the payment request.
type Item struct {
	Name  string `json:"name" v:"required|length:1,100#Item name is required|Item name must be between 1 and 100 characters"`
	Price int    `json:"price" v:"required|min:0#Price is required|Price must be non-negative"`
}

// CheckoutAddRes represents the response for a created payment request.
type CheckoutAddRes struct {
	PaymentLinkId string `json:"paymentLinkId"` // Payment link ID
	CheckoutUrl   string `json:"checkoutUrl"`   // URL for payment checkout
	QRCode        string `json:"qrCode"`        // Base64 encoded QR code for payment
}

// PaymentLinkGetReq represents the request payload for retrieving payment request details.
type PaymentLinkGetReq struct {
	g.Meta `path:"/payment-requests/:id" tags:"Payment" method:"GET" summary:"Get payment request details" description:"Retrieves details of a specific payment request by ID. Requires authentication." middleware:"middleware.Auth"`
	Id     string `json:"id" v:"required|length:1,50#Payment link ID is required|Payment link ID must be between 1 and 50 characters"`
}

// PaymentLinkGetRes represents the response for payment request details.
type PaymentLinkGetRes struct {
	PaymentLink PaymentLinkItem `json:"paymentLink"`
}

// PaymentLinkItem represents the details of a payment request.
type PaymentLinkItem struct {
	Id                 string            `json:"id"`                 // Payment link ID
	OrderCode          int64             `json:"orderCode"`          // Order code
	Amount             int               `json:"amount"`             // Total order amount
	AmountPaid         int               `json:"amountPaid"`         // Paid amount
	AmountRemaining    int               `json:"amountRemaining"`    // Remaining amount to be paid
	Status             string            `json:"status"`             // Payment status (PENDING, PAID, CANCELLED)
	CreatedAt          string            `json:"createdAt"`          // Creation time (ISO 8601)
	Transactions       []TransactionItem `json:"transactions"`       // List of transactions
	CancellationReason *string           `json:"cancellationReason"` // Reason for cancellation (if cancelled)
	CancelledAt        *string           `json:"cancelledAt"`        // Cancellation time (if cancelled)
}

// TransactionItem represents a single transaction in a payment request.
type TransactionItem struct {
	Reference              string  `json:"reference"`              // Transaction reference code
	Amount                 int     `json:"amount"`                 // Transaction amount
	AccountNumber          string  `json:"accountNumber"`          // Payment channel account number
	Description            string  `json:"description"`            // Transaction description
	TransactionDateTime    string  `json:"transactionDateTime"`    // Transaction time (ISO 8601)
	VirtualAccountName     *string `json:"virtualAccountName"`     // Virtual account name (optional)
	VirtualAccountNumber   *string `json:"virtualAccountNumber"`   // Virtual account number (optional)
	CounterAccountBankId   *string `json:"counterAccountBankId"`   // Counterparty bank ID (optional)
	CounterAccountBankName *string `json:"counterAccountBankName"` // Counterparty bank name (optional)
	CounterAccountName     *string `json:"counterAccountName"`     // Counterparty account name (optional)
	CounterAccountNumber   *string `json:"counterAccountNumber"`   // Counterparty account number (optional)
}

// RefundAddReq represents the request payload for refunding a payment.
type RefundAddReq struct {
	g.Meta    `path:"/payment-requests/:id/refund" tags:"Payment" method:"POST" summary:"Create a refund request" description:"Initiates a refund for a specific payment request. Admin only." middleware:"middleware.Auth"`
	Id        string  `json:"id" v:"required|length:1,50#Payment link ID is required|Payment link ID must be between 1 and 50 characters"`
	Amount    int     `json:"amount" v:"required|min:1000#Refund amount is required|Refund amount must be at least 1000 VND"`
	Reason    *string `json:"reason" v:"length:0,255#Reason must be less than 255 characters"`
	Signature *string `json:"signature" v:"length:0,64#Signature must be less than 64 characters"`
}

// RefundAddRes represents the response for a refund request.
type RefundAddRes struct {
	RefundId string `json:"refundId"` // Refund transaction ID
	Status   string `json:"status"`   // Refund status (PROCESSING, COMPLETED)
}

// WebhookReq represents the webhook payload received from payOS.
type WebhookReq struct {
	g.Meta    `path:"/webhook" tags:"Payment" method:"POST" summary:"Handle payment webhook" description:"Receives webhook notifications from payOS about payment status updates. Requires signature verification."`
	Code      string       `json:"code"`      // Error code
	Desc      string       `json:"desc"`      // Error description
	Success   bool         `json:"success"`   // Webhook status
	Data      *WebhookData `json:"data"`      // Webhook data
	Signature string       `json:"signature"` // Signature for data integrity
}

// WebhookData represents the data field in the webhook payload.
type WebhookData struct {
	OrderCode              int64   `json:"orderCode"`              // Order code
	Amount                 int     `json:"amount"`                 // Transaction amount
	Description            string  `json:"description"`            // Transaction description
	AccountNumber          string  `json:"accountNumber"`          // Payment channel account number
	Reference              string  `json:"reference"`              // Transaction reference code
	TransactionDateTime    string  `json:"transactionDateTime"`    // Transaction time (ISO 8601)
	Currency               string  `json:"currency"`               // Currency (e.g., VND)
	PaymentLinkId          string  `json:"paymentLinkId"`          // Payment link ID
	Code                   string  `json:"code"`                   // Transaction error code
	Desc                   string  `json:"desc"`                   // Transaction error description
	CounterAccountBankId   *string `json:"counterAccountBankId"`   // Counterparty bank ID (optional)
	CounterAccountBankName *string `json:"counterAccountBankName"` // Counterparty bank name (optional)
	CounterAccountName     *string `json:"counterAccountName"`     // Counterparty account name (optional)
	CounterAccountNumber   *string `json:"counterAccountNumber"`   // Counterparty account number (optional)
	VirtualAccountName     *string `json:"virtualAccountName"`     // Virtual account name (optional)
	VirtualAccountNumber   *string `json:"virtualAccountNumber"`   // Virtual account number (optional)
}

// WebhookRes represents the response for webhook processing.
type WebhookRes struct {
	Message string `json:"message"` // Response message (e.g., "Webhook processed")
}

// CreatePaymentLinkReq for creating payment link for orders
type CreatePaymentLinkReq struct {
	g.Meta    `path:"/create-payment-link" tags:"Payment" method:"POST" summary:"Create payment link for order" description:"Creates a PayOS payment link for parking or service orders" middleware:"middleware.Auth"`
	OrderType string `json:"order_type" v:"required|in:parking,service#Order type is required|Order type must be parking or service"`
	OrderID   int64  `json:"order_id" v:"required|min:1#Order ID is required|Order ID must be positive"`
}

// CreatePaymentLinkRes for payment link response
type CreatePaymentLinkRes struct {
	PaymentLinkId string `json:"paymentLinkId"` // Payment link ID
	CheckoutUrl   string `json:"checkoutUrl"`   // URL for payment checkout
	QRCode        string `json:"qrCode"`        // Base64 encoded QR code
	Amount        int    `json:"amount"`        // Payment amount in VND
	OrderCode     int64  `json:"orderCode"`     // Order code used in PayOS
}

// PaymentStatisticsGetReq for fetching payment statistics from PayOS
type PaymentStatisticsGetReq struct {
	g.Meta   `path:"/payment-statistics" tags:"Payment" method:"GET" summary:"Get payment statistics" description:"Fetches payment statistics from local database"`
	Page     int `json:"page" v:"min:0#Page must be non-negative" d:"0"`
	PageSize int `json:"pageSize" v:"min:1|max:100#Page size must be between 1 and 100" d:"50"`
}

// PaymentStatisticsGetRes for payment statistics response
type PaymentStatisticsGetRes struct {
	Code string                 `json:"code"` // Response code
	Desc string                 `json:"desc"` // Response description
	Data map[string]interface{} `json:"data"` // Statistics data (flexible structure)
}

// PaymentStatisticsData contains the payment orders and total count
type PaymentStatisticsData struct {
	Orders    []PaymentOrderItem `json:"orders"`    // List of payment orders
	TotalRows int                `json:"totalRows"` // Total number of orders
}

// PaymentOrderItem represents a single payment order from PayOS
type PaymentOrderItem struct {
	ID                 int64         `json:"id"`
	UUID               string        `json:"uuid"`
	OrderCode          int64         `json:"orderCode"`
	Amount             string        `json:"amount"`
	AmountPaid         *string       `json:"amountPaid"`
	AmountRemaining    string        `json:"amountRemaining"`
	Description        string        `json:"description"`
	AccountName        string        `json:"accountName"`
	AccountNumber      string        `json:"accountNumber"`
	Status             string        `json:"status"` // PENDING, PAID, CANCELLED
	Items              []PaymentItem `json:"items"`
	CancelUrl          string        `json:"cancelUrl"`
	ReturnUrl          string        `json:"returnUrl"`
	BuyerName          *string       `json:"buyerName"`
	BuyerEmail         *string       `json:"buyerEmail"`
	BuyerPhone         *string       `json:"buyerPhone"`
	Signature          string        `json:"signature"`
	CancelledAt        *string       `json:"cancelledAt"`
	CancellationReason *string       `json:"cancellationReason"`
	PaidAt             *string       `json:"paidAt"`
	CreatedAt          string        `json:"createdAt"`
	UpdatedAt          string        `json:"updatedAt"`
}

// PaymentItem represents an item in a payment order
type PaymentItem struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}
