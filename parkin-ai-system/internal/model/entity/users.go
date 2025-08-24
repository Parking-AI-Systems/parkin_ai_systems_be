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
}
