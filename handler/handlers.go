package handler

import (
	"discount/services"
	"github.com/labstack/echo/v4"
)

type BaseHandler struct {
	Voucher Voucher
}

type Voucher interface {
	RedeemVoucher() func(c echo.Context) error
	GetVoucherCodeUsed() func(c echo.Context) error
	CreateVoucher() func(c echo.Context) error
}

func NewBaseHandler(services *services.Services) *BaseHandler {
	return &BaseHandler{
		Voucher: newVoucherHandler(services),
	}
}
