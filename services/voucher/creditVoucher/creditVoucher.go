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

func (c *CreditVoucher) Redeem(ctx context.Context, userID string, code string) error {
	voucher, err := c.repository.Voucher.FindVoucherByCode(ctx, code)
	if err != nil {
		return err
	}

	err = c.repository.Voucher.RedeemVoucher(ctx, userID, voucher.ID)
	if err != nil {
		return err
	}

	err = c.sendIncreaseRequestToWallet(userID, voucher)
	if err != nil {
		return err
	}

	return nil
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

func (c *CreditVoucher) sendIncreaseRequestToWallet(userID string, voucher models.VoucherModel) error {
	d, err := json.Marshal(models.IncreaseRequestModel{
		UserID:    userID,
		VoucherID: voucher.ID,
		Amount:    voucher.Amount,
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
