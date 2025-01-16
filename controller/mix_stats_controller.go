package controller

import (
	"github.com/gofiber/fiber/v2"
)

type MixStatsController interface {
	GetUserMixStats(ctx *fiber.Ctx) error
	UpdateUserPower(ctx *fiber.Ctx) error
	GetLeaderboard(ctx *fiber.Ctx) error
}
