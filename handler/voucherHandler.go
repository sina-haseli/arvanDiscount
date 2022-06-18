package handler

import (
	"context"
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

// RedeemVoucher
// @Summary RedeemVoucher.
// @Tags Redeem
// @Accept       json
// @Produce json
// @Param voucher body models.RedeemVoucherRequest true "Voucher Code"
// @Success 200 {object} map[string]interface{}
// @Router /api/voucher/redeem [post]
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

		ctx := context.TODO()

		err = vh.service.Voucher.Redeem(ctx, rq.UserID, rq.Code)
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

// GetVoucherCodeUsed
// @Summary GetVoucherCodeUsed
// @Tags GetVoucherCodeUsed
// @Accept       json
// @Produce json
// @Param voucherCode path string true "voucherCode"
// @Success 200 {object} map[string]interface{}
// @Router /api/voucher/{voucherCode}/used [get]
func (vh *VoucherHandler) GetVoucherCodeUsed() func(c echo.Context) error {
	return func(c echo.Context) error {
		code := c.Param("voucherCode")
		ctx := context.TODO()

		voucher, errs := vh.service.Voucher.GetVoucherCodeUsed(ctx, code)
		println(errs)
		if errs != nil {
			return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"message": "Entered voucher code is invalid"})
		}

		_ = c.JSON(200, voucher)
		return nil
	}
}

// CreateVoucher
// @Summary CreateVoucher.
// @Tags Redeem
// @Accept       json
// @Produce json
// @Param voucher body models.VoucherRequestModel true "Voucher Code"
// @Success 200 {object} map[string]interface{}
// @Router /api/voucher/create [post]
func (vh *VoucherHandler) CreateVoucher() func(c echo.Context) error {
	return func(c echo.Context) error {
		var rq models.VoucherRequestModel
		err := c.Bind(&rq)
		if err != nil {
			return err
		}

		ctx := context.TODO()

		voucher, errs := vh.service.Voucher.Create(ctx, &rq)
		if errs != nil {
			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusOK, voucher)
	}
}
