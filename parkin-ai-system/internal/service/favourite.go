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
	IFavourite interface {
		FavouriteAdd(ctx context.Context, req *entity.FavoritesInput) (res *entity.FavoritesOutput, err error)
		FavouriteDelete(ctx context.Context, req *entity.FavoritesInput) (res *entity.FavoriteDelRes, err error)
		FavouriteList(ctx context.Context, req *entity.FavoriteListReq) (res *entity.FavouriteListRes, err error)
		FavouriteStatus(ctx context.Context, req *entity.FavouriteStatusReq) (res *entity.FavouriteStatusRes, err error)
	}
)

var (
	localFavourite IFavourite
)

func Favourite() IFavourite {
	if localFavourite == nil {
		panic("implement not found for interface IFavourite, forgot register?")
	}
	return localFavourite
}

func RegisterFavourite(i IFavourite) {
	localFavourite = i
}
