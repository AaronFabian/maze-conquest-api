package controller

import (
	"maze-conquest-api/exception"
	"maze-conquest-api/model/domain"
	"maze-conquest-api/model/web"
	"maze-conquest-api/repository"

	"github.com/gofiber/fiber/v2"
)

type MixStatsControllerImpl struct {
	MixStatsRepository repository.MixStatsRepository
}

func NewMixStatsRepository(mixStatsRepository repository.MixStatsRepository) MixStatsController {
	return &MixStatsControllerImpl{
		MixStatsRepository: mixStatsRepository,
	}
}

func (controller *MixStatsControllerImpl) UpdateUserPower(ctx *fiber.Ctx) error {
	params := ctx.AllParams()
	uid := params["uid"]
	if uid == "" {
		panic(exception.NewEmptyUidError())
	}

	// Get the user from other repository
	var user *domain.User = repository.NewUserRepository(controller.MixStatsRepository.GetFirebaseInstance()).FindById(ctx, uid)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   user,
	}

	return ctx.Status(200).JSON(webResponse)
}
