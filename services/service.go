package services

import (
	"discount/config"
	"discount/repositories"
	"discount/services/voucher"
	"discount/services/voucher/creditVoucher"
)

type Services struct {
	Voucher voucher.Voucher
}

func NewServices(repository *repositories.Repository, app *config.ConfiguredApp) *Services {
	return &Services{
		Voucher: creditVoucher.NewCreditVoucher(repository, app.Config.App.ComQueueName),
	}
}
