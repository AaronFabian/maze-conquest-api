package controller

import (
	"github.com/gofiber/fiber/v2"
)

type TestController interface {
	Gateway(*fiber.Ctx) error
	FindAll(ctx *fiber.Ctx) error
	FindById(ctx *fiber.Ctx) error
}
