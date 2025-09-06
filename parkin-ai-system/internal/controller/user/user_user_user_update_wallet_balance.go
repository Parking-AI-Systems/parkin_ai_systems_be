package user

import (
	"context"

	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

func (c *ControllerUser) UserUpdateWalletBalance(ctx context.Context, req *user.UserUpdateWalletBalanceReq) (res *user.UserUpdateWalletBalanceRes, err error) {
	// Map API request to entity request
	input := &entity.UserUpdateWalletBalanceReq{
		UserId:        req.UserId,
		WalletBalance: req.WalletBalance,
	}

	// Call service
	updateRes, err := service.User().UpdateWalletBalance(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &user.UserUpdateWalletBalanceRes{
		Message: updateRes.Message,
	}
	return res, nil
}
