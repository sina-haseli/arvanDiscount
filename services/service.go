package services

import (
	"discount/repositories"
	"discount/services/voucher"
	"discount/services/voucher/creditVoucher"
)

type Services struct {
	Voucher voucher.Voucher
}

func NewServices(repository *repositories.Repository) *Services {
	return &Services{
		Voucher: creditVoucher.NewCreditVoucher(repository),
	}
}
