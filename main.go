package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/skip2/go-qrcode"
)

func main() {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(requestid.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,OPTIONS",
	}))

	app.Get("/", func(context *fiber.Ctx) error {
		return context.Status(fiber.StatusBadRequest).SendString(
			"Provide some text to be encoded as a QR code",
		)
	})

	app.Get("/:text", func(context *fiber.Ctx) error {
		textToBeEncoded := context.Params("text")

		qrcodePng, err := qrcode.Encode(textToBeEncoded, qrcode.Medium, 512)

		if err != nil {
			return context.Status(fiber.StatusInternalServerError).SendString(
				"ERROR: " + err.Error(),
			)
		}

		return context.Status(fiber.StatusOK).Send(qrcodePng)
	})

	app.Get("/:text/:size", func(context *fiber.Ctx) error {
		textToBeEncoded := context.Params("text")
		qrcodeImageSize, err := strconv.ParseUint(context.Params("size"), 10, 64)

		if err != nil {
			return context.Status(fiber.StatusInternalServerError).SendString(
				"ERROR: " + err.Error(),
			)
		}

		qrcodePng, err := qrcode.Encode(textToBeEncoded, qrcode.Medium, int(qrcodeImageSize))

		if err != nil {
			return context.Status(fiber.StatusInternalServerError).SendString(
				"ERROR: " + err.Error(),
			)
		}

		return context.Status(fiber.StatusOK).Send(qrcodePng)
	})

	err := app.Listen(":3000")

	if err != nil {
		panic("Error starting server: " + err.Error())
	}
}
