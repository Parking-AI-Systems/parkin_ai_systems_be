package favourite

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

type sFavorite struct{}

func Init() {
	service.RegisterFavorite(&sFavorite{})
}
func init() {
	Init()
}

func (s *sFavorite) FavoriteAdd(ctx context.Context, req *entity.FavoriteAddReq) (*entity.FavoriteAddRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to add a favorite parking lot.")
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
		return nil, gerror.NewCode(consts.CodeNotFound, "The parking lot could not be found.")
	}

	count, err := dao.Favorites.Ctx(ctx).
		Where("user_id", userID).
		Where("lot_id", req.LotId).
		Where("deleted_at IS NULL").
		Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to check favorite parking lots. Please try again later.")
	}
	if count > 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "This parking lot is already in your favorites.")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding the favorite parking lot. Please try again later.")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	data := do.Favorites{
		UserId:    gconv.Int64(userID),
		LotId:     req.LotId,
		CreatedAt: gtime.Now(),
	}
	lastId, err := dao.Favorites.Ctx(ctx).TX(tx).Data(data).InsertAndGetId()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding the favorite parking lot. Please try again later.")
	}

	adminUsers, err := dao.Users.Ctx(ctx).Where("role", "admin").Where("deleted_at IS NULL").All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding the favorite parking lot. Please try again later.")
	}

	for _, admin := range adminUsers {
		notiData := do.Notifications{
			UserId:         gconv.Int64(admin["id"]),
			Type:           "favorite_lot_added",
			Content:        fmt.Sprintf("User #%d added parking lot #%d to favorites.", gconv.Int64(userID), req.LotId),
			RelatedOrderId: lastId,
			IsRead:         false,
			CreatedAt:      gtime.Now(),
		}
		_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding the favorite parking lot. Please try again later.")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while adding the favorite parking lot. Please try again later.")
	}

	return &entity.FavoriteAddRes{Id: lastId}, nil
}

func (s *sFavorite) FavoriteList(ctx context.Context, req *entity.FavoriteListReq) (*entity.FavoriteListRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to view your favorite parking lots.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}

	// Build base query conditions
	baseQuery := dao.Favorites.Ctx(ctx).
		LeftJoin("parking_lots", "parking_lots.id = favorites.lot_id").
		Where("favorites.deleted_at IS NULL").
		Where("parking_lots.deleted_at IS NULL")

	isAdmin := gconv.String(user.Map()["role"]) == consts.RoleAdmin
	if !isAdmin {
		baseQuery = baseQuery.Where("favorites.user_id", userID)
	}

	if req.LotName != "" {
		baseQuery = baseQuery.WhereLike("parking_lots.name", "%"+req.LotName+"%")
	}

	// Count query - use simple field
	total, err := baseQuery.Fields("favorites.id").Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load your favorite parking lots. Please try again later.")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// Data query - use joined fields
	var favorites []struct {
		entity.Favorites
		LotName    string `json:"lot_name"`
		LotAddress string `json:"lot_address"`
	}
	err = baseQuery.Fields("favorites.*, parking_lots.name as lot_name, parking_lots.address as lot_address").
		Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize).
		Order("favorites.id DESC").Scan(&favorites)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load your favorite parking lots. Please try again later.")
	}

	list := make([]entity.FavoriteItem, 0, len(favorites))
	for _, f := range favorites {
		item := entity.FavoriteItem{
			Id:         f.Id,
			UserId:     f.UserId,
			LotId:      f.LotId,
			LotName:    f.LotName,
			LotAddress: f.LotAddress,
			CreatedAt:  time.Time(f.CreatedAt.Time).Format("2006-01-02 15:04:05"),
		}
		list = append(list, item)
	}

	return &entity.FavoriteListRes{
		List:  list,
		Total: total,
	}, nil
}

func (s *sFavorite) FavoriteDelete(ctx context.Context, req *entity.FavoriteDeleteReq) (*entity.FavoriteDeleteRes, error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to delete a favorite parking lot.")
	}

	user, err := dao.Users.Ctx(ctx).Where("id", userID).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if user.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}

	favorite, err := dao.Favorites.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to find the favorite parking lot. Please try again.")
	}
	if favorite.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeNotFound, "The favorite parking lot could not be found.")
	}

	isAdmin := gconv.String(user.Map()["role"]) == consts.RoleAdmin
	if !isAdmin && gconv.Int64(favorite.Map()["user_id"]) != gconv.Int64(userID) {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "You can only remove your own favorite parking lots or must be an admin.")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while removing the favorite parking lot. Please try again later.")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = dao.Favorites.Ctx(ctx).TX(tx).Data(g.Map{
		"deleted_at": gtime.Now(),
		"updated_at": gtime.Now(),
	}).Where("id", req.Id).Where("deleted_at IS NULL").Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while removing the favorite parking lot. Please try again later.")
	}

	adminUsers, err := dao.Users.Ctx(ctx).Where("role", "admin").Where("deleted_at IS NULL").All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while removing the favorite parking lot. Please try again later.")
	}

	for _, admin := range adminUsers {
		notiData := do.Notifications{
			UserId:         gconv.Int64(admin["id"]),
			Type:           "favorite_lot_deleted",
			Content:        fmt.Sprintf("User #%d removed parking lot #%d from favorites.", gconv.Int64(userID), favorite.Map()["lot_id"]),
			RelatedOrderId: req.Id,
			IsRead:         false,
			CreatedAt:      gtime.Now(),
		}
		_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while removing the favorite parking lot. Please try again later.")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while removing the favorite parking lot. Please try again later.")
	}

	return &entity.FavoriteDeleteRes{Message: "Favorite parking lot removed successfully"}, nil
}
