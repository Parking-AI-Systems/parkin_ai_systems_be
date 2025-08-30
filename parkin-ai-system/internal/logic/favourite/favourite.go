package favourite

import (
	"context"
	"parkin-ai-system/api/favourite"
	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/model/do"
	"parkin-ai-system/internal/service"
	"github.com/gogf/gf/v2/os/gtime"
)

type sFavourite struct{}

func (s *sFavourite) FavouriteAdd(ctx context.Context, req *favourite.FavouriteAddReq) (res *favourite.FavouriteAddRes, err error) {
	userId := ctx.Value("user_id")
	fav := do.Favorites{
		UserId:    userId,
		LotId:     req.LotId,
		CreatedAt: gtime.Now(),
	}
	_, err = dao.Favorites.Ctx(ctx).Data(fav).Insert()
	if err != nil {
		return nil, err
	}
	res = &favourite.FavouriteAddRes{Success: true}
	return
}

func (s *sFavourite) FavouriteDelete(ctx context.Context, req *favourite.FavouriteDeleteReq) (res *favourite.FavouriteDeleteRes, err error) {
	userId := ctx.Value("user_id")
	_, err = dao.Favorites.Ctx(ctx).Where("user_id", userId).Where("lot_id", req.LotId).Delete()
	if err != nil {
		return nil, err
	}
	res = &favourite.FavouriteDeleteRes{Success: true}
	return
}

func (s *sFavourite) FavouriteList(ctx context.Context, req *favourite.FavouriteListReq) (res *favourite.FavouriteListRes, err error) {
	userId := ctx.Value("user_id")
	all, err := dao.Favorites.Ctx(ctx).Where("user_id", userId).All()
	if err != nil {
		return nil, err
	}
	favs := make([]favourite.FavouriteInfo, 0, len(all))
	for _, r := range all {
		favs = append(favs, favourite.FavouriteInfo{
			Id:        r["id"].Int64(),
			LotId:     r["lot_id"].Int64(),
			CreatedAt: r["created_at"].GTime().Format("Y-m-d H:i:s"),
		})
	}
	res = &favourite.FavouriteListRes{Favourites: favs}
	return
}

func (s *sFavourite) FavouriteStatus(ctx context.Context, req *favourite.FavouriteStatusReq) (res *favourite.FavouriteStatusRes, err error) {
	userId := ctx.Value("user_id")
	count, err := dao.Favorites.Ctx(ctx).Where("user_id", userId).Where("lot_id", req.LotId).Count()
	if err != nil {
		return nil, err
	}
	res = &favourite.FavouriteStatusRes{IsFavourite: count > 0}
	return
}

func InitFavourite() {
	service.RegisterFavourite(&sFavourite{})
}

func init() {
	InitFavourite()
}
