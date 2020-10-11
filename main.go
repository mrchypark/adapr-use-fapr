package main

import (
	"errors"
	"os"
	"reflect"
	"strings"
	"time"

	"fapr/function"

	"github.com/rs/zerolog"
	"github.com/stoewer/go-strcase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/helmet/v2"
	jsoniter "github.com/json-iterator/go"
)

var version = "fapr-v0.1.0"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  version,
	})

	app.Use(helmet.New())
	app.Use(recover.New())
	app.Use(pprof.New())

	app.Use(logger.New(logger.Config{
		Format:     "{\"tid\": \"${header:tid}\", \"time\": \"${time}\", \"status\": \"${status}\", \"latency\": \"${latency}\", \"method\": \"${method}\", \"path\": \"${path}\", \n\"request_body\": \n${body}}\n",
		TimeFormat: "2006-01-02T15:04:05-0700",
		TimeZone:   "Asia/Seoul",
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))

	app.Get("/", showVersion)
	app.Get("/ver", showVersion)
	app.Get("/version", showVersion)
	app.Get("/health", healthCheck)
	app.Get("/healthz", healthCheck)

	app.Post("/*", endpoint)

	app.Listen(":4000")
}

func endpoint(c *fiber.Ctx) error {
	req := function.Input{}

	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return err
	}

	if err := requiredChk(req); err != nil {
		c.SendString(err.Error())
		return c.SendStatus(400)

	}

	return c.JSON(function.Handler(req))
}

func requiredChk(req interface{}) error {
	e := reflect.ValueOf(req)
	p := []string{}
	for i := 0; i < e.NumField(); i++ {
		v := e.Field(i)
		t := e.Type().Field(i)
		if v.Interface() == "" {
			p = append(p, strcase.SnakeCase(t.Name))
		}
	}
	if len(p) != 0 {
		return errors.New("Bad Request : '" + strings.Join(p, ", ") + "' is required.")
	}
	return nil
}

func healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"State": "ok",
	})
}

func showVersion(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"Version": version,
	})
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
