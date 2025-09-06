// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"parkin-ai-system/internal/model/entity"
)

type (
	IParkingLotReview interface {
		ParkingLotReviewAdd(ctx context.Context, req *entity.ParkingLotReviewAddReq) (*entity.ParkingLotReviewAddRes, error)
		ParkingLotReviewList(ctx context.Context, req *entity.ParkingLotReviewListReq) (*entity.ParkingLotReviewListRes, error)
		ParkingLotReviewGet(ctx context.Context, req *entity.ParkingLotReviewGetReq) (*entity.ParkingLotReviewItem, error)
		ParkingLotReviewUpdate(ctx context.Context, req *entity.ParkingLotReviewUpdateReq) (*entity.ParkingLotReviewItem, error)
		ParkingLotReviewDelete(ctx context.Context, req *entity.ParkingLotReviewDeleteReq) (*entity.ParkingLotReviewDeleteRes, error)
	}
)

var (
	localParkingLotReview IParkingLotReview
)

func ParkingLotReview() IParkingLotReview {
	if localParkingLotReview == nil {
		panic("implement not found for interface IParkingLotReview, forgot register?")
	}
	return localParkingLotReview
}

func RegisterParkingLotReview(i IParkingLotReview) {
	localParkingLotReview = i
}
