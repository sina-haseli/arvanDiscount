package repositories

import (
	"database/sql"
	"discount/models"
	"github.com/go-redis/redis/v7"
)

type Voucher interface {
	FindVoucherByCode(code string) (models.VoucherModel, error)
	InsertIntoRedeemedVoucher(userID, voucherID, step int) error
	RedeemVoucher(userID int, voucher models.VoucherModel, getStep func(voucher models.VoucherModel) (int, error)) error
}

type Redis interface {
	Increase(key string) (int, error)
	Decrease(key string) (int, error)
	SetValue(key string, value interface{}) error
	GetValue(key string) (string, error)
}

type Repository struct {
	Voucher Voucher
	Redis   Redis
}

func NewRepository(db *sql.DB, re *redis.Client) *Repository {
	return &Repository{
		Voucher: NewVoucherRepository(db),
		Redis:   NewRedisRepository(re),
	}
}
