package entity

// PayOS payment types
type Item struct {
	Name  string `json:"name"`  // Tên sản phẩm
	Price int    `json:"price"` // Giá sản phẩm
}

type CheckoutRequestType struct {
	OrderCode    int64   `json:"orderCode"`    // Mã đơn hàng
	Amount       int     `json:"amount"`       // Tổng tiền đơn hàng
	Description  string  `json:"description"`  // Mô tả đơn hàng
	CancelUrl    string  `json:"cancelUrl"`    // URL hủy thanh toán
	ReturnUrl    string  `json:"returnUrl"`    // URL thành công
	Items        []Item  `json:"items"`        // Danh sách sản phẩm
	BuyerName    *string `json:"buyerName"`    // Tên người mua
	BuyerEmail   *string `json:"buyerEmail"`   // Email người mua
	BuyerPhone   *string `json:"buyerPhone"`   // Số điện thoại người mua
	BuyerAddress *string `json:"buyerAddress"` // Địa chỉ người mua
	ExpiredAt    *int64  `json:"expiredAt"`    // Thời gian hết hạn
	Signature    string  `json:"signature"`    // Chữ ký
}

type CheckoutResponseType struct {
	Code string                   `json:"code"`
	Desc string                   `json:"desc"`
	Data CheckoutResponseDataType `json:"data"`
}

type CheckoutResponseDataType struct {
	Bin           string `json:"bin"`           // Mã BIN ngân hàng
	AccountNumber string `json:"accountNumber"` // Số tài khoản
	AccountName   string `json:"accountName"`   // Tên chủ tài khoản
	Amount        int    `json:"amount"`        // Tổng tiền
	Description   string `json:"description"`   // Mô tả
	OrderCode     int64  `json:"orderCode"`     // Mã đơn hàng
	Currency      string `json:"currency"`      // Đơn vị tiền tệ
	PaymentLinkId string `json:"paymentLinkId"` // Mã link thanh toán
	Status        string `json:"status"`        // Trạng thái
	CheckoutUrl   string `json:"checkoutUrl"`   // URL thanh toán
	QRCode        string `json:"qrCode"`        // Mã QR
}

// Webhook structures
type WebhookType struct {
	Code      string           `json:"code"`      // Mã lỗi
	Desc      string           `json:"desc"`      // Mô tả lỗi
	Success   bool             `json:"success"`   // Trạng thái của webhook
	Data      *WebhookDataType `json:"data"`      // Dữ liệu webhook
	Signature string           `json:"signature"` // Chữ ký số
}

type WebhookDataType struct {
	OrderCode              int64   `json:"orderCode"`              // Mã đơn hàng
	Amount                 int     `json:"amount"`                 // Số tiền
	Description            string  `json:"description"`            // Mô tả
	AccountNumber          string  `json:"accountNumber"`          // Số tài khoản
	Reference              string  `json:"reference"`              // Mã tham chiếu
	TransactionDateTime    string  `json:"transactionDateTime"`    // Thời gian giao dịch
	Currency               string  `json:"currency"`               // Đơn vị tiền tệ
	PaymentLinkId          string  `json:"paymentLinkId"`          // Mã link thanh toán
	Code                   string  `json:"code"`                   // Mã lỗi
	Desc                   string  `json:"desc"`                   // Mô tả lỗi
	CounterAccountBankId   *string `json:"counterAccountBankId"`   // Mã ngân hàng đối ứng
	CounterAccountBankName *string `json:"counterAccountBankName"` // Tên ngân hàng đối ứng
	CounterAccountName     *string `json:"counterAccountName"`     // Tên chủ tài khoản đối ứng
	CounterAccountNumber   *string `json:"counterAccountNumber"`   // Số tài khoản đối ứng
	VirtualAccountName     *string `json:"virtualAccountName"`     // Tên tài khoản ảo
	VirtualAccountNumber   *string `json:"virtualAccountNumber"`   // Số tài khoản ảo
}

// PaymentLinkDataType represents payment link details from PayOS
type PaymentLinkDataType struct {
	Id                 string            `json:"id"`                 // Payment link ID
	OrderCode          int64             `json:"orderCode"`          // Order code
	Amount             int               `json:"amount"`             // Total amount
	AmountPaid         int               `json:"amountPaid"`         // Paid amount
	AmountRemaining    int               `json:"amountRemaining"`    // Remaining amount
	Status             string            `json:"status"`             // Payment status
	CreatedAt          string            `json:"createdAt"`          // Creation time
	Transactions       []TransactionItem `json:"transactions"`       // Transaction list
	CancellationReason *string           `json:"cancellationReason"` // Cancellation reason
	CancelledAt        *string           `json:"cancelledAt"`        // Cancellation time
}

// TransactionItem represents a transaction in payment link
type TransactionItem struct {
	Reference              string  `json:"reference"`              // Transaction reference
	Amount                 int     `json:"amount"`                 // Transaction amount
	AccountNumber          string  `json:"accountNumber"`          // Account number
	Description            string  `json:"description"`            // Description
	TransactionDateTime    string  `json:"transactionDateTime"`    // Transaction time
	VirtualAccountName     *string `json:"virtualAccountName"`     // Virtual account name
	VirtualAccountNumber   *string `json:"virtualAccountNumber"`   // Virtual account number
	CounterAccountBankId   *string `json:"counterAccountBankId"`   // Counter bank ID
	CounterAccountBankName *string `json:"counterAccountBankName"` // Counter bank name
	CounterAccountName     *string `json:"counterAccountName"`     // Counter account name
	CounterAccountNumber   *string `json:"counterAccountNumber"`   // Counter account number
}

// RefundResponseType represents refund response from PayOS
type RefundResponseType struct {
	Code string         `json:"code"`
	Desc string         `json:"desc"`
	Data RefundDataType `json:"data"`
}

// RefundDataType represents refund data
type RefundDataType struct {
	RefundId string `json:"refundId"` // Refund transaction ID
	Status   string `json:"status"`   // Refund status
}
