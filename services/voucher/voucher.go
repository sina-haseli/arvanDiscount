package voucher

import "discount/models"

type Voucher interface {
	Redeem(userID int, code string) error
	Create(rq *models.VoucherRequestModel) (*models.VoucherModel, error)
	GetVoucherCodeUsed(code string) (*models.RedeemVoucherRequest, error)
}
