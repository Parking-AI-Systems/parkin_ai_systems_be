package user

import (
	"context"

	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerUser) UserById(ctx context.Context, req *user.UserByIdReq) (res *user.UserByIdRes, err error) {
	// Map API request to entity request
	input := &entity.UserByIdReq{
		Id: req.Id,
	}

	// Call service
	userRes, err := service.User().UserById(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &user.UserByIdRes{
		UserId:        userRes.UserId,
		Username:      userRes.Username,
		Email:         userRes.Email,
		Phone:         userRes.Phone,
		FullName:      userRes.FullName,
		Gender:        userRes.Gender,
		BirthDate:     userRes.BirthDate,
		Role:          userRes.Role,
		AvatarUrl:     userRes.AvatarUrl,
		WalletBalance: userRes.WalletBalance,
		CreatedAt:     userRes.CreatedAt,
		UpdatedAt:     userRes.UpdatedAt,
		DeletedAt:     userRes.DeletedAt,
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}

	return res, nil
}
