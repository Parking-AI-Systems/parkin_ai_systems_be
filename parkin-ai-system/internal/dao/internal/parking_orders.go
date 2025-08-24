// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ParkingOrdersDao is the data access object for the table parking_orders.
type ParkingOrdersDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  ParkingOrdersColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// ParkingOrdersColumns defines and stores column names for the table parking_orders.
type ParkingOrdersColumns struct {
	Id            string //
	UserId        string //
	SlotId        string //
	LotId         string //
	StartTime     string //
	EndTime       string //
	Status        string //
	Price         string //
	PaymentStatus string //
	CreatedAt     string //
	UpdatedAt     string //
}

// parkingOrdersColumns holds the columns for the table parking_orders.
var parkingOrdersColumns = ParkingOrdersColumns{
	Id:            "id",
	UserId:        "user_id",
	SlotId:        "slot_id",
	LotId:         "lot_id",
	StartTime:     "start_time",
	EndTime:       "end_time",
	Status:        "status",
	Price:         "price",
	PaymentStatus: "payment_status",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewParkingOrdersDao creates and returns a new DAO object for table data access.
func NewParkingOrdersDao(handlers ...gdb.ModelHandler) *ParkingOrdersDao {
	return &ParkingOrdersDao{
		group:    "default",
		table:    "parking_orders",
		columns:  parkingOrdersColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ParkingOrdersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ParkingOrdersDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ParkingOrdersDao) Columns() ParkingOrdersColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ParkingOrdersDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ParkingOrdersDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ParkingOrdersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
