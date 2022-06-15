package creditVoucher

import (
	"discount/models"
	"discount/repositories"
	"encoding/json"
	"fmt"
	"strconv"
)

type CreditVoucher struct {
	repository             *repositories.Repository
	communicationQueueName string
}

func NewCreditVoucher(repository *repositories.Repository, comQueue string) *CreditVoucher {
	return &CreditVoucher{
		repository:             repository,
		communicationQueueName: comQueue,
	}
}

func (c *CreditVoucher) Redeem(userID int, code string) error {
	voucher, err := c.repository.Voucher.FindVoucherByCode(code)
	if err != nil {
		return err
	}

	return c.repository.Voucher.RedeemVoucher(userID, voucher, c.sendIncreaseRequestToWallet)

}

func (c *CreditVoucher) Create(rq *models.VoucherRequestModel) (*models.VoucherModel, error) {
	voucher, err := c.repository.Voucher.Create(rq)
	if err != nil {
		return nil, err
	}

	return voucher, nil
}

func (c *CreditVoucher) GetVoucherCodeUsed(code string) (*models.RedeemVoucherRequest, error) {
	return c.repository.Voucher.GetVoucherCodeUsed(code)
}

func (c *CreditVoucher) getUsedCount(voucherID int) (int, error) {
	var v int
	val, err := c.repository.Redis.GetValue(getRedisCacheKeyForVoucher(voucherID))
	if err != nil {
		v, err = c.repository.Voucher.GetRedeemedCount(voucherID)
		if err != nil {
			return 0, err
		}

		_ = c.repository.Redis.SetValue(getRedisCacheKeyForVoucher(voucherID), v)
	} else {
		v, err = strconv.Atoi(val)
		if err != nil {
			return 0, err
		}
	}
	return v, nil
}

func (c *CreditVoucher) sendIncreaseRequestToWallet(userID int, voucher models.VoucherModel) error {
	d, err := json.Marshal(models.IncreaseRequestModel{
		UserID: userID,
		Amount: voucher.Amount,
	})
	if err != nil {
		return err
	}

	err = c.repository.Redis.Enqueue(d, c.communicationQueueName)
	if err != nil {
		return err
	}

	return nil
}

func getRedisCacheKeyForVoucher(voucherID int) string {
	return fmt.Sprintf("voucher:%d", voucherID)
}
