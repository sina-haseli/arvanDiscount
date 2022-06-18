package cmd

import (
	"discount/config"
	"discount/handler"
	"discount/repositories"
	"discount/routes"
	"discount/services"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve arvan discount application",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func init() {
	rootCMD.AddCommand(serveCmd)
}

func serve() {
	ca := config.InitializeConfig()
	rep := repositories.NewRepository(ca.DB, ca.RDB)
	ser := services.NewServices(rep, ca)
	hndl := handler.NewBaseHandler(ser)
	initializeHttpServer(hndl, ca.PORT)
}

func initializeHttpServer(handler *handler.BaseHandler, port string) {
	e := echo.New()
	e.HideBanner = true
	p := prometheus.NewPrometheus("ArvanVoucher", nil)
	p.Use(e)
	routes.RegisterRoutes(e, handler)
	e.Logger.Fatal(e.Start(":" + port))
}
