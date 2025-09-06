package favourite

import (
	"context"

	"parkin-ai-system/api/favourite/favourite"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

func (c *ControllerFavourite) FavoriteList(ctx context.Context, req *favourite.FavoriteListReq) (res *favourite.FavoriteListRes, err error) {
	// Map API request to entity request
	input := &entity.FavoriteListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		LotName:  req.LotName,
	}

	// Call service
	listRes, err := service.Favorite().FavoriteList(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity list to API response
	res = &favourite.FavoriteListRes{
		List:  make([]favourite.FavoriteItem, 0, len(listRes.List)),
		Total: listRes.Total,
	}
	for _, item := range listRes.List {
		res.List = append(res.List, favourite.FavoriteItem{
			Id:         item.Id,
			LotId:      item.LotId,
			LotName:    item.LotName,
			LotAddress: item.LotAddress,
			CreatedAt:  item.CreatedAt,
		})
	}
	return res, nil
}
