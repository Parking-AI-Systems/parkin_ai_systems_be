// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ParkingSlotsDao is the data access object for the table parking_slots.
type ParkingSlotsDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  ParkingSlotsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// ParkingSlotsColumns defines and stores column names for the table parking_slots.
type ParkingSlotsColumns struct {
	Id          string //
	LotId       string //
	Code        string //
	IsAvailable string //
	SlotType    string //
	Floor       string //
	CreatedAt   string //
}

// parkingSlotsColumns holds the columns for the table parking_slots.
var parkingSlotsColumns = ParkingSlotsColumns{
	Id:          "id",
	LotId:       "lot_id",
	Code:        "code",
	IsAvailable: "is_available",
	SlotType:    "slot_type",
	Floor:       "floor",
	CreatedAt:   "created_at",
}

// NewParkingSlotsDao creates and returns a new DAO object for table data access.
func NewParkingSlotsDao(handlers ...gdb.ModelHandler) *ParkingSlotsDao {
	return &ParkingSlotsDao{
		group:    "default",
		table:    "parking_slots",
		columns:  parkingSlotsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ParkingSlotsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ParkingSlotsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ParkingSlotsDao) Columns() ParkingSlotsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ParkingSlotsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ParkingSlotsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ParkingSlotsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
