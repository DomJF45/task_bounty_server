package routes

import (
	"github.com/gofiber/fiber/v2"

	"task_bounty_server/controllers"
	"task_bounty_server/middleware"
)

func UserRouter(app *fiber.App) {
	app.Post("/user/register", controllers.RegisterUser)
	app.Post("/user/login", controllers.LoginUser)
	app.Get("/user", middleware.ProtectMiddleware, controllers.GetUser)
}
