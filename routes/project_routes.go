package routes

import (
	"github.com/gofiber/fiber/v2"

	"task_bounty_server/controllers"
	"task_bounty_server/middleware"
)

func ProjectRouter(app *fiber.App) {
	app.Post("/project", middleware.ProtectMiddleware, controllers.CreateProject)
	app.Get("/projects", middleware.ProtectMiddleware, controllers.GetUserProjects)
	app.Get("/project/:projectID", middleware.ProtectMiddleware, controllers.GetProjectById)
	app.Put("/projects/:projectID", middleware.ProtectMiddleware, controllers.TakeProject)
	app.Post("/projects/:projectID/tasks", middleware.ProtectMiddleware, controllers.AddTask)
	app.Get("/projects/:projectID/getTasks", middleware.ProtectMiddleware, controllers.GetTasks)
}
