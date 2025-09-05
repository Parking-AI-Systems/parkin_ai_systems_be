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
