package favourite

import (
	"context"

	"parkin-ai-system/api/favourite/favourite"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerFavourite) FavoriteDelete(ctx context.Context, req *favourite.FavoriteDeleteReq) (res *favourite.FavoriteDeleteRes, err error) {
	// Map API request to entity request
	input := &entity.FavoriteDeleteReq{
		Id: req.Id,
	}

	// Call service
	deleteRes, err := service.Favorite().FavoriteDelete(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &favourite.FavoriteDeleteRes{
		Message: deleteRes.Message,
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
