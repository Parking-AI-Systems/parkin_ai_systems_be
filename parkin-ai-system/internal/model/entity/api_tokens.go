// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ApiTokens is the golang structure for table api_tokens.
type ApiTokens struct {
	Id          int64       `json:"id"          orm:"id"          description:""`
	UserId      int64       `json:"userId"      orm:"user_id"     description:""`
	Token       string      `json:"token"       orm:"token"       description:""`
	Description string      `json:"description" orm:"description" description:""`
	IsActive    bool        `json:"isActive"    orm:"is_active"   description:""`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"  description:""`
}
