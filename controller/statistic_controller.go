package controller

import (
	"github.com/gofiber/fiber/v2"
)

type StatisticController interface {
	GetUsers(*fiber.Ctx) error
	GetMixStats(*fiber.Ctx) error
	GetUserPercentileFromLevel(*fiber.Ctx) error
	GetUserPercentileFromPower(*fiber.Ctx) error
	GetUserLeaderboard(*fiber.Ctx) error
}
