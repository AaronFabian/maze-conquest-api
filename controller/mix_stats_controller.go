package controller

import (
	"github.com/gofiber/fiber/v2"
)

type MixStatsController interface {
	UpdateUserPower(ctx *fiber.Ctx) error
}
