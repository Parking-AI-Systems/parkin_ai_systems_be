// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// VehiclesDao is the data access object for the table vehicles.
type VehiclesDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  VehiclesColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// VehiclesColumns defines and stores column names for the table vehicles.
type VehiclesColumns struct {
	Id           string //
	UserId       string //
	LicensePlate string //
	Brand        string //
	Model        string //
	Color        string //
	Type         string //
	CreatedAt    string //
}

// vehiclesColumns holds the columns for the table vehicles.
var vehiclesColumns = VehiclesColumns{
	Id:           "id",
	UserId:       "user_id",
	LicensePlate: "license_plate",
	Brand:        "brand",
	Model:        "model",
	Color:        "color",
	Type:         "type",
	CreatedAt:    "created_at",
}

// NewVehiclesDao creates and returns a new DAO object for table data access.
func NewVehiclesDao(handlers ...gdb.ModelHandler) *VehiclesDao {
	return &VehiclesDao{
		group:    "default",
		table:    "vehicles",
		columns:  vehiclesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *VehiclesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *VehiclesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *VehiclesDao) Columns() VehiclesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *VehiclesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *VehiclesDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *VehiclesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
