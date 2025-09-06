package parking_lot_review

import (
	"context"
	"fmt"
	"parkin-ai-system/internal/consts"
	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/model/do"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
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
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}

	lot, err := dao.ParkingLots.Ctx(ctx).Where("id", req.LotId).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking parking lot")
	}
	if lot.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeParkingLotNotFound, "Parking lot not found")
	}

	count, err := dao.ParkingOrders.Ctx(ctx).
		Where("user_id", userID).
		Where("lot_id", req.LotId).
		Where("status", "completed").
		Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking completed orders")
	}
	if count == 0 {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User must have a completed order to review this parking lot")
	}

	if req.Rating < 1 || req.Rating > 5 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Rating must be between 1 and 5")
	}
	if req.Comment != "" && len(req.Comment) > 1000 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Comment must be less than 1000 characters")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error starting transaction")
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
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating review")
	}

	adminUsers, err := dao.Users.Ctx(ctx).Where("role", consts.RoleAdmin).All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving admins")
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
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating notification")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error committing transaction")
	}

	return &entity.ParkingLotReviewAddRes{Id: lastId}, nil
}

func (s *sParkingLotReview) ParkingLotReviewList(ctx context.Context, req *entity.ParkingLotReviewListReq) (*entity.ParkingLotReviewListRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}

	m := dao.ParkingLotReviews.Ctx(ctx).
		Fields("parking_lot_reviews.*, users.username as username, parking_lots.name as lot_name").
		LeftJoin("users", "users.id = parking_lot_reviews.user_id").
		LeftJoin("parking_lots", "parking_lots.id = parking_lot_reviews.lot_id").
		Where("parking_lots.deleted_at IS NULL")

	if req.LotId != 0 {
		m = m.Where("parking_lot_reviews.lot_id", req.LotId)
	}
	if req.Rating != 0 {
		m = m.Where("parking_lot_reviews.rating", req.Rating)
	}

	total, err := m.Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error counting reviews")
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
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving reviews")
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
			CreatedAt: r.CreatedAt.Format("2006-01-02 15:04:05"),
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
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
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
		Where("parking_lots.deleted_at IS NULL").
		Scan(&review)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving review")
	}
	if review.Id == 0 {
		return nil, gerror.NewCode(consts.CodeNotFound, "Review not found")
	}

	isAdmin := gconv.String(user.Map()["role"]) == "admin"
	if !isAdmin && review.UserId != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "You can only view your own reviews or must be an admin")
	}

	item := entity.ParkingLotReviewItem{
		Id:        review.Id,
		LotId:     review.LotId,
		LotName:   review.LotName,
		UserId:    review.UserId,
		Username:  review.Username,
		Rating:    review.Rating,
		Comment:   review.Comment,
		CreatedAt: review.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return &item, nil
}

func (s *sParkingLotReview) ParkingLotReviewUpdate(ctx context.Context, req *entity.ParkingLotReviewUpdateReq) (*entity.ParkingLotReviewItem, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}

	review, err := dao.ParkingLotReviews.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking review")
	}
	if review.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "Review not found")
	}

	isAdmin := gconv.String(user.Map()["role"]) == "admin"
	if !isAdmin && gconv.Int64(review.Map()["user_id"]) != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "You can only update your own reviews or must be an admin")
	}

	if req.LotId != 0 {
		lot, err := dao.ParkingLots.Ctx(ctx).Where("id", req.LotId).Where("deleted_at IS NULL").One()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking parking lot")
		}
		if lot.IsEmpty() {
			return nil, gerror.NewCode(consts.CodeParkingLotNotFound, "Parking lot not found")
		}
	}

	if req.Rating != 0 && (req.Rating < 1 || req.Rating > 5) {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Rating must be between 1 and 5")
	}
	if req.Comment != "" && len(req.Comment) > 1000 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Comment must be less than 1000 characters")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error starting transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	updateData := g.Map{}
	if req.LotId != 0 {
		updateData["lot_id"] = req.LotId
	}
	if req.Rating != 0 {
		updateData["rating"] = req.Rating
	}
	if req.Comment != "" {
		updateData["comment"] = req.Comment
	}

	_, err = dao.ParkingLotReviews.Ctx(ctx).TX(tx).Data(updateData).Where("id", req.Id).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error updating review")
	}

	adminUsers, err := dao.Users.Ctx(ctx).Where("role", "admin").All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving admins")
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
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating notification")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error committing transaction")
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
		Scan(&updatedReview)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving updated review")
	}

	item := entity.ParkingLotReviewItem{
		Id:        updatedReview.Id,
		LotId:     updatedReview.LotId,
		LotName:   updatedReview.LotName,
		UserId:    updatedReview.UserId,
		Username:  updatedReview.Username,
		Rating:    updatedReview.Rating,
		Comment:   updatedReview.Comment,
		CreatedAt: updatedReview.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return &item, nil
}

func (s *sParkingLotReview) ParkingLotReviewDelete(ctx context.Context, req *entity.ParkingLotReviewDeleteReq) (*entity.ParkingLotReviewDeleteRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}

	review, err := dao.ParkingLotReviews.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking review")
	}
	if review.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "Review not found")
	}

	isAdmin := gconv.String(user.Map()["role"]) == "admin"
	if !isAdmin && gconv.Int64(review.Map()["user_id"]) != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "You can only delete your own reviews or must be an admin")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error starting transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = dao.ParkingLotReviews.Ctx(ctx).TX(tx).Where("id", req.Id).Delete()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error deleting review")
	}

	adminUsers, err := dao.Users.Ctx(ctx).Where("role", "admin").All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving admins")
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
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating notification")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error committing transaction")
	}

	return &entity.ParkingLotReviewDeleteRes{Message: "Review deleted successfully"}, nil
}
