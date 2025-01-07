package exception

import (
	"maze-conquest-api/model/web"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	code := 500
	message := "Internal Server Error"

	if e, ok := err.(*NotFoundError); ok {
		return ctx.Status(e.Code).JSON(web.WebResponse{
			Code:   e.Code,
			Status: "Not Found",
			Data: fiber.Map{
				"message": e.Error(),
			},
		})
	}

	if e, ok := err.(*EmptyUidError); ok {
		return ctx.Status(e.Code).JSON(web.WebResponse{
			Code:   e.Code,
			Status: "Empty Uid",
			Data: fiber.Map{
				"message": e.Error(),
			},
		})
	}

	// Fatal error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return ctx.Status(code).JSON(web.WebResponse{
		Code:   code,
		Status: message,
		Data: fiber.Map{
			"message": err.Error(),
		},
	})
}
