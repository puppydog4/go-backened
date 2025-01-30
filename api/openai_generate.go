package api

import "github.com/gofiber/fiber/v2"

func OpenaiGenerate(app *fiber.App) {
	app.Post("/generate", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
}
