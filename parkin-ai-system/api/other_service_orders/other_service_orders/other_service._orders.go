package other_service_orders

import "github.com/gogf/gf/v2/frame/g"

type OthersServiceOrderAddReq struct {
	g.Meta        `path:"/service-orders" tags:"Service Order" method:"POST" summary:"Create a new service order" description:"Creates a service order for a user with the specified vehicle, lot, service, and scheduled time." middleware:"middleware.Auth"`
	VehicleId     int64  `json:"vehicle_id" v:"required|min:1#Vehicle ID is required|Vehicle ID must be positive"`
	LotId         int64  `json:"lot_id" v:"required|min:1#Parking lot ID is required|Parking lot ID must be positive"`
	ServiceId     int64  `json:"service_id" v:"required|min:1#Service ID is required|Service ID must be positive"`
	ScheduledTime string `json:"scheduled_time" v:"required|date#Scheduled time is required|Invalid scheduled time format"`
}

type OthersServiceOrderAddRes struct {
	Id int64 `json:"id"`
}

type OthersServiceOrderListReq struct {
	g.Meta   `path:"/service-orders" tags:"Service Order" method:"GET" summary:"List service orders" description:"Retrieves a paginated list of service orders with optional filters for user, lot, and status." middleware:"middleware.Auth"`
	UserId   int64  `json:"user_id" v:"min:0#User ID must be non-negative"`
	LotId    int64  `json:"lot_id" v:"min:0#Parking lot ID must be non-negative"`
	Page     int    `json:"page" v:"min:1#Page must be at least 1"`
	PageSize int    `json:"page_size" v:"min:1|max:100#Page size must be between 1 and 100"`
	Status   string `json:"status" v:"in:pending,confirmed,canceled,completed#Invalid status value"`
}

type OthersServiceOrderItem struct {
	Id            int64   `json:"id"`
	UserId        int64   `json:"user_id"`
	LotId         int64   `json:"lot_id"`
	ServiceId     int64   `json:"service_id"`
	VehicleId     int64   `json:"vehicle_id"`
	LotName       string  `json:"lot_name"`
	ServiceName   string  `json:"service_name"`
	VehiclePlate  string  `json:"vehicle_plate"`
	ScheduledTime string  `json:"scheduled_time"`
	Status        string  `json:"status"`
	Price         float64 `json:"price"`
	PaymentStatus string  `json:"payment_status"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type OthersServiceOrderListRes struct {
	List  []OthersServiceOrderItem `json:"list"`
	Total int                      `json:"total"`
}

type OthersServiceOrderGetReq struct {
	g.Meta `path:"/service-orders/:id" tags:"Service Order" method:"GET" summary:"Get service order details" description:"Retrieves details of a specific service order by ID." middleware:"middleware.Auth"`
	Id     int64 `json:"id" v:"required|min:1#Order ID is required|Order ID must be positive"`
}

type OthersServiceOrderGetRes struct {
	Order OthersServiceOrderItem `json:"order"`
}

type OthersServiceOrderUpdateReq struct {
	g.Meta        `path:"/service-orders/:id" tags:"Service Order" method:"PATCH" summary:"Update a service order" description:"Updates the scheduled time or status of a service order." middleware:"middleware.Auth"`
	Id            int64  `json:"id" v:"required|min:1#Order ID is required|Order ID must be positive"`
	ScheduledTime string `json:"scheduled_time" v:"date#Invalid scheduled time format"`
	Status        string `json:"status" v:"in:pending,confirmed,canceled,completed#Invalid status value"`
}

type OthersServiceOrderUpdateRes struct {
	Order OthersServiceOrderItem `json:"order"`
}

type OthersServiceOrderCancelReq struct {
	g.Meta `path:"/service-orders/:id/cancel" tags:"Service Order" method:"PATCH" summary:"Cancel a service order" description:"Cancels a service order." middleware:"middleware.Auth"`
	Id     int64 `json:"id" v:"required|min:1#Order ID is required|Order ID must be positive"`
}

type OthersServiceOrderCancelRes struct {
	Order OthersServiceOrderItem `json:"order"`
}

type OthersServiceOrderDeleteReq struct {
	g.Meta `path:"/service-orders/:id" tags:"Service Order" method:"DELETE" summary:"Delete a service order" description:"Soft deletes a service order." middleware:"middleware.Auth"`
	Id     int64 `json:"id" v:"required|min:1#Order ID is required|Order ID must be positive"`
}

type OthersServiceOrderDeleteRes struct {
	Message string `json:"message"`
}

type OthersServiceOrderPaymentReq struct {
	g.Meta `path:"/service-orders/:id/payment" tags:"Service Order" method:"PATCH" summary:"Confirm payment for a service order" description:"Updates the payment status to paid and processes wallet transaction." middleware:"middleware.Auth"`
	Id     int64 `json:"id" v:"required|min:1#Order ID is required|Order ID must be positive"`
}

type OthersServiceOrderPaymentRes struct {
	Order OthersServiceOrderItem `json:"order"`
}
