/*
Copyright (C) 2024  Kaan Barmore-Genc

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	_ "embed"
	"errors"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"github.com/justinian/dice"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/oklog/ulid/v2"
	"github.com/seriousbug/random/v2/dictionaries"
)

func randInt(c *fiber.Ctx) error {
	min := c.QueryInt("min", 0)
	max := c.QueryInt("max", 100)

	if min >= max {
		var err = c.SendStatus(400)
		if err != nil {
			return err
		}
		return c.SendString(fmt.Sprintf("min %d should be less than max %d", min, max))
	}

	offset := 0
	if min < 0 {
		offset = min
		min = 0
		max = max - offset
	}

	return c.SendString(fmt.Sprint(rand.IntN(max-min) + offset + min))
}

func randFloat(c *fiber.Ctx) error {
	min := c.QueryFloat("min", 0)
	max := c.QueryFloat("max", 1)

	if min >= max {
		var err = c.SendStatus(400)
		if err != nil {
			return err
		}
		return c.SendString(fmt.Sprintf("min %f should be less than max %f", min, max))
	}

	return c.SendString(fmt.Sprint(rand.Float64()*(max-min) + min))
}

func randWord(c *fiber.Ctx) error {
	category := c.Query("category", "words")
	count := c.QueryInt("count", 1)
	separator := c.Query("separator", " ")

	dict, ok := dictionaries.Dictionaries[category]
	if !ok {
		c.SendStatus(400)
		return c.SendString("Invalid category, must be one of " + strings.Join(dictionaries.Keys, ", "))
	}

	words := make([]string, 0, count)
	for i := 0; i < count; i++ {
		words = append(words, dict[rand.IntN(len(dict))])
	}

	return c.SendString(strings.Join(words, separator))
}

func randDice(c *fiber.Ctx) error {
	input := c.Query("input", "1d6")
	output := c.Query("output", "sum")

	if len(input) == 0 || len(input) > 100 {
		c.SendStatus(400)
		return c.SendString("Invalid input")
	}

	result, _, err := dice.Roll(input)
	if err != nil {
		c.SendStatus(400)
		return c.SendString(err.Error())
	}

	if output == "sum" {
		return c.SendString(fmt.Sprint(result.Int()))
	}
	if output == "full" {
		return c.SendString(result.String())
	}

	return errors.New("invalid output, must be sum or full")
}

func randUlid(c *fiber.Ctx) error {
	return c.SendString(ulid.Make().String())
}

func randNanoId(c *fiber.Ctx) error {
	size := c.QueryInt("size", 21)

	if size < 1 || size > 200 {
		c.SendStatus(400)
		return c.SendString("Invalid size, must be between 1 and 200")
	}

	id, err := gonanoid.ID(size)
	if err != nil {
		return err
	}
	return c.SendString(id)
}

func randUuid(c *fiber.Ctx) error {
	version := c.Query("version", "4")

	var id uuid.UUID

	switch version {
	case "4":
		var err error
		id, err = uuid.NewRandom()
		if err != nil {
			return err
		}
	case "7":
		var err error
		id, err = uuid.NewV7()
		if err != nil {
			return err
		}
	default:
		c.SendStatus(400)
		return c.SendString("Invalid UUID version, must be 4 or 7")
	}
	return c.SendString(id.String())
}

func main() {
	app := fiber.New()

	app.Use(recover.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Use(limiter.New(limiter.Config{
		Max:        180,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for", c.IP())
		},
		LimitReached: func(c *fiber.Ctx) error {
			c.SendStatus(429)
			return c.SendString("Too many requests")
		},
	}))

	v1 := app.Group("/v1")
	v1.Get("/int", randInt)
	v1.Get("/float", randFloat)
	v1.Get("/word", randWord)
	v1.Get("/dice", randDice)
	v1.Get("/ulid", randUlid)
	v1.Get("/nanoid", randNanoId)
	v1.Get("/uuid", randUuid)

	app.Mount("/", SetupViews())

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
