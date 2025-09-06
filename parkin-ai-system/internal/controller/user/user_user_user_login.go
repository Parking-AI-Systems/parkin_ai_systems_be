package user

import (
	"context"

	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

func (c *ControllerUser) UserLogin(ctx context.Context, req *user.UserLoginReq) (res *user.UserLoginRes, err error) {
	// Map API request to entity request
	input := &entity.UserLoginReq{
		Account:  req.Account,
		Password: req.Password,
	}

	// Call service
	loginRes, err := service.User().Login(ctx, input)
	if err != nil {
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
	return res, nil
}
