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
	IFavorite interface {
		FavoriteAdd(ctx context.Context, req *entity.FavoriteAddReq) (*entity.FavoriteAddRes, error)
		FavoriteList(ctx context.Context, req *entity.FavoriteListReq) (*entity.FavoriteListRes, error)
		FavoriteDelete(ctx context.Context, req *entity.FavoriteDeleteReq) (*entity.FavoriteDeleteRes, error)
	}
)

var (
	localFavorite IFavorite
)

func Favorite() IFavorite {
	if localFavorite == nil {
		panic("implement not found for interface IFavorite, forgot register?")
	}
	return localFavorite
}

func RegisterFavorite(i IFavorite) {
	localFavorite = i
}
