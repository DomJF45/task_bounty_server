package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"task_bounty_server/configs"
	"task_bounty_server/routes"
)

func main() {
	app := fiber.New()

	configs.ConnectDB()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173, https://task-bounty.vercel.app",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	routes.UserRouter(app)
	routes.ProjectRouter(app)

	app.Listen(":8080")
}
