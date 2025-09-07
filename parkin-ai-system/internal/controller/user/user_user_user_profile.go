package user

import (
	"context"

	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerUser) UserProfile(ctx context.Context, req *user.UserProfileReq) (res *user.UserProfileRes, err error) {
	// Map API request to entity request
	input := &entity.UserProfileReq{}

	// Call service
	profileRes, err := service.User().UserProfile(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &user.UserProfileRes{
		UserId:        profileRes.UserId,
		Username:      profileRes.Username,
		Email:         profileRes.Email,
		Phone:         profileRes.Phone,
		FullName:      profileRes.FullName,
		Gender:        profileRes.Gender,
		BirthDate:     profileRes.BirthDate,
		Role:          profileRes.Role,
		AvatarUrl:     profileRes.AvatarUrl,
		WalletBalance: profileRes.WalletBalance,
		CreatedAt:     profileRes.CreatedAt,
		UpdatedAt:     profileRes.UpdatedAt,
		DeletedAt:     profileRes.DeletedAt,
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}

	return res, nil
}
