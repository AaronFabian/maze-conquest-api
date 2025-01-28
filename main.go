package main

import (
	"maze-conquest-api/controller"
	"maze-conquest-api/exception"
	"maze-conquest-api/module"
	"maze-conquest-api/repository"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		IdleTimeout:             time.Second * 5,
		ReadTimeout:             time.Second * 5,
		WriteTimeout:            time.Second * 5,
		Prefork:                 false,
		EnableTrustedProxyCheck: true,
		ErrorHandler:            exception.ErrorHandler,
	})

	// middleware
	app.Use(cors.New())
	// app.Use(csrf.New())
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	firebaseApp, err := module.InitFirebase()
	if err != nil {
		panic(err)
	}

	// html file, etc..
	app.Static("/", "./public")

	apiV1 := app.Group("/api/v1")

	testRepository := repository.NewTestRepository(firebaseApp)
	testController := controller.NewTestController(testRepository)
	apiV1.Get("/testDB/", testController.FindAll)
	apiV1.Get("/testDB/:uid", testController.FindById)

	userRepository := repository.NewUserRepository(firebaseApp)
	userController := controller.NewUserController(userRepository)
	apiV1.Get("/users/", userController.FindAll)
	apiV1.Get("/users/:uid", userController.FindById)
	apiV1.Get("/users/:uid/strongest_hero", userController.FindStrongestHero)
	apiV1.Get("/users/:uid/maze_level", userController.MazeLevel)
	apiV1.Get("/users/:uid/power", userController.Power) // * Deprecated
	apiV1.Patch("/users", userController.UpdateItem)

	mixStatsRepository := repository.NewMixStatsRepository(firebaseApp)
	mixStatsController := controller.NewMixStatsController(mixStatsRepository)
	apiV1.Post("/mix_stats/leaderboard", mixStatsController.GetLeaderboard)
	apiV1.Get("/mix_stats/:uid", mixStatsController.GetUserMixStats)
	apiV1.Patch("/mix_stats/:uid/power", mixStatsController.UpdateUserPower)

	statisticRepository := repository.NewStatisticRepository(firebaseApp)
	statisticController := controller.NewStatisticController(statisticRepository)
	apiV1.Get("/statistics/users", statisticController.GetUsers)
	apiV1.Get("/statistics/users/percentile_from_level/:uid", statisticController.GetUserPercentileFromLevel)
	apiV1.Get("/statistics/users/percentile_from_power/:uid", statisticController.GetUserPercentileFromPower)
	apiV1.Get("/statistics/mix_stats", statisticController.GetMixStats)

	// * Websocket in test mode
	apiV1.Use("/ws", module.SocketUpgradeMidl)
	apiV1.Get("/ws/:id", websocket.New(module.SocketImpl))

	// Just a gateway test
	app.Get("/api/v1", testController.Gateway)

	// Check server condition
	go module.CheckEmptyClient()

	url := "localhost:8000"
	if os.Getenv("MODE") == "prod" {
		url = ":8000"
	}
	err = app.Listen(url)
	if err != nil {
		panic(err)
	}
}
