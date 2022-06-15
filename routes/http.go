package routes

import (
	"discount/handler"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, handler *handler.BaseHandler) {
	api := e.Group("/api")
	api.POST("/voucher/redeem", handler.Voucher.RedeemVoucher())
}
