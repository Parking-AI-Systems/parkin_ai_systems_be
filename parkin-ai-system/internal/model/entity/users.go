// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Users is the golang structure for table users.
type Users struct {
	Id            int64       `json:"id"            orm:"id"             description:""`
	Username      string      `json:"username"      orm:"username"       description:""`
	PasswordHash  string      `json:"passwordHash"  orm:"password_hash"  description:""`
	FullName      string      `json:"fullName"      orm:"full_name"      description:""`
	Email         string      `json:"email"         orm:"email"          description:""`
	Phone         string      `json:"phone"         orm:"phone"          description:""`
	Role          string      `json:"role"          orm:"role"           description:""`
	AvatarUrl     string      `json:"avatarUrl"     orm:"avatar_url"     description:""`
	WalletBalance float64     `json:"walletBalance" orm:"wallet_balance" description:""`
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:""`
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     description:""`
	DeletedAt     *gtime.Time `json:"deletedAt"     orm:"deleted_at"     description:""`
}

type UserRegisterReq struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FullName  string `json:"fullName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birthDate"`
	AvatarUrl string `json:"avatarUrl"`
}

type UserRegisterRes struct {
	UserId    int64  `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	FullName  string `json:"full_name"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birth_date"`
}

type UserLoginReq struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type UserLoginRes struct {
	AccessToken   string  `json:"access_token"`
	RefreshToken  string  `json:"refresh_token"`
	UserId        int64   `json:"user_id"`
	Username      string  `json:"username"`
	Role          string  `json:"role"`
	WalletBalance float64 `json:"wallet_balance"`
}

type UserRefreshTokenReq struct {
	RefreshToken string `json:"refresh_token"`
}

type UserRefreshTokenRes struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserLogoutReq struct {
}

type UserLogoutRes struct {
	Message string `json:"message"`
}

type UserProfileReq struct {
}

type UserProfileRes struct {
	UserId        int64   `json:"user_id"`
	Username      string  `json:"username"`
	Email         string  `json:"email"`
	Phone         string  `json:"phone"`
	FullName      string  `json:"full_name"`
	Gender        string  `json:"gender"`
	BirthDate     string  `json:"birth_date"`
	Role          string  `json:"role"`
	AvatarUrl     string  `json:"avatar_url"`
	WalletBalance float64 `json:"wallet_balance"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	DeletedAt     string  `json:"deleted_at"`
}

type UserByIdReq struct {
	Id int64 `json:"id"`
}

type UserByIdRes struct {
	UserId        int64   `json:"user_id"`
	Username      string  `json:"username"`
	Email         string  `json:"email"`
	Phone         string  `json:"phone"`
	FullName      string  `json:"full_name"`
	Gender        string  `json:"gender"`
	BirthDate     string  `json:"birth_date"`
	Role          string  `json:"role"`
	AvatarUrl     string  `json:"avatar_url"`
	WalletBalance float64 `json:"wallet_balance"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	DeletedAt     string  `json:"deleted_at"`
}

type UserUpdateProfileReq struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birth_date"`
	AvatarUrl string `json:"avatar_url"`
}

type UserUpdateProfileRes struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type UserListReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

type UserItem struct {
	UserId        int64   `json:"user_id"`
	Username      string  `json:"username"`
	Email         string  `json:"email"`
	Phone         string  `json:"phone"`
	FullName      string  `json:"full_name"`
	Gender        string  `json:"gender"`
	BirthDate     string  `json:"birth_date"`
	Role          string  `json:"role"`
	AvatarUrl     string  `json:"avatar_url"`
	WalletBalance float64 `json:"wallet_balance"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	DeletedAt     string  `json:"deleted_at"`
}

type UserListRes struct {
	Users []UserItem `json:"users"`
	Total int        `json:"total"`
	Page  int        `json:"page"`
	Size  int        `json:"size"`
}

type UserDeleteReq struct {
	UserId int64 `json:"user_id"`
}

type UserDeleteRes struct {
	Message string `json:"message"`
}

type UserUpdateRoleReq struct {
	UserId int64  `json:"user_id"`
	Role   string `json:"role"`
}

type UserUpdateRoleRes struct {
	Message string `json:"message"`
}

type UserUpdateWalletBalanceReq struct {
	UserId        int64   `json:"user_id"`
	WalletBalance float64 `json:"wallet_balance"`
}

type UserUpdateWalletBalanceRes struct {
	Message string `json:"message"`
}

type UserCountReq struct {
	Period    string `json:"period"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type UserCountRes struct {
	TotalUsers int64 `json:"total_users"`
}

type UserRoleDistributionReq struct {
	Period    string `json:"period"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type UserRoleDistributionRes struct {
	Roles []UserRoleItem `json:"roles"`
}

type UserRoleItem struct {
	Role  string `json:"role"`
	Count int64  `json:"count"`
}

type UserRecentRegistrationsReq struct {
	Period    string `json:"period"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type UserRecentRegistrationsRes struct {
	Registrations []UserRegistrationItem `json:"registrations"`
}

type UserRegistrationItem struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}