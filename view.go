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
	"embed"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

//go:embed dist/*
var dist embed.FS

func SetupViews() *fiber.App {
	app := fiber.New()
	app.Use("/", filesystem.New(filesystem.Config{
		Root: http.FS(dist),
		PathPrefix: "dist",
		Browse: false,
	}))
	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return true
		},
		CacheControl: true,
		Expiration: 60 * time.Second,
	}))
	
	return app
}
