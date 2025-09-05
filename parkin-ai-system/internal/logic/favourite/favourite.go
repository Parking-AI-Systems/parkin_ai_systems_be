package favourite

import (
	"context"
	"parkin-ai-system/internal/consts"
	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/model/do"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

type sFavourite struct{}

func InitFavourite() {
	service.RegisterFavourite(&sFavourite{})
}

func init() {
	InitFavourite()
}
func (s *sFavourite) FavouriteAdd(ctx context.Context, req *entity.FavoritesInput) (res *entity.FavoritesOutput, err error) {
	userId, ok := ctx.Value("user_id").(int64)
	if !ok {
		return nil, gerror.NewCode(consts.CodeUnauthorized)
	}
	lot, err := dao.ParkingLots.Ctx(ctx).Where("id", req.LotId).One()
	if err != nil {
		return nil, err
	}
	if lot.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeParkingLotNotFound)
	}
	exist, err := dao.Favorites.Ctx(ctx).
		Where("user_id", userId).
		Where("lot_id", req.LotId).
		One()
	if err != nil {
		return nil, err
	}
	if !exist.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeAlreadyFavorited)
	}
	now := gtime.Now()
	fav := do.Favorites{
		UserId:    userId,
		LotId:     req.LotId,
		CreatedAt: now,
	}
	_, err = dao.Favorites.Ctx(ctx).Data(fav).Insert()
	if err != nil {
		return nil, err
	}
	res = &entity.FavoritesOutput{
		UserId:    userId,
		LotId:     req.LotId,
		CreatedAt: now,
	}
	return
}

func (s *sFavourite) FavouriteDelete(ctx context.Context, req *entity.FavoritesInput) (res *entity.FavoriteDelRes, err error) {
	userId, ok := ctx.Value("user_id").(int64)
	if !ok {
		return nil, gerror.NewCode(consts.CodeUnauthorized)
	}
	r, err := dao.Favorites.Ctx(ctx).
		Where("user_id", userId).
		Where("lot_id", req.LotId).
		Delete()
	if err != nil {
		return nil, err
	}

	affected, err := r.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, gerror.NewCode(consts.CodeFavoriteNotFound)
	}

	return &entity.FavoriteDelRes{
		Success: true,
		LotId:   req.LotId,
		UserId:  userId,
	}, nil
}

func (s *sFavourite) FavouriteList(ctx context.Context, req *entity.FavoriteListReq) (res *entity.FavouriteListRes, err error) {
	userId, ok := ctx.Value("user_id").(int64)
	if !ok {
		return nil, gerror.NewCode(consts.CodeUnauthorized)
	}
	total, err := dao.Favorites.Ctx(ctx).Where("user_id", userId).Count()
	if err != nil {
		return nil, err
	}

	all, err := dao.Favorites.Ctx(ctx).
		Fields("favorites.id, favorites.lot_id, favorites.created_at, parking_lots.name as lot_name, parking_lots.address").
		LeftJoin("parking_lots", "favorites.lot_id = parking_lots.id").
		Where("favorites.user_id", userId).
		OrderDesc("favorites.created_at").
		Limit(req.PageSize).
		Offset((req.Page - 1) * req.PageSize).
		All()
	if err != nil {
		return nil, err
	}

	favs := make([]entity.FavouriteInfo, 0, len(all))
	for _, r := range all {
		favs = append(favs, entity.FavouriteInfo{
			Id:        r["id"].Int64(),
			LotId:     r["lot_id"].Int64(),
			LotName:   r["lot_name"].String(),
			Address:   r["address"].String(),
			CreatedAt: r["created_at"].GTime().Format("Y-m-d H:i:s"),
		})
	}

	return &entity.FavouriteListRes{
		Favourites: favs,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}, nil
}

func (s *sFavourite) FavouriteStatus(ctx context.Context, req *entity.FavouriteStatusReq) (res *entity.FavouriteStatusRes, err error) {
	userId, ok := ctx.Value("user_id").(int64)
	if !ok {
		return nil, gerror.NewCode(consts.CodeUnauthorized)
	}
	exist, err := dao.Favorites.Ctx(ctx).
		Where("user_id", userId).
		Where("lot_id", req.LotId).
		One()
	if err != nil {
		return nil, err
	}
	res = &entity.FavouriteStatusRes{IsFavourite: !exist.IsEmpty()}
	return
}
