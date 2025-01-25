package controller

import (
	"maze-conquest-api/exception"
	"maze-conquest-api/model/web"
	"maze-conquest-api/repository"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type StatisticControllerImpl struct {
	StatisticRepository repository.StatisticRepository
}

func NewStatisticController(repository repository.StatisticRepository) StatisticController {
	return &StatisticControllerImpl{
		StatisticRepository: repository,
	}
}

func (controller *StatisticControllerImpl) GetUsers(ctx *fiber.Ctx) error {
	statistic := controller.StatisticRepository.GetUsers(ctx)

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   statistic,
	}
	return ctx.Status(200).JSON(webResponse)
}

func (controller *StatisticControllerImpl) GetMixStats(ctx *fiber.Ctx) error {
	webResponse := web.WebResponse{
		Code:   http.StatusAccepted,
		Status: "ACCEPTED",
		Data:   fiber.Map{},
	}
	return ctx.Status(202).JSON(webResponse)
}

func (controller *StatisticControllerImpl) GetUserPercentileFromLevel(ctx *fiber.Ctx) error {
	params := ctx.AllParams()
	uid := params["uid"]
	if uid == "" {
		panic(exception.NewEmptyUidError())
	}

	statistics := controller.StatisticRepository.GetUserPercentileFromLevel(ctx, uid)
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   statistics,
	}

	return ctx.Status(200).JSON(webResponse)
}

func (controller *StatisticControllerImpl) GetUserPercentileFromPower(ctx *fiber.Ctx) error {
	params := ctx.AllParams()
	uid := params["uid"]
	if uid == "" {
		panic(exception.NewEmptyUidError())
	}

	statistics := controller.StatisticRepository.GetUserPercentileFromPower(ctx, uid)
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   statistics,
	}

	return ctx.Status(200).JSON(webResponse)
}
