package main

import (
	"fmt"
	"maze-conquest-api/controller"
	"maze-conquest-api/exception"
	"maze-conquest-api/module"
	"maze-conquest-api/module/handlers"
	"maze-conquest-api/module/webrtc"
	"maze-conquest-api/repository"
	"net/http"
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

	mixStatsRepository := repository.NewMixStatsRepositoryImpl(firebaseApp)
	mixStatsController := controller.NewMixStatsController(mixStatsRepository)
	apiV1.Get("/mix_stats/leaderboard", mixStatsController.GetLeaderboard)
	apiV1.Get("/mix_stats/:uid", mixStatsController.GetUserMixStats)
	apiV1.Patch("/mix_stats/:uid/power", mixStatsController.UpdateUserPower)

	// * In test mode
	webrtc.Streams = make(map[string]*webrtc.Room)
	webrtc.Rooms = make(map[string]*webrtc.Room)
	apiV1.Get("/room/create", handlers.RoomCreate)
	apiV1.Get("/room/:uuid", handlers.Room)
	apiV1.Get("/room/:uuid/chat/websocket", websocket.New(handlers.RoomChatWebsocket))
	apiV1.Get("/room/:uuid/websocket", websocket.New(handlers.RoomWebsocket, websocket.Config{
		HandshakeTimeout: 10 * time.Second,
	}))

	// Just a gateway test
	app.Get("/api/v1", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).JSON(fiber.Map{
			"code":   http.StatusOK,
			"status": "OK",
			"data": fiber.Map{
				"message": "Welcome to API Gateway",
			},
		})
	})

	// Check server condition
	go checkEmptyClient()

	err = app.Listen("localhost:8000")
	if err != nil {
		panic(err)
	}
}

func checkEmptyClient() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		for roomUUID := range webrtc.Rooms {
			fmt.Println(roomUUID)
			room := webrtc.Rooms[roomUUID]
			totalClient := len(room.Hub.Clients)

			// Delete room if no one there
			if totalClient <= 0 {
				fmt.Println("[System] Delete ", roomUUID)
				delete(webrtc.Rooms, roomUUID)
			}
		}

		// fmt.Println("Total Room created: ", len(webrtc.Rooms))
	}
}
