package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/nikitamirzani323/wl_apiagen/controllers"
	"github.com/nikitamirzani323/wl_apiagen/middleware"
)

func Init() *fiber.App {
	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Path() == "/"
		},
		Format: "${time} | ${status} | ${latency} | ${ips[0]} | ${method} | ${path} - ${queryParams} ${body}\n",
	}))
	app.Use(recover.New())
	app.Use(compress.New())

	app.Get("/dashboard", monitor.New())
	app.Post("/api/login", controllers.CheckLogin)
	api := app.Group("/api/", middleware.JWTProtected())

	api.Post("valid", controllers.Home)
	api.Post("alladmin", controllers.Adminhome)
	api.Post("saveadmin", controllers.AdminSave)
	api.Post("alladminrule", controllers.Adminrulehome)
	api.Post("saveadminrule", controllers.AdminruleSave)
	api.Post("allcurr", controllers.Currhome)
	api.Post("savecurr", controllers.CurrSave)
	api.Post("agen", controllers.Agenhome)
	api.Post("saveagen", controllers.AgenSave)
	api.Post("agenlistadmin", controllers.AgenListadmin)
	api.Post("saveagenlistadmin", controllers.AgenSavelistadmin)
	api.Post("log", controllers.Loghome)

	return app
}
