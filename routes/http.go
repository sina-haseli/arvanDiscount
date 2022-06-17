package routes

import (
	"discount/handler"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "discount/docs"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2

func RegisterRoutes(e *echo.Echo, handler *handler.BaseHandler) {
	api := e.Group("/api")
	api.POST("/voucher/redeem", handler.Voucher.RedeemVoucher())
	api.GET("/voucher/:voucherCode/used", handler.Voucher.GetVoucherCodeUsed())
	api.POST("/voucher/create", handler.Voucher.CreateVoucher())
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
