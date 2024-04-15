package main

import (
	"flag"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/zhinea/umamigo-server/entity"
	"github.com/zhinea/umamigo-server/libs/database"
	"github.com/zhinea/umamigo-server/utils"
	"log"
)

func main() {
	cfgFilename := flag.String("config", utils.GetEnvPath(), "Config file path.")

	flag.Parse()
	// load configuration file
	utils.LoadConfig(*cfgFilename)

	app := fiber.New()
	app.Use(cors.New())

	database.Connect()

	app.Get("/", func(c fiber.Ctx) error {

		return c.SendString("Hello, World!")
	})

	app.Post("/api/send", func(c fiber.Ctx) error {

		// TODO: Implement check bot request

		payload := entity.RequestPayload{}
		
		if err := c.Bind().Body(&payload); err != nil {
			log.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request",
			})
		}

		log.Println("Payload: ", payload)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Success",
			"data":    payload.Events,
		})

	})

	err := app.Listen(":3022", fiber.ListenConfig{
		EnablePrefork: true,
	})
	if err != nil {
		log.Println("Error: ", err)
		return
	}
	log.Println("Server started on port 3022")
}
