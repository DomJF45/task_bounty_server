package middleware

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"task_bounty_server/controllers"
	"task_bounty_server/responses"
)

func ProtectMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": "unauthorized auth header"}})
	}

	fmt.Printf("\nauthHeader: %v", authHeader)
	token := authHeader[len("Bearer "):]

	fmt.Printf("\ntoken: %v", token)
	userID, err := controllers.VerifyToken(token)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{
			Status:  http.StatusUnauthorized,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}
	fmt.Printf("\nuserID from token: %v", userID)
	c.Locals("UserID", userID)

	return c.Next()
}
