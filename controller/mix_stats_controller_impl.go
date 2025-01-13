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

func NewMixStatsController(mixStatsRepository repository.MixStatsRepository) MixStatsController {
	return &MixStatsControllerImpl{
		MixStatsRepository: mixStatsRepository,
	}
}

func (controller *MixStatsControllerImpl) GetUserMixStats(ctx *fiber.Ctx) error {
	params := ctx.AllParams()
	uid := params["uid"]
	if uid == "" {
		panic(exception.NewEmptyUidError())
	}

	var userMixStats *domain.MixStats = controller.MixStatsRepository.GetMixStats(ctx, uid)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   userMixStats,
	}

	return ctx.Status(200).JSON(webResponse)
}

/*
The build this is for later caching, rather than we must look every user in our database
and looping one by one to calculate
*/
func (controller *MixStatsControllerImpl) UpdateUserPower(ctx *fiber.Ctx) error {
	params := ctx.AllParams()
	uid := params["uid"]
	if uid == "" {
		panic(exception.NewEmptyUidError())
	}

	// Get firebase instance reference
	fbApp := controller.MixStatsRepository.GetFirebaseInstance()

	// Get the hereos from other repository
	var heroes []*domain.Hero = repository.NewUserRepository(fbApp).GetAllHeroes(ctx, uid)
	var totalPower = 0
	var levelPoint = 0
	for _, hero := range heroes {
		// Calculate the level only
		levelPoint += hero.Level * 10

		// etc ...
		// ...
	}

	// Calculate for total power
	totalPower = levelPoint

	controller.MixStatsRepository.UpdatePower(ctx, uid, totalPower)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data: fiber.Map{
			"totalPower": totalPower,
		},
	}

	return ctx.Status(200).JSON(webResponse)
}
