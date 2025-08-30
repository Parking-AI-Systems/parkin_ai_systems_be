package favourite

import "context"

type IFavourite interface {
	FavouriteAdd(ctx context.Context, req *FavouriteAddReq) (res *FavouriteAddRes, err error)
	FavouriteDelete(ctx context.Context, req *FavouriteDeleteReq) (res *FavouriteDeleteRes, err error)
	FavouriteList(ctx context.Context, req *FavouriteListReq) (res *FavouriteListRes, err error)
	FavouriteStatus(ctx context.Context, req *FavouriteStatusReq) (res *FavouriteStatusRes, err error)
}
