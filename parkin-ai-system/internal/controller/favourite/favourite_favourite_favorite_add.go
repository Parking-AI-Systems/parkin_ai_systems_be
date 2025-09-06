package favourite

import (
	"context"

	"parkin-ai-system/api/favourite/favourite"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

func (c *ControllerFavourite) FavoriteAdd(ctx context.Context, req *favourite.FavoriteAddReq) (res *favourite.FavoriteAddRes, err error) {
	// Map API request to entity request
	input := &entity.FavoriteAddReq{
		LotId: req.LotId,
	}

	// Call service
	addRes, err := service.Favorite().FavoriteAdd(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &favourite.FavoriteAddRes{
		Id: addRes.Id,
	}
	return res, nil
}
