package controller

import (
	"github.com/gofiber/fiber/v2"
)

type UserController interface {
	UpdateItem(*fiber.Ctx) error
	FindAll(*fiber.Ctx) error
	FindById(*fiber.Ctx) error
	FindStrongestHero(*fiber.Ctx) error
	MazeLevel(*fiber.Ctx) error
	Power(*fiber.Ctx) error
	IsExist(*fiber.Ctx) error
}
