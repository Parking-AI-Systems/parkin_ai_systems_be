package user

import (
	"github.com/gogf/gf/v2/frame/g"
)

type UserRegisterReq struct {
	g.Meta    `path:"/register" tags:"User" method:"POST" summary:"Register a new user" description:"Creates a new user account."`
	Username  string `json:"username" v:"required|length:3,50#Username is required|Username must be 3-50 characters"`
	Password  string `json:"password" v:"required|length:6,50#Password is required|Password must be 6-50 characters"`
	FullName  string `json:"full_name" v:"length:0,100#Full name must be 0-100 characters"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Gender    string `json:"gender" v:"length:0,20#Gender must be 0-20 characters"`
	BirthDate string `json:"birth_date" v:"length:0,10|regex:^[0-9]{4}-[0-9]{2}-[0-9]{2}$#Birth date must be 0-10 characters|Invalid date format (YYYY-MM-DD)"`
	AvatarUrl string `json:"avatar_url" v:"url|length:0,255#Invalid URL format|Avatar URL must be 0-255 characters"`
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
	g.Meta   `path:"/login" tags:"User" method:"POST" summary:"User login" description:"Authenticates a user and returns access and refresh tokens."`
	Account  string `json:"account" v:"required#Account is required"`
	Password string `json:"password" v:"required#Password is required"`
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
	g.Meta       `path:"/refresh-token" tags:"User" method:"POST" summary:"Refresh token" description:"Refreshes access and refresh tokens."`
	RefreshToken string `json:"refresh_token" v:"required#Refresh token is required"`
}

type UserRefreshTokenRes struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserLogoutReq struct {
	g.Meta `path:"/logout" tags:"User" method:"POST" summary:"User logout" description:"Logs out the user by deactivating tokens." middleware:"middleware.Auth"`
}

type UserLogoutRes struct {
	Message string `json:"message"`
}

type UserProfileReq struct {
	g.Meta `path:"/profile" tags:"User" method:"GET" summary:"Get user profile" description:"Retrieves the authenticated user's profile." middleware:"middleware.Auth"`
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
}

type UserByIdReq struct {
	g.Meta `path:"/users/:id" tags:"User" method:"GET" summary:"Get user by ID" description:"Retrieves a user's profile by ID (admin or self)." middleware:"middleware.Auth"`
	Id     int64 `json:"id" v:"required|min:1#User ID is required|User ID must be positive"`
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
}

type UserUpdateProfileReq struct {
	g.Meta    `path:"/profile" tags:"User" method:"PUT" summary:"Update user profile" description:"Updates the authenticated user's profile." middleware:"middleware.Auth"`
	Username  string `json:"username" v:"length:3,50#Username must be 3-50 characters"`
	Password  string `json:"password" v:"length:6,50#Password must be 6-50 characters"`
	FullName  string `json:"full_name" v:"length:0,100#Full name must be 0-100 characters"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Gender    string `json:"gender" v:"length:0,20#Gender must be 0-20 characters"`
	BirthDate string `json:"birth_date" v:"length:0,10|regex:^[0-9]{4}-[0-9]{2}-[0-9]{2}$#Birth date must be 0-10 characters|Invalid date format (YYYY-MM-DD)"`
	AvatarUrl string `json:"avatar_url" v:"url|length:0,255#Invalid URL format|Avatar URL must be 0-255 characters"`
}

type UserUpdateProfileRes struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type UserListReq struct {
	g.Meta   `path:"/users" tags:"User" method:"GET" summary:"List users" description:"Retrieves a paginated list of users (admin only)." middleware:"middleware.Auth"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Page     int    `json:"page" v:"min:1#Page must be at least 1"`
	PageSize int    `json:"page_size" v:"min:1|max:100#Page size must be between 1 and 100"`
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
}

type UserListRes struct {
	Users []UserItem `json:"users"`
	Total int        `json:"total"`
	Page  int        `json:"page"`
	Size  int        `json:"size"`
}

type UserDeleteReq struct {
	g.Meta `path:"/users/:user_id" tags:"User" method:"DELETE" summary:"Delete a user" description:"Deletes a user and their associated data (admin only)." middleware:"middleware.Auth"`
	UserId int64 `json:"user_id" v:"required|min:1#User ID is required|User ID must be positive"`
}

type UserDeleteRes struct {
	Message string `json:"message"`
}

type UserUpdateRoleReq struct {
	g.Meta `path:"/users/:user_id/role" tags:"User" method:"PUT" summary:"Update user role" description:"Updates a user's role (admin only)." middleware:"middleware.Auth"`
	UserId int64  `json:"user_id" v:"required|min:1#User ID is required|User ID must be positive"`
	Role   string `json:"role" v:"required#Role is required"`
}

type UserUpdateRoleRes struct {
	Message string `json:"message"`
}

type UserUpdateWalletBalanceReq struct {
	g.Meta        `path:"/users/:user_id/wallet" tags:"User" method:"PUT" summary:"Update wallet balance" description:"Updates a user's wallet balance and creates a transaction (admin only)." middleware:"middleware.Auth"`
	UserId        int64   `json:"user_id" v:"required|min:1#User ID is required|User ID must be positive"`
	WalletBalance float64 `json:"wallet_balance" v:"required|min:0#Wallet balance is required|Wallet balance must be non-negative"`
}

type UserUpdateWalletBalanceRes struct {
	Message string `json:"message"`
}
