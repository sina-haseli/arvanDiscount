package voucher

import (
	"context"
	"discount/models"
)

type Voucher interface {
	Redeem(ctx context.Context, userID string, code string) error
	Create(ctx context.Context, rq *models.VoucherRequestModel) (*models.VoucherModel, error)
	GetVoucherCodeUsed(ctx context.Context, code string) (*[]models.RedeemVoucherRequest, error)
}
