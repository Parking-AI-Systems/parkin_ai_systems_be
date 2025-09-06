package others_service

import (
	"context"

	"parkin-ai-system/api/others_service/others_service"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerOthers_service) OthersServiceDelete(ctx context.Context, req *others_service.OthersServiceDeleteReq) (res *others_service.OthersServiceDeleteRes, err error) {
	// Map API request to entity request
	input := &entity.OthersServiceDeleteReq{
		Id: req.Id,
	}

	// Call service
	deleteRes, err := service.OthersService().OthersServiceDelete(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &others_service.OthersServiceDeleteRes{
		Message: deleteRes.Message,
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
