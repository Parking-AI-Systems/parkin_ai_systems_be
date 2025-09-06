package vehicle

import "github.com/gogf/gf/v2/frame/g"

type VehicleAddReq struct {
	g.Meta       `path:"/vehicles" tags:"Vehicle" method:"POST" summary:"Add a new vehicle" description:"Creates a new vehicle for the authenticated user." middleware:"middleware.Auth"`
	LicensePlate string `json:"license_plate" v:"required|length:1,20#License plate is required|License plate must be between 1 and 20 characters"`
	Brand        string `json:"brand" v:"length:0,50#Brand must be less than 50 characters"`
	Model        string `json:"model" v:"length:0,50#Model must be less than 50 characters"`
	Color        string `json:"color" v:"length:0,50#Color must be less than 50 characters"`
	Type         string `json:"type" v:"required#Vehicle type is required"`
}

type VehicleAddRes struct {
	Id int64 `json:"id"`
}

type VehicleListReq struct {
	g.Meta   `path:"/vehicles" tags:"Vehicle" method:"GET" summary:"List vehicles" description:"Retrieves a paginated list of vehicles for the authenticated user or all vehicles for admins." middleware:"middleware.Auth"`
	Type     string `json:"type"`
	Page     int    `json:"page" v:"min:1#Page must be at least 1"`
	PageSize int    `json:"page_size" v:"min:1|max:100#Page size must be between 1 and 100"`
}

type VehicleItem struct {
	Id           int64  `json:"id"`
	UserId       int64  `json:"user_id"`
	Username     string `json:"username"`
	LicensePlate string `json:"license_plate"`
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	Color        string `json:"color"`
	Type         string `json:"type"`
	CreatedAt    string `json:"created_at"`
}

type VehicleListRes struct {
	List  []VehicleItem `json:"list"`
	Total int           `json:"total"`
}

type VehicleGetReq struct {
	g.Meta `path:"/vehicles/:id" tags:"Vehicle" method:"GET" summary:"Get vehicle details" description:"Retrieves details of a specific vehicle by ID." middleware:"middleware.Auth"`
	Id     int64 `json:"id" v:"required|min:1#Vehicle ID is required|Vehicle ID must be positive"`
}

type VehicleGetRes struct {
	Vehicle VehicleItem `json:"vehicle"`
}

type VehicleUpdateReq struct {
	g.Meta       `path:"/vehicles/:id" tags:"Vehicle" method:"PATCH" summary:"Update a vehicle" description:"Updates the details of a vehicle. Only the owner or admin can update." middleware:"middleware.Auth"`
	Id           int64  `json:"id" v:"required|min:1#Vehicle ID is required|Vehicle ID must be positive"`
	LicensePlate string `json:"license_plate" v:"length:0,20#License plate must be less than 20 characters"`
	Brand        string `json:"brand" v:"length:0,50#Brand must be less than 50 characters"`
	Model        string `json:"model" v:"length:0,50#Model must be less than 50 characters"`
	Color        string `json:"color" v:"length:0,50#Color must be less than 50 characters"`
	Type         string `json:"type"`
}

type VehicleUpdateRes struct {
	Vehicle VehicleItem `json:"vehicle"`
}

type VehicleDeleteReq struct {
	g.Meta `path:"/vehicles/:id" tags:"Vehicle" method:"DELETE" summary:"Delete a vehicle" description:"Permanently deletes a vehicle. Only the owner or admin can delete." middleware:"middleware.Auth"`
	Id     int64 `json:"id" v:"required|min:1#Vehicle ID is required|Vehicle ID must be positive"`
}

type VehicleDeleteRes struct {
	Message string `json:"message"`
}
