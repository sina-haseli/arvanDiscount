package repositories

import (
	"context"
	"database/sql"
	"discount/models"
	"github.com/go-redis/redis/v7"
)

type Voucher interface {
	FindVoucherByCode(ctx context.Context, code string) (models.VoucherModel, error)
	FindVoucherByCodeAndNotUsed(ctx context.Context, code string) (models.VoucherModel, error)
	InsertIntoRedeemedVoucher(userID, voucherID int) error
	IsUserRedeemedVoucherBefore(userID, voucherID int) (bool, error)
	RedeemVoucher(ctx context.Context, userID int, voucher models.VoucherModel) error
	Create(ctx context.Context, rq *models.VoucherRequestModel) (*models.VoucherModel, error)
	GetVoucherCodeUsed(ctx context.Context, code string) (*[]models.RedeemVoucherRequest, error)
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
