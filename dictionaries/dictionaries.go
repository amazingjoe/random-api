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
package dictionaries

import (
	_ "embed"
	"strings"
)

//go:embed animals.txt
var animals string
//go:embed words.txt
var words string
//go:embed cities.txt
var cities string
//go:embed countries.txt
var countries string
//go:embed fruits.txt
var fruits string
//go:embed vegetables.txt
var vegetables string
//go:embed lorem-ipsum.txt
var loremIpsum string
//go:embed nouns.txt
var nouns string

var Dictionaries = map[string][]string{
	"animals": strings.Split(animals, "\n"),
	"words": strings.Split(words, "\n"),
	"cities": strings.Split(cities, "\n"),
	"countries": strings.Split(countries, "\n"),
	"fruits": strings.Split(fruits, "\n"),
	"vegetables": strings.Split(vegetables, "\n"),
	"lorem-ipsum": strings.Split(loremIpsum, "\n"),
	"nouns": strings.Split(nouns, "\n"),
}

func mapKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

var Keys = mapKeys(Dictionaries)
