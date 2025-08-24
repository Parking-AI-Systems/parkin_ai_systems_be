// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ParkingLotReviewsDao is the data access object for the table parking_lot_reviews.
type ParkingLotReviewsDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  ParkingLotReviewsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// ParkingLotReviewsColumns defines and stores column names for the table parking_lot_reviews.
type ParkingLotReviewsColumns struct {
	Id        string //
	LotId     string //
	UserId    string //
	Rating    string //
	Comment   string //
	CreatedAt string //
}

// parkingLotReviewsColumns holds the columns for the table parking_lot_reviews.
var parkingLotReviewsColumns = ParkingLotReviewsColumns{
	Id:        "id",
	LotId:     "lot_id",
	UserId:    "user_id",
	Rating:    "rating",
	Comment:   "comment",
	CreatedAt: "created_at",
}

// NewParkingLotReviewsDao creates and returns a new DAO object for table data access.
func NewParkingLotReviewsDao(handlers ...gdb.ModelHandler) *ParkingLotReviewsDao {
	return &ParkingLotReviewsDao{
		group:    "default",
		table:    "parking_lot_reviews",
		columns:  parkingLotReviewsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ParkingLotReviewsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ParkingLotReviewsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ParkingLotReviewsDao) Columns() ParkingLotReviewsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ParkingLotReviewsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ParkingLotReviewsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ParkingLotReviewsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
