package parking_lot_review

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"parkin-ai-system/internal/consts"
	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/model/do"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

type sParkingLotReview struct{}

func Init() {
	service.RegisterParkingLotReview(&sParkingLotReview{})
}
func init() {
	Init()
}

func (s *sParkingLotReview) ParkingLotReviewAdd(ctx context.Context, req *entity.ParkingLotReviewAddReq) (*entity.ParkingLotReviewAddRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to add a review.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}

	lot, err := dao.ParkingLots.Ctx(ctx).Where("id", req.LotId).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to find the parking lot. Please try again.")
	}
	if lot.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeParkingLotNotFound, "The parking lot could not be found.")
	}

	count, err := dao.ParkingOrders.Ctx(ctx).
		Where("user_id", userID).
		Where("lot_id", req.LotId).
		Where("status", "completed").
		Where("deleted_at IS NULL").
		Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your parking history. Please try again later.")
	}
	if count == 0 {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "You need a completed booking to review this parking lot.")
	}

	if req.Rating < 1 || req.Rating > 5 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Please provide a rating between 1 and 5.")
	}
	if req.Comment != "" && len(req.Comment) > 1000 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Your comment must be less than 1000 characters.")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding your review. Please try again later.")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	data := do.ParkingLotReviews{
		LotId:     req.LotId,
		UserId:    gconv.Int64(userID),
		Rating:    req.Rating,
		Comment:   req.Comment,
		CreatedAt: gtime.Now(),
	}
	lastId, err := dao.ParkingLotReviews.Ctx(ctx).TX(tx).Data(data).InsertAndGetId()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding your review. Please try again later.")
	}

	adminUsers, err := dao.Users.Ctx(ctx).Where("role", consts.RoleAdmin).Where("deleted_at IS NULL").All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding your review. Please try again later.")
	}

	for _, admin := range adminUsers {
		notiData := do.Notifications{
			UserId:         admin["id"].Int64(),
			Type:           "parking_lot_review_added",
			Content:        fmt.Sprintf("New review #%d for parking lot #%d with rating %d.", lastId, req.LotId, req.Rating),
			RelatedOrderId: lastId,
			IsRead:         false,
			CreatedAt:      gtime.Now(),
		}
		_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding your review. Please try again later.")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding your review. Please try again later.")
	}

	return &entity.ParkingLotReviewAddRes{Id: lastId}, nil
}

func (s *sParkingLotReview) ParkingLotReviewList(ctx context.Context, req *entity.ParkingLotReviewListReq) (*entity.ParkingLotReviewListRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to view reviews.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}

	// Build base query for count
	baseQuery := dao.ParkingLotReviews.Ctx(ctx).
		LeftJoin("users", "users.id = parking_lot_reviews.user_id").
		LeftJoin("parking_lots", "parking_lots.id = parking_lot_reviews.lot_id").
		Where("parking_lot_reviews.deleted_at IS NULL").
		Where("users.deleted_at IS NULL OR users.id IS NULL").
		Where("parking_lots.deleted_at IS NULL OR parking_lots.id IS NULL")

	if req.LotId != 0 {
		baseQuery = baseQuery.Where("parking_lot_reviews.lot_id", req.LotId)
	}
	if req.Rating != 0 {
		baseQuery = baseQuery.Where("parking_lot_reviews.rating", req.Rating)
	}

	total, err := baseQuery.Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load reviews. Please try again later.")
	}

	// Build data query with fields
	m := dao.ParkingLotReviews.Ctx(ctx).
		Fields("parking_lot_reviews.*, users.username as username, parking_lots.name as lot_name").
		LeftJoin("users", "users.id = parking_lot_reviews.user_id").
		LeftJoin("parking_lots", "parking_lots.id = parking_lot_reviews.lot_id").
		Where("parking_lot_reviews.deleted_at IS NULL").
		Where("users.deleted_at IS NULL OR users.id IS NULL").
		Where("parking_lots.deleted_at IS NULL OR parking_lots.id IS NULL")

	if req.LotId != 0 {
		m = m.Where("parking_lot_reviews.lot_id", req.LotId)
	}
	if req.Rating != 0 {
		m = m.Where("parking_lot_reviews.rating", req.Rating)
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	m = m.Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize)

	var reviews []struct {
		entity.ParkingLotReviews
		Username string `json:"username"`
		LotName  string `json:"lot_name"`
	}
	err = m.Order("parking_lot_reviews.id DESC").Scan(&reviews)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load reviews. Please try again later.")
	}

	list := make([]entity.ParkingLotReviewItem, 0, len(reviews))
	for _, r := range reviews {
		item := entity.ParkingLotReviewItem{
			Id:        r.Id,
			LotId:     r.LotId,
			LotName:   r.LotName,
			UserId:    r.UserId,
			Username:  r.Username,
			Rating:    r.Rating,
			Comment:   r.Comment,
			CreatedAt: time.Time(r.CreatedAt.Time).Format("2006-01-02 15:04:05"),
			UpdatedAt: time.Time(r.UpdatedAt.Time).Format("2006-01-02 15:04:05"),
		}
		list = append(list, item)
	}

	return &entity.ParkingLotReviewListRes{
		List:  list,
		Total: total,
	}, nil
}

func (s *sParkingLotReview) ParkingLotReviewGet(ctx context.Context, req *entity.ParkingLotReviewGetReq) (*entity.ParkingLotReviewItem, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to view the review.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}

	var review struct {
		entity.ParkingLotReviews
		Username string `json:"username"`
		LotName  string `json:"lot_name"`
	}
	err = dao.ParkingLotReviews.Ctx(ctx).
		Fields("parking_lot_reviews.*, users.username as username, parking_lots.name as lot_name").
		LeftJoin("users", "users.id = parking_lot_reviews.user_id").
		LeftJoin("parking_lots", "parking_lots.id = parking_lot_reviews.lot_id").
		Where("parking_lot_reviews.id", req.Id).
		Where("parking_lot_reviews.deleted_at IS NULL").
		Where("users.deleted_at IS NULL OR users.id IS NULL").
		Where("parking_lots.deleted_at IS NULL OR parking_lots.id IS NULL").
		Scan(&review)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load the review. Please try again later.")
	}
	if review.Id == 0 {
		return nil, gerror.NewCode(consts.CodeNotFound, "The review could not be found.")
	}

	isAdmin := gconv.String(user.Map()["role"]) == consts.RoleAdmin
	if !isAdmin && review.UserId != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "You can only view your own reviews or must be an admin.")
	}

	item := entity.ParkingLotReviewItem{
		Id:        review.Id,
		LotId:     review.LotId,
		LotName:   review.LotName,
		UserId:    review.UserId,
		Username:  review.Username,
		Rating:    review.Rating,
		Comment:   review.Comment,
		CreatedAt: time.Time(review.CreatedAt.Time).Format("2006-01-02 15:04:05"),
		UpdatedAt: time.Time(review.UpdatedAt.Time).Format("2006-01-02 15:04:05"),
	}

	return &item, nil
}

func (s *sParkingLotReview) ParkingLotReviewUpdate(ctx context.Context, req *entity.ParkingLotReviewUpdateReq) (*entity.ParkingLotReviewItem, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to update the review.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}

	review, err := dao.ParkingLotReviews.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to find the review. Please try again.")
	}
	if review.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "The review could not be found.")
	}

	isAdmin := gconv.String(user.Map()["role"]) == consts.RoleAdmin
	if !isAdmin && gconv.Int64(review.Map()["user_id"]) != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "You can only update your own reviews or must be an admin.")
	}

	if req.LotId != 0 {
		lot, err := dao.ParkingLots.Ctx(ctx).Where("id", req.LotId).Where("deleted_at IS NULL").One()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to find the parking lot. Please try again.")
		}
		if lot.IsEmpty() {
			return nil, gerror.NewCode(consts.CodeParkingLotNotFound, "The parking lot could not be found.")
		}
	}

	if req.Rating != 0 && (req.Rating < 1 || req.Rating > 5) {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Please provide a rating between 1 and 5.")
	}
	if req.Comment != "" && len(req.Comment) > 1000 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Your comment must be less than 1000 characters.")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the review. Please try again later.")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	updateData := g.Map{
		"updated_at": gtime.Now(),
	}
	if req.LotId != 0 {
		updateData["lot_id"] = req.LotId
	}
	if req.Rating != 0 {
		updateData["rating"] = req.Rating
	}
	if req.Comment != "" {
		updateData["comment"] = req.Comment
	}

	_, err = dao.ParkingLotReviews.Ctx(ctx).TX(tx).Data(updateData).Where("id", req.Id).Where("deleted_at IS NULL").Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the review. Please try again later.")
	}

	adminUsers, err := dao.Users.Ctx(ctx).Where("role", consts.RoleAdmin).Where("deleted_at IS NULL").All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the review. Please try again later.")
	}

	for _, admin := range adminUsers {
		notiData := do.Notifications{
			UserId:         admin["id"].Int64(),
			Type:           "parking_lot_review_updated",
			Content:        fmt.Sprintf("Review #%d for parking lot #%d has been updated.", req.Id, req.LotId),
			RelatedOrderId: req.Id,
			IsRead:         false,
			CreatedAt:      gtime.Now(),
		}
		_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the review. Please try again later.")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the review. Please try again later.")
	}

	var updatedReview struct {
		entity.ParkingLotReviews
		Username string `json:"username"`
		LotName  string `json:"lot_name"`
	}
	err = dao.ParkingLotReviews.Ctx(ctx).
		Fields("parking_lot_reviews.*, users.username as username, parking_lots.name as lot_name").
		LeftJoin("users", "users.id = parking_lot_reviews.user_id").
		LeftJoin("parking_lots", "parking_lots.id = parking_lot_reviews.lot_id").
		Where("parking_lot_reviews.id", req.Id).
		Where("parking_lot_reviews.deleted_at IS NULL").
		Where("users.deleted_at IS NULL OR users.id IS NULL").
		Where("parking_lots.deleted_at IS NULL OR parking_lots.id IS NULL").
		Scan(&updatedReview)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the review. Please try again later.")
	}

	item := entity.ParkingLotReviewItem{
		Id:        updatedReview.Id,
		LotId:     updatedReview.LotId,
		LotName:   updatedReview.LotName,
		UserId:    updatedReview.UserId,
		Username:  updatedReview.Username,
		Rating:    updatedReview.Rating,
		Comment:   updatedReview.Comment,
		CreatedAt: time.Time(updatedReview.CreatedAt.Time).Format("2006-01-02 15:04:05"),
		UpdatedAt: time.Time(updatedReview.UpdatedAt.Time).Format("2006-01-02 15:04:05"),
	}

	return &item, nil
}

func (s *sParkingLotReview) ParkingLotReviewDelete(ctx context.Context, req *entity.ParkingLotReviewDeleteReq) (*entity.ParkingLotReviewDeleteRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to delete the review.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}

	review, err := dao.ParkingLotReviews.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to find the review. Please try again.")
	}
	if review.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "The review could not be found.")
	}

	isAdmin := gconv.String(user.Map()["role"]) == consts.RoleAdmin
	if !isAdmin && gconv.Int64(review.Map()["user_id"]) != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "You can only delete your own reviews or must be an admin.")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the review. Please try again later.")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = dao.ParkingLotReviews.Ctx(ctx).TX(tx).Data(g.Map{
		"deleted_at": gtime.Now(),
		"updated_at": gtime.Now(),
	}).Where("id", req.Id).Where("deleted_at IS NULL").Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the review. Please try again later.")
	}

	adminUsers, err := dao.Users.Ctx(ctx).Where("role", consts.RoleAdmin).Where("deleted_at IS NULL").All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the review. Please try again later.")
	}

	for _, admin := range adminUsers {
		notiData := do.Notifications{
			UserId:         admin["id"].Int64(),
			Type:           "parking_lot_review_deleted",
			Content:        fmt.Sprintf("Review #%d for parking lot #%d has been deleted.", req.Id, review.Map()["lot_id"]),
			RelatedOrderId: req.Id,
			IsRead:         false,
			CreatedAt:      gtime.Now(),
		}
		_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the review. Please try again later.")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the review. Please try again later.")
	}

	return &entity.ParkingLotReviewDeleteRes{Message: "Review deleted successfully"}, nil
}
func (s *sParkingLotReview) GetMyLotReview(ctx context.Context, req *entity.MyParkingLotReviewReq) (*entity.MyParkingLotReviewRes, error) {
	// userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	return nil, nil
}
