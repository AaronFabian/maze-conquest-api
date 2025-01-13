package repository

import (
	"maze-conquest-api/model/domain"

	firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
)

type MixStatsRepository interface {
	GetMixStats(ctx *fiber.Ctx, uid string) *domain.MixStats
	UpdatePower(ctx *fiber.Ctx, uid string, newPower int) bool
	GetFirebaseInstance() *firebase.App
}
