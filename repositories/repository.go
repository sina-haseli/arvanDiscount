package repositories

import (
	"database/sql"
	"discount/models"
	"github.com/go-redis/redis/v7"
)

type Voucher interface {
	FindVoucherByCode(code string) (models.VoucherModel, error)
	FindVoucherByCodeAndNotUsed(code string) (models.VoucherModel, error)
	InsertIntoRedeemedVoucher(userID, voucherID int) error
	GetRedeemedCount(voucherID int) (int, error)
	IsUserRedeemedVoucherBefore(userID, voucherID int) (bool, error)
	RedeemVoucher(userID int, voucher models.VoucherModel, success func(userID int, voucher models.VoucherModel) error) error
	Create(rq *models.VoucherRequestModel) (*models.VoucherModel, error)
	GetVoucherCodeUsed(code string) (*models.RedeemVoucherRequest, error)
}

type Redis interface {
	Dequeue(queueName string) (string, error)
	Enqueue(message []byte, queueName string) error
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
