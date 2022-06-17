package creditVoucher

import (
	"context"
	"discount/models"
	"discount/repositories"
	"encoding/json"
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

func (c *CreditVoucher) Redeem(ctx context.Context, userID int, code string) error {
	voucher, err := c.repository.Voucher.FindVoucherByCodeAndNotUsed(ctx, code)
	if err != nil {
		return err
	}

	result := c.repository.Voucher.RedeemVoucher(ctx, userID, voucher)
	if result == nil {
		err := c.sendIncreaseRequestToWallet(userID, voucher)
		if err != nil {
			return err
		}
	}
	return result
}

func (c *CreditVoucher) Create(ctx context.Context, rq *models.VoucherRequestModel) (*models.VoucherModel, error) {
	voucher, err := c.repository.Voucher.Create(ctx, rq)
	if err != nil {
		return nil, err
	}

	return voucher, nil
}

func (c *CreditVoucher) GetVoucherCodeUsed(ctx context.Context, code string) (*[]models.RedeemVoucherRequest, error) {
	return c.repository.Voucher.GetVoucherCodeUsed(ctx, code)
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
