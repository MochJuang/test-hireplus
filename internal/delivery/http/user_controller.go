package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	e "hireplus-project/internal/exception"
	"hireplus-project/internal/model"
	"hireplus-project/internal/service"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService}
}

func (h *UserController) Register(c *fiber.Ctx) error {
	var req model.UserRegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return e.Validation(err)
	}

	user, err := h.userService.Register(req.FirstName, req.LastName, req.PhoneNumber, req.Address, req.Pin)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.JSON(fiber.Map{"status": "SUCCESS", "result": user})
}

func (h *UserController) Login(c *fiber.Ctx) error {
	var req model.UserLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	accessToken, refreshToken, err := h.userService.Login(req.PhoneNumber, req.Pin)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"status": "SUCCESS", "access_token": accessToken, "refresh_token": refreshToken})
}

func (h *UserController) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var req model.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	user, err := h.userService.UpdateProfile(userID, req.FirstName, req.LastName, req.Address)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "SUCCESS", "result": user})
}
