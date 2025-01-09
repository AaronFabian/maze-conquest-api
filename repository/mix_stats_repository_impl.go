package repository

import (
	"maze-conquest-api/model/domain"

	firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
)

type MixStatsRepositoryImpl struct {
	FbApp *firebase.App
}

func NewMixStatsRepositoryImpl(fbApp *firebase.App) MixStatsRepository {
	return &MixStatsRepositoryImpl{
		FbApp: fbApp,
	}
}

func (mixStatsRepository *MixStatsRepositoryImpl) UpdatePower(ctx *fiber.Ctx, uid string) *domain.MixStats {
	return nil
}

func (mixStatsRepository *MixStatsRepositoryImpl) GetFirebaseInstance() *firebase.App {
	return mixStatsRepository.FbApp
}
