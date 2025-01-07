package controller

import (
	"maze-conquest-api/exception"
	"maze-conquest-api/model/domain"
	"maze-conquest-api/model/web"
	"maze-conquest-api/repository"

	"github.com/gofiber/fiber/v2"
)

type UserControllerImpl struct {
	UserRepository repository.UserRepository
}

func NewUserController(userRepository repository.UserRepository) UserController {
	return &UserControllerImpl{
		UserRepository: userRepository,
	}
}

func (controller *UserControllerImpl) FindById(ctx *fiber.Ctx) error {
	params := ctx.AllParams()
	uid := params["uid"]
	if uid == "" {
		panic(exception.NewEmptyUidError(400, "UID is not provided !"))
	}

	var user *domain.User = controller.UserRepository.FindById(ctx, uid)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   user,
	}

	return ctx.Status(200).JSON(webResponse)
}

func (controller *UserControllerImpl) FindAll(ctx *fiber.Ctx) error {
	var users []*domain.User = controller.UserRepository.FindAll(ctx)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   users,
	}

	return ctx.Status(200).JSON(webResponse)
}

type ItemRequest struct {
	Uid      string `json:"uid"`
	ItemName string `json:"itemName"`
	Quantity int    `json:"quantity"`
}

func (controller *UserControllerImpl) UpdateItem(ctx *fiber.Ctx) error {
	item := new(ItemRequest)
	err := ctx.BodyParser(item)
	if err != nil {
		panic(err)
	}

	var users *domain.User = controller.UserRepository.UpdateItem(ctx, item.Uid, item.ItemName, item.Quantity)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   users,
	}

	return ctx.Status(200).JSON(webResponse)
}

func (controller *UserControllerImpl) FindStrongestHero(ctx *fiber.Ctx) error {
	params := ctx.AllParams()
	uid := params["uid"]
	if uid == "" {
		panic(exception.NewEmptyUidError(400, "UID is not provided !"))
	}

	var heroes []*domain.Hero = controller.UserRepository.GetAllHeroes(ctx, uid)
	var userStrongestHero *domain.Hero = nil
	for _, hero := range heroes {
		// Here, `hero` is a pointer to a `domain.Hero`, so you can access its fields
		if userStrongestHero == nil || hero.Level > userStrongestHero.Level {
			userStrongestHero = hero
		}
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   userStrongestHero,
	}

	return ctx.Status(200).JSON(webResponse)
}

func (controller *UserControllerImpl) MazeLevel(ctx *fiber.Ctx) error {
	params := ctx.AllParams()
	uid := params["uid"]
	if uid == "" {
		panic(exception.NewEmptyUidError(400, "UID is not provided !"))
	}

	var maze = controller.UserRepository.MazeLevel(ctx, uid)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   maze,
	}

	return ctx.Status(200).JSON(webResponse)
}

func (controller *UserControllerImpl) Power(ctx *fiber.Ctx) error {
	params := ctx.AllParams()
	uid := params["uid"]
	if uid == "" {
		panic(exception.NewEmptyUidError(400, "UID is not provided !"))
	}

	var heroes []*domain.Hero = controller.UserRepository.GetAllHeroes(ctx, uid)

	var totalPower = 0
	var levelPoint = 0
	for _, hero := range heroes {
		// Calculate the level only
		levelPoint += hero.Level * 10

		// ... etc
		//
	}

	// Calculate for total power
	totalPower = levelPoint

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data: fiber.Map{
			"power":      totalPower,
			"levelPoint": levelPoint,
		},
	}

	return ctx.Status(200).JSON(webResponse)
}
