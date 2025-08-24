package user

import (
	"context"
	"fmt"
	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/gogf/gf/v2/errors/gerror"
	"golang.org/x/crypto/bcrypt"
)

type sUser struct {
}

func Init() {
	service.RegisterUser(&sUser{})
}

func (s *sUser) SignUp(ctx context.Context, req *user.RegisterReq) (res *user.RegisterRes, err error) {
	if req.Email == "" || req.Password == "" || req.Username == "" || req.Phone == "" {
		return nil, gerror.New("Invalid input")
	}

	count, err := dao.Users.Ctx(ctx).Where("email", req.Email).Count()
	if err != nil {
		return nil, gerror.New("Database error")
	}
	if count > 0 {
		return nil, gerror.New("Email already exists")
	}

	hashedPwd, err := s.HashPassword(req.Password)
	if err != nil {
		return nil, gerror.New("Failed to hash password")
	}

	userId, err := dao.Users.Ctx(ctx).Data(g.Map{
		"email":    req.Email,
		"password": hashedPwd,
		"username": req.Username,
		"phone":    req.Phone,
	}).InsertAndGetId()
	if err != nil {
		return nil, gerror.New("Failed to create user")
	}
	res = &user.RegisterRes{
		UserID:    fmt.Sprintf("%d", userId),
		Username:  req.Username,
		Email:     req.Email,
		Phone:     req.Phone,
		FullName:  req.FullName,
		Gender:    req.Gender,
		BirthDate: req.BirthDate,
	}
	return
}

func (s *sUser) Login(ctx context.Context, req *user.UserLoginReq) (res *user.UserLoginRes, err error) {
	if req.Email == "" || req.Password == "" {
		return nil, gerror.New("Invalid input")
	}
	return
}

func (s *sUser) Logout(ctx context.Context, req *user.UserLogoutReq) (res *user.UserLogoutRes, err error) {
	return
}

func (s *sUser) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
