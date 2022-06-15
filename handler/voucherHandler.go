package handler

import (
	"discount/models"
	"discount/repositories"
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
			return err
		}

		if rq.Code == "" || rq.UserID == 0 {
			return echo.ErrBadRequest
		}

		err = vh.service.Voucher.Redeem(rq.UserID, rq.Code)
		switch err {
		case repositories.VoucherSoldOut:
			return c.JSON(http.StatusNotAcceptable, map[string]interface{}{"message": "Sorry voucher code sold out! :("})
		case repositories.InvalidVoucherCode:
			return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"message": "Entered voucher code is invalid"})
		case repositories.VoucherAlreadyUsed:
			return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"message": "You have already used this code"})
		default:
			fmt.Println(err)
			return echo.ErrInternalServerError
		case nil:
			return c.JSON(http.StatusOK, map[string]interface{}{"message": "Congratulation your credit will be added to your wallet soon"})
		}
	}
}

func (vh *VoucherHandler) GetVoucherCodeUsed() func(c echo.Context) error {
	return func(c echo.Context) error {
		code := c.Param("voucherCode")

		voucher, errs := vh.service.Voucher.GetVoucherCodeUsed(code)
		if errs != nil {
			return echo.ErrInternalServerError
		}

		_ = c.JSON(200, voucher)
		return nil
	}
}

func (vh *VoucherHandler) CreateVoucher() func(c echo.Context) error {
	return func(c echo.Context) error {
		var rq models.VoucherRequestModel
		err := c.Bind(&rq)
		if err != nil {
			return err
		}

		voucher, errs := vh.service.Voucher.Create(&rq)
		if errs != nil {
			return errs
		}

		return c.JSON(http.StatusOK, voucher)
	}
}
