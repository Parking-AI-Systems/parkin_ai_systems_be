package user

import (
	"context"

	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerUser) UserRefreshToken(ctx context.Context, req *user.UserRefreshTokenReq) (res *user.UserRefreshTokenRes, err error) {
	// Map API request to entity request
	input := &entity.UserRefreshTokenReq{
		RefreshToken: req.RefreshToken,
	}

	// Call service
	refreshRes, err := service.User().RefreshToken(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &user.UserRefreshTokenRes{
		AccessToken:  refreshRes.AccessToken,
		RefreshToken: refreshRes.RefreshToken,
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}

	return res, nil
}
