package service

import (
	"context"
	"parkin-ai-system/api/favourite"
)

type IFavourite interface {
	FavouriteAdd(ctx context.Context, req *favourite.FavouriteAddReq) (res *favourite.FavouriteAddRes, err error)
	FavouriteDelete(ctx context.Context, req *favourite.FavouriteDeleteReq) (res *favourite.FavouriteDeleteRes, err error)
	FavouriteList(ctx context.Context, req *favourite.FavouriteListReq) (res *favourite.FavouriteListRes, err error)
	FavouriteStatus(ctx context.Context, req *favourite.FavouriteStatusReq) (res *favourite.FavouriteStatusRes, err error)
}

var localFavourite IFavourite

func Favourite() IFavourite {
	if localFavourite == nil {
		panic("implement not found for interface IFavourite, forgot register?")
	}
	return localFavourite
}

func RegisterFavourite(i IFavourite) {
	localFavourite = i
}
