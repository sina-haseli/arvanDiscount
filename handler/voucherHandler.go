package handler

import (
	"discount/models"
	"discount/services"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type VoucherHandler struct {
	service *services.Services
}

func newVoucherHandler(service *services.Services) *VoucherHandler {
	return &VoucherHandler{
		service: service,
	}
}

func (vh *VoucherHandler) RedeemVoucher() func(c echo.Context) error {
	return func(c echo.Context) error {
		var rq models.RedeemVoucherRequest
		err := c.Bind(&rq)
		if err != nil {
			fmt.Println("could not bind request")
			return err
		}

		if rq.Code == "" || rq.UserID == 0 {
			return echo.ErrBadRequest
		}

		err = vh.service.Voucher.Redeem(rq.UserID, rq.Code)
		switch err {
		default:
			fmt.Println(err)
			return echo.ErrInternalServerError
		case nil:
			return c.JSON(http.StatusOK, map[string]interface{}{"message": "Congratulations! You have redeemed your voucher"})
		}
	}
}
