package repository

import (
	"maze-conquest-api/model/domain"

	"github.com/gofiber/fiber/v2"
)

type UserRepository interface {
	FindAll(ctx *fiber.Ctx) []*domain.User
	FindById(ctx *fiber.Ctx, uid string) *domain.User
	UpdateItem(ctx *fiber.Ctx, uid string, itemName string, quantity int) *domain.User
	GetAllHeroes(ctx *fiber.Ctx, uid string) []*domain.Hero
	MazeLevel(ctx *fiber.Ctx, uid string) *domain.World
}
