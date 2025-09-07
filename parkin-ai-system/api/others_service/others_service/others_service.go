package others_service

import "github.com/gogf/gf/v2/frame/g"

type OthersServiceAddReq struct {
	g.Meta          `path:"/others-services" tags:"Others Service" method:"POST" summary:"Add a new service" description:"Creates a new service for a parking lot. Admin only." middleware:"middleware.Auth"`
	LotId           int64   `json:"lot_id" v:"required|min:1#Parking lot ID is required|Parking lot ID must be positive"`
	Name            string  `json:"name" v:"required#Service name is required"`
	Description     string  `json:"description" v:"length:0,1000#Description must be less than 1000 characters"`
	Price           float64 `json:"price" v:"required|min:0#Price is required|Price must be non-negative"`
	DurationMinutes int     `json:"duration_minutes" v:"required|min:1#Duration is required|Duration must be positive"`
	IsActive        bool    `json:"is_active"`
}

type OthersServiceAddRes struct {
	Id int64 `json:"id"`
}

type OthersServiceListReq struct {
	g.Meta   `path:"/others-services" tags:"Others Service" method:"GET" summary:"List services" description:"Retrieves a paginated list of services with optional filters for parking lot and active status." middleware:"middleware.Auth"`
	LotId    int64 `json:"lot_id" v:"min:0#Parking lot ID must be non-negative"`
	IsActive bool  `json:"is_active"`
	Page     int   `json:"page" v:"min:1#Page must be at least 1"`
	PageSize int   `json:"page_size" v:"min:1|max:100#Page size must be between 1 and 100"`
}

type OthersServiceItem struct {
	Id              int64   `json:"id"`
	LotId           int64   `json:"lot_id"`
	LotName         string  `json:"lot_name"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	DurationMinutes int     `json:"duration_minutes"`
	IsActive        bool    `json:"is_active"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
	DeletedAt       string  `json:"deleted_at"`
}

type OthersServiceListRes struct {
	List  []OthersServiceItem `json:"list"`
	Total int                 `json:"total"`
}

type OthersServiceGetReq struct {
	g.Meta `path:"/others-services/:id" tags:"Others Service" method:"GET" summary:"Get service details" description:"Retrieves details of a specific service by ID." middleware:"middleware.Auth"`
	Id     int64 `json:"id" v:"required|min:1#Service ID is required|Service ID must be positive"`
}

type OthersServiceGetRes struct {
	Service OthersServiceItem `json:"service"`
}

type OthersServiceUpdateReq struct {
	g.Meta          `path:"/others-services/:id" tags:"Others Service" method:"PATCH" summary:"Update a service" description:"Updates the details of a service. Admin only." middleware:"middleware.Auth"`
	Id              int64   `json:"id" v:"required|min:1#Service ID is required|Service ID must be positive"`
	LotId           int64   `json:"lot_id" v:"min:0#Parking lot ID must be non-negative"`
	Name            string  `json:"name" v:"length:0,255#Name must be less than 255 characters"`
	Description     string  `json:"description" v:"length:0,1000#Description must be less than 1000 characters"`
	Price           float64 `json:"price" v:"min:-1#Price must be non-negative"`
	DurationMinutes int     `json:"duration_minutes" v:"min:-1#Duration must be non-negative"`
	IsActive        *bool   `json:"is_active"`
}

type OthersServiceUpdateRes struct {
	Service OthersServiceItem `json:"service"`
}

type OthersServiceDeleteReq struct {
	g.Meta `path:"/others-services/:id" tags:"Others Service" method:"DELETE" summary:"Delete a service" description:"Permanently deletes a service. Admin only." middleware:"middleware.Auth"`
	Id     int64 `json:"id" v:"required|min:1#Service ID is required|Service ID must be positive"`
}

type OthersServiceDeleteRes struct {
	Message string `json:"message"`
}
