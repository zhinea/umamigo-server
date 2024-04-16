package main

import (
	"flag"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zhinea/umamigo-server/entity"
	"github.com/zhinea/umamigo-server/libs/database"
	libjwt "github.com/zhinea/umamigo-server/libs/jwt"
	"github.com/zhinea/umamigo-server/libs/session"
	"github.com/zhinea/umamigo-server/utils"
	"log"
	"time"
)

func main() {
	cfgFilename := flag.String("config", utils.GetEnvPath(), "Config file path.")

	flag.Parse()
	// load configuration file
	utils.LoadConfig(*cfgFilename)

	app := fiber.New(fiber.Config{
		Immutable:   false,
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	app.Use(cors.New())
	//app.Use(func(c fiber.Ctx) error {
	//	// start timer
	//	start := time.Now()
	//	// next routes
	//	c.Next()
	//	// stop timer
	//	stop := time.Now()
	//	// Do something with response
	//	log.Println("Request took ", stop.Sub(start).String())
	//	return c.Next()
	//})

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))
	//app.Use(logger.New())

	database.Connect()

	app.Get("/", func(c fiber.Ctx) error {

		return c.JSON(fiber.Map{
			"halo":    "hall",
			"headers": c.GetReqHeaders()["x-umami-cache"],
		})
	})

	app.Post("/api/send", func(c fiber.Ctx) error {
		// TODO: Implement check bot request
		totalTime := time.Now()
		validationTimer := time.Now()

		payload := entity.RequestPayload{}

		if err := c.Bind().Body(&payload); err != nil {
			log.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request",
			})
		}

		log.Println("main: validation timer took", time.Now().Sub(validationTimer).String())

		sessionTimer := time.Now()
		headers := c.GetReqHeaders()

		sess, err := session.UseSession(&entity.UseSessionPayloadData{
			Headers: headers,
			Body:    payload,
			IP:      c.IP(),
			IsLocal: c.IsFromLocal(),
		})

		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Server cannot process the request",
				"e":     err,
			})
		}

		if sess.RegisteredClaims.IssuedAt != nil && time.Now().Unix()-sess.RegisteredClaims.IssuedAt.Unix() > 1800 {
			sess.VisitID = utils.UUID(sess.ID, utils.VisitSalt())
		}

		sess.RegisteredClaims.IssuedAt = &jwt.NumericDate{Time: time.Now()}

		log.Println("main: session timer took", time.Now().Sub(sessionTimer).String())

		jwtCreationTimer := time.Now()
		jwtToken := libjwt.CreateToken(sess)
		log.Println("main: jwt creation took", time.Now().Sub(jwtCreationTimer).String())
		log.Println("main: total time took", time.Now().Sub(totalTime).String())
		log.Println("============")
		return c.SendString(jwtToken)

	})

	err := app.Listen(":3022", fiber.ListenConfig{
		EnablePrefork: utils.Cfg.Server.Prefork,
	})
	if err != nil {
		log.Println("Error: ", err)
		return
	}
	log.Println("Server started on port 3022")
}
