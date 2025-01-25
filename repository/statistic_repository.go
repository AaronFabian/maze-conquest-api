package repository

import (
	"maze-conquest-api/model/domain"

	"github.com/gofiber/fiber/v2"
)

type StatisticRepository interface {
	GetUsers(ctx *fiber.Ctx) []*domain.Statistic
	GetMixStats(ctx *fiber.Ctx) []*domain.Statistic
	GetUserPercentileFromLevel(ctx *fiber.Ctx, uid string) []*domain.Statistic
	GetUserPercentileFromPower(ctx *fiber.Ctx, uid string) []*domain.Statistic
}
