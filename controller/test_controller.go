package controller

import (
	"github.com/gofiber/fiber/v2"
)

type TestController interface {
	FindAll(ctx *fiber.Ctx) error
	FindById(ctx *fiber.Ctx) error
}
