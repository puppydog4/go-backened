package api

import (
	"github.com/gofiber/fiber/v2"
)

func App() {
	app := fiber.New()
	OpenaiGenerate(app)
	app.Listen(":3000")
}
