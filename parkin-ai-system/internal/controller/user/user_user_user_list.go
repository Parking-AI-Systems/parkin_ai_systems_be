package user

import (
	"context"

	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

func (c *ControllerUser) UserList(ctx context.Context, req *user.UserListReq) (res *user.UserListRes, err error) {
	// Map API request to entity request
	input := &entity.UserListReq{
		Username: req.Username,
		Email:    req.Email,
		Role:     req.Role,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	// Call service
	listRes, err := service.User().GetAllUsers(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity list to API response
	res = &user.UserListRes{
		Users: make([]user.UserItem, 0, len(listRes.Users)),
		Total: listRes.Total,
		Page:  req.Page,
		Size:  req.PageSize,
	}
	for _, item := range listRes.Users {
		res.Users = append(res.Users, user.UserItem{
			UserId:        item.UserId,
			Username:      item.Username,
			Email:         item.Email,
			Phone:         item.Phone,
			FullName:      item.FullName,
			Gender:        item.Gender,
			BirthDate:     item.BirthDate,
			Role:          item.Role,
			AvatarUrl:     item.AvatarUrl,
			WalletBalance: item.WalletBalance,
			CreatedAt:     item.CreatedAt,
			UpdatedAt:     item.UpdatedAt,
		})
	}
	return res, nil
}
