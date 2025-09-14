package parking_order

import "github.com/gogf/gf/v2/frame/g"

// ParkingOrderAddReq defines the request structure for creating a new parking order.
type ParkingOrderAddReq struct {
	g.Meta    `path:"/parking-orders" tags:"Parking Order" method:"POST" summary:"Create a new parking order" description:"Creates a parking order for a user with the specified vehicle, lot, slot, and time period." middleware:"middleware.Auth"`
	VehicleId int64  `json:"vehicle_id" v:"required|min:1#Vehicle ID is required|Vehicle ID must be positive"`        // ID of the vehicle
	LotId     int64  `json:"lot_id" v:"required|min:1#Parking lot ID is required|Parking lot ID must be positive"`    // ID of the parking lot
	SlotId    int64  `json:"slot_id" v:"required|min:1#Parking slot ID is required|Parking slot ID must be positive"` // ID of the parking slot
	StartTime string `json:"start_time" v:"required|date#Start time is required|Invalid start time format"`           // Start time in string format (e.g., "2025-09-06 10:00:00")
	EndTime   string `json:"end_time" v:"required|date#End time is required|Invalid end time format"`                 // End time in string format (e.g., "2025-09-06 12:00:00")
}

// ParkingOrderAddRes defines the response structure for creating a new parking order.
type ParkingOrderAddRes struct {
	Id int64 `json:"id"` // ID of the created parking order
}

// ParkingOrderListReq defines the request structure for listing parking orders.
type ParkingOrderListReq struct {
	g.Meta   `path:"/parking-orders" tags:"Parking Order" method:"GET" summary:"List parking orders" description:"Retrieves a paginated list of parking orders with optional filters for user, lot, and status." middleware:"middleware.Auth"`
	UserId   int64  `json:"user_id" v:"min:0#User ID must be non-negative"`                          // Filter by user ID (optional, 0 means no filter)
	LotId    int64  `json:"lot_id" v:"min:0#Parking lot ID must be non-negative"`                    // Filter by parking lot ID (optional, 0 means no filter)
	Page     int    `json:"page" v:"min:1#Page must be at least 1"`                                  // Page number for pagination
	PageSize int    `json:"page_size" v:"min:1|max:100#Page size must be between 1 and 100"`         // Number of items per page
	Status   string `json:"status" v:"in:pending,confirmed,canceled,completed#Invalid status value"` // Filter by order status (optional)
}

// ParkingOrderItem defines the structure for an individual parking order in the response.
type ParkingOrderItem struct {
	Id            int64   `json:"id"`             // Order ID
	UserId        int64   `json:"user_id"`        // User ID
	LotId         int64   `json:"lot_id"`         // Parking lot ID
	SlotId        int64   `json:"slot_id"`        // Parking slot ID
	VehicleId     int64   `json:"vehicle_id"`     // Vehicle ID
	LotName       string  `json:"lot_name"`       // Parking lot name (from join)
	SlotCode      string  `json:"slot_code"`      // Parking slot code (from join)
	VehiclePlate  string  `json:"vehicle_plate"`  // Vehicle license plate (from join)
	StartTime     string  `json:"start_time"`     // Start time (formatted)
	EndTime       string  `json:"end_time"`       // End time (formatted)
	Status        string  `json:"status"`         // Order status
	Price         float64 `json:"price"`          // Order price
	PaymentStatus string  `json:"payment_status"` // Payment status
	CreatedAt     string  `json:"created_at"`     // Creation timestamp (formatted)
	UpdatedAt     string  `json:"updated_at"`     // Update timestamp (formatted, optional)
}

// ParkingOrderListRes defines the response structure for listing parking orders.
type ParkingOrderListRes struct {
	List  []ParkingOrderItem `json:"list"`  // List of parking orders
	Total int                `json:"total"` // Total number of matching orders
}

// ParkingOrderGetReq defines the request structure for retrieving a parking order.
type ParkingOrderGetReq struct {
	g.Meta `path:"/parking-orders/:id" tags:"Parking Order" method:"GET" summary:"Get parking order details" description:"Retrieves details of a specific parking order by ID." middleware:"middleware.Auth"`
	Id     int64 `json:"id" v:"required|min:1#Order ID is required|Order ID must be positive"` // Parking order ID
}

// ParkingOrderGetRes defines the response structure for retrieving a parking order.
type ParkingOrderGetRes struct {
	Order ParkingOrderItem `json:"order"` // Details of the parking order
}

// ParkingOrderUpdateReq defines the request structure for updating a parking order.
type ParkingOrderUpdateReq struct {
	g.Meta    `path:"/parking-orders/:id" tags:"Parking Order" method:"PATCH" summary:"Update a parking order" description:"Updates the start time, end time, or status of a parking order." middleware:"middleware.Auth"`
	Id        int64  `json:"id" v:"required|min:1#Order ID is required|Order ID must be positive"`    // Parking order ID
	StartTime string `json:"start_time" v:"date#Invalid start time format"`                           // New start time (optional)
	EndTime   string `json:"end_time" v:"date#Invalid end time format"`                               // New end time (optional)
	Status    string `json:"status" v:"in:pending,confirmed,canceled,completed#Invalid status value"` // New status (optional)
}

// ParkingOrderUpdateRes defines the response structure for updating a parking order.
type ParkingOrderUpdateRes struct {
	Order ParkingOrderItem `json:"order"` // Updated parking order details
}

// ParkingOrderCancelReq defines the request structure for canceling a parking order.
type ParkingOrderCancelReq struct {
	g.Meta `path:"/parking-orders/:id/cancel" tags:"Parking Order" method:"PATCH" summary:"Cancel a parking order" description:"Cancels a parking order and updates slot availability." middleware:"middleware.Auth"`
	Id     int64 `json:"id" v:"required|min:1#Order ID is required|Order ID must be positive"` // Parking order ID
}

// ParkingOrderCancelRes defines the response structure for canceling a parking order.
type ParkingOrderCancelRes struct {
	Order ParkingOrderItem `json:"order"` // Canceled parking order details
}

// ParkingOrderDeleteReq defines the request structure for deleting a parking order.
type ParkingOrderDeleteReq struct {
	g.Meta `path:"/parking-orders/:id" tags:"Parking Order" method:"DELETE" summary:"Delete a parking order" description:"Soft deletes a parking order and updates slot availability." middleware:"middleware.Auth"`
	Id     int64 `json:"id" v:"required|min:1#Order ID is required|Order ID must be positive"` // Parking order ID
}

// ParkingOrderDeleteRes defines the response structure for deleting a parking order.
type ParkingOrderDeleteRes struct {
	Message string `json:"message"` // Confirmation message
}

// ParkingOrderPaymentReq defines the request structure for confirming payment of a parking order.
type ParkingOrderPaymentReq struct {
	g.Meta `path:"/parking-orders/:id/payment" tags:"Parking Order" method:"PATCH" summary:"Confirm payment for a parking order" description:"Updates the payment status to paid and processes wallet transaction." middleware:"middleware.Auth"`
	Id     int64 `json:"id" v:"required|min:1#Order ID is required|Order ID must be positive"` // Parking order ID
}

// ParkingOrderPaymentRes defines the response structure for confirming payment of a parking order.
type ParkingOrderPaymentRes struct {
	Order ParkingOrderItem `json:"order"` // Paid parking order details
}

// New dashboard APIs

// ParkingOrderRevenueReq defines the request structure for retrieving total revenue from parking orders.
type ParkingOrderRevenueReq struct {
	g.Meta    `path:"/parking-orders/revenue" tags:"Parking Order" method:"GET" summary:"Get total parking order revenue" description:"Retrieves total revenue from parking orders filtered by period." middleware:"middleware.Auth"`
	Period    string `json:"period" v:"in:1h,1d,1w,1m,custom#Invalid period value"`                                                                // Period filter (1h, 1d, 1w, 1m, custom)
	StartTime string `json:"start_time" v:"datetime|required-if:Period,custom#Invalid start time format|Start time is required for custom period"` // Start time (YYYY-MM-DD HH:MM:SS) for custom period
	EndTime   string `json:"end_time" v:"datetime|required-if:Period,custom#Invalid end time format|End time is required for custom period"`       // End time (YYYY-MM-DD HH:MM:SS) for custom period
}

// ParkingOrderRevenueRes defines the response structure for total revenue.
type ParkingOrderRevenueRes struct {
	TotalRevenue float64 `json:"total_revenue"` // Total revenue from parking orders
}

// ParkingOrderTrendsReq defines the request structure for retrieving parking order trends.
type ParkingOrderTrendsReq struct {
	g.Meta    `path:"/parking-orders/trends" tags:"Parking Order" method:"GET" summary:"Get parking order trends" description:"Retrieves parking order trends over time filtered by period, suitable for line charts." middleware:"middleware.Auth"`
	Period    string `json:"period" v:"in:1h,1d,1w,1m,custom#Invalid period value"`                                                                // Period filter (1h, 1d, 1w, 1m, custom)
	StartTime string `json:"start_time" v:"datetime|required-if:Period,custom#Invalid start time format|Start time is required for custom period"` // Start time (YYYY-MM-DD HH:MM:SS) for custom period
	EndTime   string `json:"end_time" v:"datetime|required-if:Period,custom#Invalid end time format|End time is required for custom period"`       // End time (YYYY-MM-DD HH:MM:SS) for custom period
}

// ParkingOrderTrendsRes defines the response structure for parking order trends.
type ParkingOrderTrendsRes struct {
	Orders []ParkingOrderTrendsItem `json:"orders"` // List of trend data points
	Total  int64                    `json:"total"`  // Total number of orders
}

// ParkingOrderTrendsItem defines an individual trend data point.
type ParkingOrderTrendsItem struct {
	Date  string `json:"date"`  // Date (formatted based on period)
	Count int64  `json:"count"` // Number of orders for the date
}

// ParkingOrderStatusBreakdownReq defines the request structure for retrieving parking order status breakdown.
type ParkingOrderStatusBreakdownReq struct {
	g.Meta    `path:"/parking-orders/status" tags:"Parking Order" method:"GET" summary:"Get parking order status breakdown" description:"Retrieves the distribution of parking order statuses filtered by period." middleware:"middleware.Auth"`
	Period    string `json:"period" v:"in:1h,1d,1w,1m,custom#Invalid period value"`                                                                // Period filter (1h, 1d, 1w, 1m, custom)
	StartTime string `json:"start_time" v:"datetime|required-if:Period,custom#Invalid start time format|Start time is required for custom period"` // Start time (YYYY-MM-DD HH:MM:SS) for custom period
	EndTime   string `json:"end_time" v:"datetime|required-if:Period,custom#Invalid end time format|End time is required for custom period"`       // End time (YYYY-MM-DD HH:MM:SS) for custom period
}

// ParkingOrderStatusBreakdownRes defines the response structure for status breakdown.
type ParkingOrderStatusBreakdownRes struct {
	Statuses []ParkingOrderStatusItem `json:"statuses"` // List of status counts
	Total    int64                    `json:"total"`    // Total number of orders
}

// ParkingOrderStatusItem defines an individual status count.
type ParkingOrderStatusItem struct {
	Status string `json:"status"` // Order status (pending, confirmed, canceled, completed)
	Count  int64  `json:"count"`  // Number of orders with this status
}

type MyParkingLotOrderListReq struct {
	g.Meta   `path:"/my/parking-orders" tags:"My Parking Orders" method:"GET" summary:"List my parking orders" description:"Retrieves a paginated list of parking orders for the authenticated user with optional filters for lot and status." middleware:"middleware.Auth"`
	LotId    int64  `json:"lot_id" v:"min:0#Parking lot ID must be non-negative"`                    // Filter by parking lot ID (optional, 0 means no filter)
	Page     int    `json:"page" v:"min:1#Page must be at least 1"`                                  // Page number for pagination
	PageSize int    `json:"page_size" v:"min:1|max:100#Page size must be between 1 and 100"`         // Number of items per page
	Status   string `json:"status" v:"in:pending,confirmed,canceled,completed#Invalid status value"` // Filter by order status (optional)

}
type MyParkingLotOrderListRes struct {
	List  []ParkingOrderItem `json:"list"`  // List of parking orders
	Total int                `json:"total"` // Total number of matching orders
}
