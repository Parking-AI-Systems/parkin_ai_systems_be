// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// NotificationsDao is the data access object for the table notifications.
type NotificationsDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  NotificationsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// NotificationsColumns defines and stores column names for the table notifications.
type NotificationsColumns struct {
	Id             string //
	UserId         string //
	Type           string //
	Content        string //
	RelatedOrderId string //
	IsRead         string //
	CreatedAt      string //
}

// notificationsColumns holds the columns for the table notifications.
var notificationsColumns = NotificationsColumns{
	Id:             "id",
	UserId:         "user_id",
	Type:           "type",
	Content:        "content",
	RelatedOrderId: "related_order_id",
	IsRead:         "is_read",
	CreatedAt:      "created_at",
}

// NewNotificationsDao creates and returns a new DAO object for table data access.
func NewNotificationsDao(handlers ...gdb.ModelHandler) *NotificationsDao {
	return &NotificationsDao{
		group:    "default",
		table:    "notifications",
		columns:  notificationsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *NotificationsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *NotificationsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *NotificationsDao) Columns() NotificationsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *NotificationsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *NotificationsDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *NotificationsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
