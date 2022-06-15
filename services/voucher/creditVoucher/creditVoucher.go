package creditVoucher

import (
	"discount/models"
	"discount/repositories"
	"errors"
	"fmt"
	"strconv"
)

type CreditVoucher struct {
	repository             *repositories.Repository
	communicationQueueName string
}

var VoucherSoldOut = errors.New("voucher sold out")

func NewCreditVoucher(repository *repositories.Repository) *CreditVoucher {
	return &CreditVoucher{
		repository: repository,
	}
}

func (cv *CreditVoucher) Redeem(userID int, code string) error {
	voucher, err := cv.repository.Voucher.FindVoucherByCode(code)
	if err != nil {
		return err
	}

	v, err := cv.getUsedCount(voucher.ID)
	if err != nil {
		return err
	}

	if voucher.Usable <= v {
		return VoucherSoldOut
	}

	return cv.repository.Voucher.RedeemVoucher(userID, voucher, cv.getStep)
}

func (cv *CreditVoucher) getUsedCount(voucherID int) (int, error) {
	var v int
	val, err := cv.repository.Redis.GetValue(getRedisCacheKeyForVoucher(voucherID))
	if err != nil {
		return 0, err
	} else {
		v, err = strconv.Atoi(val)
		if err != nil {
			return 0, err
		}
	}
	return v, nil
}

func (cv *CreditVoucher) getStep(voucher models.VoucherModel) (int, error) {
	v, err := cv.repository.Redis.Increase(getRedisCacheKeyForVoucher(voucher.ID))
	if err != nil {
		return 0, err
	}

	if v > voucher.Usable {
		return 0, VoucherSoldOut
	}

	return v, nil
}

func getRedisCacheKeyForVoucher(voucherID int) string {
	return fmt.Sprintf("voucher:%d", voucherID)
}
