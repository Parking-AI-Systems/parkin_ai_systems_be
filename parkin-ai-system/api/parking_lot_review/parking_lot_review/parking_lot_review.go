package parking_lot_review

import "github.com/gogf/gf/v2/frame/g"

type ParkingLotReviewAddReq struct {
	g.Meta  `path:"/parking-lot-reviews" tags:"Parking Lot Review" method:"POST" summary:"Add a new review" description:"Creates a new review for a parking lot. User must have a completed order for the lot." middleware:"middleware.Auth"`
	LotId   int64  `json:"lot_id" v:"required|min:1#Parking lot ID is required|Parking lot ID must be positive"`
	Rating  int    `json:"rating" v:"required|min:1|max:5#Rating is required|Rating must be between 1 and 5"`
	Comment string `json:"comment" v:"length:0,1000#Comment must be less than 1000 characters"`
}

type ParkingLotReviewAddRes struct {
	Id int64 `json:"id"`
}

type ParkingLotReviewListReq struct {
	g.Meta   `path:"/parking-lot-reviews" tags:"Parking Lot Review" method:"GET" summary:"List reviews" description:"Retrieves a paginated list of reviews for a parking lot with optional filters." middleware:"middleware.Auth"`
	LotId    int64 `json:"lot_id" v:"min:0#Parking lot ID must be non-negative"`
	Rating   int   `json:"rating" v:"min:0|max:5#Rating must be between 0 and 5"`
	Page     int   `json:"page" v:"min:1#Page must be at least 1"`
	PageSize int   `json:"page_size" v:"min:1|max:100#Page size must be between 1 and 100"`
}

type ParkingLotReviewItem struct {
	Id        int64  `json:"id"`
	LotId     int64  `json:"lot_id"`
	LotName   string `json:"lot_name"`
	UserId    int64  `json:"user_id"`
	Username  string `json:"username"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"created_at"`
}

type ParkingLotReviewListRes struct {
	List  []ParkingLotReviewItem `json:"list"`
	Total int                    `json:"total"`
}

type ParkingLotReviewGetReq struct {
	g.Meta `path:"/parking-lot-reviews/:id" tags:"Parking Lot Review" method:"GET" summary:"Get review details" description:"Retrieves details of a specific review by ID." middleware:"middleware.Auth"`
	Id     int64 `json:"id" v:"required|min:1#Review ID is required|Review ID must be positive"`
}

type ParkingLotReviewGetRes struct {
	Review ParkingLotReviewItem `json:"review"`
}

type ParkingLotReviewUpdateReq struct {
	g.Meta  `path:"/parking-lot-reviews/:id" tags:"Parking Lot Review" method:"PATCH" summary:"Update a review" description:"Updates the details of a review. Only the owner or admin can update." middleware:"middleware.Auth"`
	Id      int64  `json:"id" v:"required|min:1#Review ID is required|Review ID must be positive"`
	LotId   int64  `json:"lot_id" v:"min:0#Parking lot ID must be non-negative"`
	Rating  int    `json:"rating" v:"min:0|max:5#Rating must be between 0 and 5"`
	Comment string `json:"comment" v:"length:0,1000#Comment must be less than 1000 characters"`
}

type ParkingLotReviewUpdateRes struct {
	Review ParkingLotReviewItem `json:"review"`
}

type ParkingLotReviewDeleteReq struct {
	g.Meta `path:"/parking-lot-reviews/:id" tags:"Parking Lot Review" method:"DELETE" summary:"Delete a review" description:"Permanently deletes a review. Only the owner or admin can delete." middleware:"middleware.Auth"`
	Id     int64 `json:"id" v:"required|min:1#Review ID is required|Review ID must be positive"`
}

type ParkingLotReviewDeleteRes struct {
	Message string `json:"message"`
}

type MyParkingLotReviewListReq struct {
	g.Meta   `path:"/my-parking-lot-reviews" tags:"Parking Lot Review" method:"GET" summary:"List my reviews" description:"Retrieves a paginated list of reviews submitted by the authenticated user with optional filters." middleware:"middleware.Auth"`
	LotId    int64 `json:"lot_id" v:"min:0#Parking lot ID must be non-negative"`
	Page     int   `json:"page" v:"min:1#Page must be at least 1"`
	PageSize int   `json:"page_size" v:"min:1|max:100#Page size must be between 1 and 100"`
}
type MyParkingLotReviewListRes struct {
	List  []ParkingLotReviewItem `json:"list"`
	Total int                    `json:"total"`
}
