// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// WalletTransactionsDao is the data access object for the table wallet_transactions.
type WalletTransactionsDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  WalletTransactionsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// WalletTransactionsColumns defines and stores column names for the table wallet_transactions.
type WalletTransactionsColumns struct {
	Id             string //
	UserId         string //
	Amount         string //
	Type           string //
	Description    string //
	RelatedOrderId string //
	CreatedAt      string //
}

// walletTransactionsColumns holds the columns for the table wallet_transactions.
var walletTransactionsColumns = WalletTransactionsColumns{
	Id:             "id",
	UserId:         "user_id",
	Amount:         "amount",
	Type:           "type",
	Description:    "description",
	RelatedOrderId: "related_order_id",
	CreatedAt:      "created_at",
}

// NewWalletTransactionsDao creates and returns a new DAO object for table data access.
func NewWalletTransactionsDao(handlers ...gdb.ModelHandler) *WalletTransactionsDao {
	return &WalletTransactionsDao{
		group:    "default",
		table:    "wallet_transactions",
		columns:  walletTransactionsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *WalletTransactionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *WalletTransactionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *WalletTransactionsDao) Columns() WalletTransactionsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *WalletTransactionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *WalletTransactionsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *WalletTransactionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
