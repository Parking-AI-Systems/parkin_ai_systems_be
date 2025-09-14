// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ParkingLotReviews is the golang structure for table parking_lot_reviews.
type ParkingLotReviews struct {
	Id        int64       `json:"id"        orm:"id"         description:""`
	LotId     int64       `json:"lotId"     orm:"lot_id"     description:""`
	UserId    int64       `json:"userId"    orm:"user_id"    description:""`
	Rating    int         `json:"rating"    orm:"rating"     description:""`
	Comment   string      `json:"comment"   orm:"comment"    description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at" description:""`
}

type ParkingLotReviewAddReq struct {
	LotId   int64  `json:"lotId"`
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

type ParkingLotReviewAddRes struct {
	Id int64 `json:"id"`
}

type ParkingLotReviewListReq struct {
	LotId    int64 `json:"lotId"`
	Rating   int   `json:"rating"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
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
	UpdatedAt string `json:"updated_at"`
}

type ParkingLotReviewListRes struct {
	List  []ParkingLotReviewItem `json:"list"`
	Total int                    `json:"total"`
}

type ParkingLotReviewGetReq struct {
	Id int64 `json:"id"`
}

type ParkingLotReviewGetRes struct {
	Review ParkingLotReviewItem `json:"review"`
}

type ParkingLotReviewUpdateReq struct {
	Id      int64  `json:"id"`
	LotId   int64  `json:"lotId"`
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

type ParkingLotReviewUpdateRes struct {
	Review ParkingLotReviewItem `json:"review"`
}

type ParkingLotReviewDeleteReq struct {
	Id int64 `json:"id"`
}

type ParkingLotReviewDeleteRes struct {
	Message string `json:"message"`
}
type MyParkingLotReviewReq struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}
type MyParkingLotReviewRes struct {
	List  []ParkingLotReviewItem `json:"list"`
	Total int                    `json:"total"`
}