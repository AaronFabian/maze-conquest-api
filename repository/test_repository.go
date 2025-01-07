package repository

import (
	"maze-conquest-api/model/domain"

	"github.com/gofiber/fiber/v2"
)

type TestRepository interface {
	FindAll(ctx *fiber.Ctx) []*domain.User
	FindById(ctx *fiber.Ctx, uuid string) *domain.User
}
