package controller

import (
	"github.com/gofiber/fiber/v2"
)

type UserController interface {
	UpdateItem(ctx *fiber.Ctx) error
	FindAll(ctx *fiber.Ctx) error
	FindById(ctx *fiber.Ctx) error
	FindStrongestHero(ctx *fiber.Ctx) error
	MazeLevel(ctx *fiber.Ctx) error
	Power(ctx *fiber.Ctx) error
}
