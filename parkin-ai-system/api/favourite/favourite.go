// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package favourite

import (
	"context"

	"parkin-ai-system/api/favourite/favourite"
)

type IFavouriteFavourite interface {
	FavoriteAdd(ctx context.Context, req *favourite.FavoriteAddReq) (res *favourite.FavoriteAddRes, err error)
	FavoriteList(ctx context.Context, req *favourite.FavoriteListReq) (res *favourite.FavoriteListRes, err error)
	FavoriteDelete(ctx context.Context, req *favourite.FavoriteDeleteReq) (res *favourite.FavoriteDeleteRes, err error)
}
