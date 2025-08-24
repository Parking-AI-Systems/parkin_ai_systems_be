// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// OthersServiceDao is the data access object for the table others_service.
type OthersServiceDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  OthersServiceColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// OthersServiceColumns defines and stores column names for the table others_service.
type OthersServiceColumns struct {
	Id              string //
	LotId           string //
	Name            string //
	Description     string //
	Price           string //
	DurationMinutes string //
	IsActive        string //
	CreatedAt       string //
}

// othersServiceColumns holds the columns for the table others_service.
var othersServiceColumns = OthersServiceColumns{
	Id:              "id",
	LotId:           "lot_id",
	Name:            "name",
	Description:     "description",
	Price:           "price",
	DurationMinutes: "duration_minutes",
	IsActive:        "is_active",
	CreatedAt:       "created_at",
}

// NewOthersServiceDao creates and returns a new DAO object for table data access.
func NewOthersServiceDao(handlers ...gdb.ModelHandler) *OthersServiceDao {
	return &OthersServiceDao{
		group:    "default",
		table:    "others_service",
		columns:  othersServiceColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *OthersServiceDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *OthersServiceDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *OthersServiceDao) Columns() OthersServiceColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *OthersServiceDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *OthersServiceDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *OthersServiceDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
