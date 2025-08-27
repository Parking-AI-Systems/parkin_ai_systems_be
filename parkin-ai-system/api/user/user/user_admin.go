package user

import "github.com/gogf/gf/v2/frame/g"

// GetAllUsersReq - Admin only endpoint to get all users
type GetAllUsersReq struct {
	g.Meta `path:"/admin/users" method:"get" tags:"Admin" summary:"Get all users (Admin only)"`
	Page   int `json:"page" v:"min:1" d:"1"`
	Size   int `json:"size" v:"min:1,max:100" d:"10"`
}

type GetAllUsersRes struct {
	Users []UserInfo `json:"users"`
	Total int        `json:"total"`
	Page  int        `json:"page"`
	Size  int        `json:"size"`
}

type UserInfo struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	FullName  string `json:"full_name"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

// DeleteUserReq - Admin only endpoint to delete users
type DeleteUserReq struct {
	g.Meta `path:"/admin/users/:id" method:"delete" tags:"Admin" summary:"Delete user (Admin only)"`
	Id     int64 `json:"id" in:"path" name:"id"`
}

type DeleteUserRes struct {
	Message string `json:"message"`
}

// UpdateUserRoleReq - Admin only endpoint to update user roles
type UpdateUserRoleReq struct {
	g.Meta `path:"/admin/users/:id/role" method:"put" tags:"Admin" summary:"Update user role (Admin only)"`
	Id     int64  `json:"id" in:"path" name:"id"`
	Role   string `json:"role" v:"required|in:admin,user#Role must be admin or user"`
}

type UpdateUserRoleRes struct {
	Message string `json:"message"`
}
