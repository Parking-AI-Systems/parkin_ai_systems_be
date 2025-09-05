package favourite

import (
	"context"
	"parkin-ai-system/api/favourite"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerFavourite) FavouriteAdd(ctx context.Context, req *favourite.FavouriteAddReq) (res *favourite.FavouriteAddRes, err error) {
	input := &entity.FavoritesInput{
		LotId: req.LotId,
	}

	_, err = service.Favourite().FavouriteAdd(ctx, input)
	if err != nil {
		return nil, err
	}

	res = &favourite.FavouriteAddRes{
		Success: true,
	}
	return
}

func (c *ControllerFavourite) FavouriteDelete(ctx context.Context, req *favourite.FavouriteDeleteReq) (res *favourite.FavouriteDeleteRes, err error) {
	input := &entity.FavoritesInput{
		LotId: req.LotId,
	}
	_, err = service.Favourite().FavouriteDelete(ctx, input)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}

func (c *ControllerFavourite) FavouriteList(ctx context.Context, req *favourite.FavouriteListReq) (res *favourite.FavouriteListRes, err error) {
	input := &entity.FavoriteListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	entityRes, err := service.Favourite().FavouriteList(ctx, input)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	favourites := make([]favourite.FavouriteInfo, 0, len(entityRes.Favourites))
	for _, fav := range entityRes.Favourites {
		favourites = append(favourites, favourite.FavouriteInfo{
			LotId:     fav.LotId,
			LotName:   fav.LotName,
			Address:   fav.Address,
			CreatedAt: fav.CreatedAt,
			Id:        fav.Id,
		})
	}
	res = &favourite.FavouriteListRes{
		Favourites: favourites,
		Page:       entityRes.Page,
		PageSize:   entityRes.PageSize,
		Total:      entityRes.Total,
	}
	return res, nil
}

func (c *ControllerFavourite) FavouriteStatus(ctx context.Context, req *favourite.FavouriteStatusReq) (res *favourite.FavouriteStatusRes, err error) {
	input := &entity.FavouriteStatusReq{
		LotId: req.LotId,
	}
	output, err := service.Favourite().FavouriteStatus(ctx, input)
	if err != nil {
		return nil, err
	}

	res = &favourite.FavouriteStatusRes{
		IsFavourite: output.IsFavourite,
	}
	return res, nil
}
