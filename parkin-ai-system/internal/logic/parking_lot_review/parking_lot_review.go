package parking_lot_review

import (
	"context"
	"parkin-ai-system/api/parking_lot_review"
	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/model/do"
	"github.com/gogf/gf/v2/os/gtime"
)

type sParkingLotReview struct{}

func (s *sParkingLotReview) ParkingLotReviewAdd(ctx context.Context, req *parking_lot_review.ParkingLotReviewAddReq) (res *parking_lot_review.ParkingLotReviewAddRes, err error) {
	review := do.ParkingLotReviews{
		LotId:   req.LotId,
		UserId:  ctx.Value("user_id"), // cần lấy user_id từ context thực tế
		Rating:  req.Rating,
		Comment: req.Comment,
		CreatedAt: gtime.Now(),
	}
	result, err := dao.ParkingLotReviews.Ctx(ctx).Data(review).InsertAndGetId()
	if err != nil {
		return nil, err
	}
	res = &parking_lot_review.ParkingLotReviewAddRes{ReviewId: result.(int64)}
	return
}

func (s *sParkingLotReview) ParkingLotReviewUpdate(ctx context.Context, req *parking_lot_review.ParkingLotReviewUpdateReq) (res *parking_lot_review.ParkingLotReviewUpdateRes, err error) {
	data := do.ParkingLotReviews{
		Rating:  req.Rating,
		Comment: req.Comment,
	}
	_, err = dao.ParkingLotReviews.Ctx(ctx).Where("id", req.Id).Data(data).Update()
	if err != nil {
		return nil, err
	}
	res = &parking_lot_review.ParkingLotReviewUpdateRes{Success: true}
	return
}

func (s *sParkingLotReview) ParkingLotReviewDelete(ctx context.Context, req *parking_lot_review.ParkingLotReviewDeleteReq) (res *parking_lot_review.ParkingLotReviewDeleteRes, err error) {
	_, err = dao.ParkingLotReviews.Ctx(ctx).Where("id", req.Id).Delete()
	if err != nil {
		return nil, err
	}
	res = &parking_lot_review.ParkingLotReviewDeleteRes{Success: true}
	return
}

func (s *sParkingLotReview) ParkingLotReviewDetail(ctx context.Context, req *parking_lot_review.ParkingLotReviewDetailReq) (res *parking_lot_review.ParkingLotReviewDetailRes, err error) {
	review, err := dao.ParkingLotReviews.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, err
	}
	if review.IsEmpty() {
		return nil, nil
	}
	info := &parking_lot_review.ParkingLotReviewInfo{
		Id:        review["id"].Int64(),
		LotId:     review["lot_id"].Int64(),
		UserId:    review["user_id"].Int64(),
		Rating:    review["rating"].Int(),
		Comment:   review["comment"].String(),
		CreatedAt: review["created_at"].GTime().Format("Y-m-d H:i:s"),
	}
	res = &parking_lot_review.ParkingLotReviewDetailRes{Review: info}
	return
}

func (s *sParkingLotReview) ParkingLotReviewList(ctx context.Context, req *parking_lot_review.ParkingLotReviewListReq) (res *parking_lot_review.ParkingLotReviewListRes, err error) {
	m := dao.ParkingLotReviews.Ctx(ctx)
	if req.LotId != 0 {
		m = m.Where("lot_id", req.LotId)
	}
	all, err := m.All()
	if err != nil {
		return nil, err
	}
	reviews := make([]parking_lot_review.ParkingLotReviewInfo, 0, len(all))
	for _, r := range all {
		reviews = append(reviews, parking_lot_review.ParkingLotReviewInfo{
			Id:        r["id"].Int64(),
			LotId:     r["lot_id"].Int64(),
			UserId:    r["user_id"].Int64(),
			Rating:    r["rating"].Int(),
			Comment:   r["comment"].String(),
			CreatedAt: r["created_at"].GTime().Format("Y-m-d H:i:s"),
		})
	}
	res = &parking_lot_review.ParkingLotReviewListRes{Reviews: reviews}
	return
}

func InitParkingLotReview() {
	service.RegisterParkingLotReview(&sParkingLotReview{})
}

func init() {
	InitParkingLotReview()
}
