package middleware

import (
	"github.com/gofiber/fiber/v2"
	e "hireplus-project/internal/exception"
)

func ErrorHandlerMiddleware(c *fiber.Ctx) error {

	err := c.Next()

	if err != nil {
		return e.HandleHttpErrorFiber(c, err)
	}

	return nil
}
