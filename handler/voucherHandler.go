package handler

import (
	"discount/models"
	"discount/repositories"
	"discount/services"
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

		if rq.Code == "" || rq.UserID == "" {
			return echo.ErrBadRequest
		}

		err = vh.service.Voucher.Redeem(c.Request().Context(), rq.UserID, rq.Code)
		switch err {
		case repositories.VoucherSoldOut:
			return c.JSON(http.StatusNotAcceptable, map[string]interface{}{"message": "Sorry voucher code sold out! :("})
		case repositories.InvalidVoucherCode:
			return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"message": "Entered voucher code is invalid"})
		case repositories.VoucherAlreadyUsed:
			return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"message": "You have already used this code"})
		default:
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

		vouchers, errs := vh.service.Voucher.GetVoucherCodeUsed(c.Request().Context(), code)
		if errs != nil {
			switch errs {
			case repositories.InvalidVoucherCode:
				return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"message": "Entered voucher code is invalid"})
			default:
				return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Entered voucher code is invalid"})
			}
		}

		return c.JSON(200, vouchers)
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

		voucher, errs := vh.service.Voucher.Create(c.Request().Context(), &rq)
		if errs != nil {
			switch errs {
			case repositories.VoucherExist:
				return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"message": "voucher code is exist"})
			default:
				return echo.ErrInternalServerError
			}
		}

		return c.JSON(http.StatusOK, voucher)
	}
}
