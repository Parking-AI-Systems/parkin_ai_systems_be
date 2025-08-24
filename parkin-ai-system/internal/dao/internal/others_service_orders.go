// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// OthersServiceOrdersDao is the data access object for the table others_service_orders.
type OthersServiceOrdersDao struct {
	table    string                     // table is the underlying table name of the DAO.
	group    string                     // group is the database configuration group name of the current DAO.
	columns  OthersServiceOrdersColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler         // handlers for customized model modification.
}

// OthersServiceOrdersColumns defines and stores column names for the table others_service_orders.
type OthersServiceOrdersColumns struct {
	Id            string //
	UserId        string //
	VehicleId     string //
	ServiceId     string //
	LotId         string //
	ScheduledTime string //
	Status        string //
	Price         string //
	PaymentStatus string //
	CreatedAt     string //
}

// othersServiceOrdersColumns holds the columns for the table others_service_orders.
var othersServiceOrdersColumns = OthersServiceOrdersColumns{
	Id:            "id",
	UserId:        "user_id",
	VehicleId:     "vehicle_id",
	ServiceId:     "service_id",
	LotId:         "lot_id",
	ScheduledTime: "scheduled_time",
	Status:        "status",
	Price:         "price",
	PaymentStatus: "payment_status",
	CreatedAt:     "created_at",
}

// NewOthersServiceOrdersDao creates and returns a new DAO object for table data access.
func NewOthersServiceOrdersDao(handlers ...gdb.ModelHandler) *OthersServiceOrdersDao {
	return &OthersServiceOrdersDao{
		group:    "default",
		table:    "others_service_orders",
		columns:  othersServiceOrdersColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *OthersServiceOrdersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *OthersServiceOrdersDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *OthersServiceOrdersDao) Columns() OthersServiceOrdersColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *OthersServiceOrdersDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *OthersServiceOrdersDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *OthersServiceOrdersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
