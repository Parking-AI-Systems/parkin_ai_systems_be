package favourite

import (
	"context"
	"parkin-ai-system/api/favourite"
	"parkin-ai-system/internal/service"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerFavourite) FavouriteAdd(ctx context.Context, req *favourite.FavouriteAddReq) (res *favourite.FavouriteAddRes, err error) {
	res, err = service.Favourite().FavouriteAdd(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}

func (c *ControllerFavourite) FavouriteDelete(ctx context.Context, req *favourite.FavouriteDeleteReq) (res *favourite.FavouriteDeleteRes, err error) {
	res, err = service.Favourite().FavouriteDelete(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}

func (c *ControllerFavourite) FavouriteList(ctx context.Context, req *favourite.FavouriteListReq) (res *favourite.FavouriteListRes, err error) {
	res, err = service.Favourite().FavouriteList(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}

func (c *ControllerFavourite) FavouriteStatus(ctx context.Context, req *favourite.FavouriteStatusReq) (res *favourite.FavouriteStatusRes, err error) {
	res, err = service.Favourite().FavouriteStatus(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}
