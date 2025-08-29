package parking_lot_review

import (
	"github.com/gogf/gf/v2/frame/g"
)

type ParkingLotReviewAddReq struct {
	g.Meta   `path:"/parking-lot-reviews" method:"post" tags:"ParkingLotReview" summary:"Thêm đánh giá bãi xe" security:"BearerAuth"`
	LotId    int64  `json:"lot_id" v:"required"`
	Rating   int    `json:"rating" v:"required|min:1|max:5"`
	Comment  string `json:"comment"`
}

type ParkingLotReviewAddRes struct {
	ReviewId int64 `json:"review_id"`
}

type ParkingLotReviewUpdateReq struct {
	g.Meta   `path:"/parking-lot-reviews/{id}" method:"put" tags:"ParkingLotReview" summary:"Cập nhật đánh giá bãi xe" security:"BearerAuth"`
	Id       int64  `json:"id" v:"required"`
	Rating   int    `json:"rating" v:"min:1|max:5"`
	Comment  string `json:"comment"`
}

type ParkingLotReviewUpdateRes struct {
	Success bool `json:"success"`
}

type ParkingLotReviewDeleteReq struct {
	g.Meta   `path:"/parking-lot-reviews/{id}" method:"delete" tags:"ParkingLotReview" summary:"Xoá đánh giá bãi xe" security:"BearerAuth"`
	Id       int64  `json:"id" v:"required"`
}

type ParkingLotReviewDeleteRes struct {
	Success bool `json:"success"`
}

type ParkingLotReviewDetailReq struct {
	g.Meta   `path:"/parking-lot-reviews/{id}" method:"get" tags:"ParkingLotReview" summary:"Chi tiết đánh giá bãi xe" security:"BearerAuth"`
	Id       int64  `json:"id" v:"required"`
}

type ParkingLotReviewDetailRes struct {
	Review *ParkingLotReviewInfo `json:"review"`
}

type ParkingLotReviewListReq struct {
	g.Meta   `path:"/parking-lot-reviews" method:"get" tags:"ParkingLotReview" summary:"Danh sách đánh giá bãi xe" security:"BearerAuth"`
	LotId    int64  `json:"lot_id"`
}

type ParkingLotReviewListRes struct {
	Reviews []ParkingLotReviewInfo `json:"reviews"`
}

type ParkingLotReviewInfo struct {
	Id        int64  `json:"id"`
	LotId     int64  `json:"lot_id"`
	UserId    int64  `json:"user_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"created_at"`
}
