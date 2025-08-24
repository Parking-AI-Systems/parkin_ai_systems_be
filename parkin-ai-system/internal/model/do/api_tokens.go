// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ApiTokens is the golang structure of table api_tokens for DAO operations like Where/Data.
type ApiTokens struct {
	g.Meta      `orm:"table:api_tokens, do:true"`
	Id          interface{} //
	UserId      interface{} //
	Token       interface{} //
	Description interface{} //
	IsActive    interface{} //
	CreatedAt   *gtime.Time //
}
