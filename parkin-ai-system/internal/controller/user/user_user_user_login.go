package user

import (
	"context"

	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerUser) UserLogin(ctx context.Context, req *user.UserLoginReq) (res *user.UserLoginRes, err error) {
	// Map API request to entity request
	input := &entity.UserLoginReq{
		Account:  req.Account,
		Password: req.Password,
	}

	loginRes, err := service.User().Login(ctx, input)
	if err != nil {
		g.Log().Error(ctx, "UserLogin - Service error:", err)
		return nil, err
	}

	// Map entity response to API response
	res = &user.UserLoginRes{
		AccessToken:   loginRes.AccessToken,
		RefreshToken:  loginRes.RefreshToken,
		UserId:        loginRes.UserId,
		Username:      loginRes.Username,
		Role:          loginRes.Role,
		WalletBalance: loginRes.WalletBalance,
	}

	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}

	return res, nil
}
