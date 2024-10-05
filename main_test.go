package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// Utils
func getResponse(t *testing.T, url string) *http.Response {
	app := App(fiber.Config{
		// 1 MB. The fuzzer sometimes gets too eager and sends a lot of data.
		ReadBufferSize: 1024 * 1024,
	})

	req := httptest.NewRequest("GET", url, nil)
	
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	return resp
}

func readResponse(t *testing.T, resp http.Response) string {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}

func getReadResponse(t *testing.T, url string) string {
	resp := getResponse(t, url)
	return readResponse(t, *resp)
}

// Int

func readInt(t *testing.T, body string) int {
	value, err := strconv.ParseInt(body, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	return int(value)
}

func getInt(t *testing.T, params string) int {
	body := getReadResponse(t, "/v1/int?"+params)
	return readInt(t, body)
}

func TestDefaultRange(t *testing.T) {
	value := getInt(t, "")
	if value < 1 || value > 100 {
		t.Fatalf("Expected value between 1 and 100, got %d", value)
	}
}

func TestMin(t *testing.T) {
	value := getInt(t, "min=80")
	if value < 80 || value > 100 {
		t.Fatalf("Expected value between 80 and 100, got %d", value)
	}
}

func TestMax(t *testing.T) {
	value := getInt(t, "max=300")
	if value < 1 || value > 300 {
		t.Fatalf("Expected value between 1 and 300, got %d", value)
	}
}

func TestMinAndMax(t *testing.T) {
	value := getInt(t, "min=200&max=300")
	if value < 200 || value > 300 {
		t.Fatalf("Expected value between 200 and 300, got %d", value)
	}
}

func FuzzInt(f *testing.F) {
	f.Add(0, 100)
	f.Add(20, 5)
	f.Add(10_000, 50_000)
	f.Add(-345, 921)
	f.Add(-500, -100)
	f.Add(-123, -345)
	
	f.Fuzz(func(t *testing.T, min int, max int) {
		resp := getResponse(t, fmt.Sprintf("/v1/int?min=%d&max=%d", min, max))
		body := readResponse(t, *resp)
		if resp.StatusCode != 200 {
			match, err := regexp.MatchString(`should be less than max`, string(body))
			if !match || err != nil {
				t.Fatalf("Unexpected response: %s\n%v", body, err)
			}
		} else {
			match, err := regexp.MatchString(`^-?\d+$`, string(body))
			if !match || err != nil {
				t.Fatalf("Unexpected response: %s\n%v", body, err)
			}
		}
	})
}

// Float
func readFloat(t *testing.T, body string) float64 {
	value, err := strconv.ParseFloat(body, 64)
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func getFloat(t *testing.T, params string) float64 {
	body := getReadResponse(t, "/v1/float?"+params)
	return readFloat(t, body)
}

func TestDefaultFloatRange(t *testing.T) {
	value := getFloat(t, "")
	if value < 0 || value > 1 {
		t.Fatalf("Expected value between 0 and 1, got %f", value)
	}
}

func TestFloatMinMax(t *testing.T) {
	value := getFloat(t, "min=-1.5&max=2.5")
	if value < -1.5 || value > 2.5 {
		t.Fatalf("Expected value between -1.5 and 2.5, got %f", value)
	}
}

func FuzzFloat(f *testing.F) {
	f.Add(0.0, 1.0)
	f.Add(-10.5, 10.5)
	f.Add(100.0, -2.0)
	f.Add(-4.2136, -3.2136)
	f.Add(2.0, 1.0)
	
	f.Fuzz(func(t *testing.T, min float64, max float64) {
		resp := getResponse(t, fmt.Sprintf("/v1/float?min=%f&max=%f", min, max))
		body := readResponse(t, *resp)
		if resp.StatusCode != 200 {
			match, err := regexp.MatchString(`should be less than max`, string(body))
			if !match || err != nil {
				t.Fatalf("Unexpected response: %s\n%v", body, err)
			}
		} else {
			match, err := regexp.MatchString(`^-?\d+(\.\d+)?$`, string(body))
			if !match || err != nil {
				t.Fatalf("Unexpected response: %s\n%v", body, err)
			}
		}
	})
}

// Word
func getWord(t *testing.T, params string) string {
	return getReadResponse(t, "/v1/word?"+params)
}

func TestDefaultWord(t *testing.T) {
	word := getWord(t, "")
	if word == "" {
		t.Fatal("Expected a non-empty word")
	}
}

func TestWordCategory(t *testing.T) {
	categories := []string{"animals", "cities", "countries", "fruits", "vegetables", "lorem-ipsum", "nouns"}
	for _, category := range categories {
		word := getWord(t, "category="+category)
		if word == "" {
			t.Fatalf("Expected a non-empty word for category %s, got %s", category, word)
		}
	}
}

func TestWordSeparator(t *testing.T) {
	words := getWord(t, "count=5&separator=,")
	
	if strings.Count(words, ",") != 4 {
		t.Fatalf("Expected 4 commas, got output %s", words)
	}
}

func TestWordCount(t *testing.T) {
	words := getWord(t, "count=5")
	
	if strings.Count(words, " ") != 4 {
		t.Fatalf("Expected 5 words, got output %s", words)
	}
}

// Dice
func getDice(t *testing.T, params string) string {
	return getReadResponse(t, "/v1/dice?"+params)
}

func TestDefaultDice(t *testing.T) {
	result := getDice(t, "")
	value := readInt(t, result)
	if value < 1 || value > 6 {
		t.Fatalf("Expected value between 1 and 6, got %d", value)
	}
}

func TestCustomDice(t *testing.T) {
	result := getDice(t, "input=2d20")
	value := readInt(t, result)
	if value < 2 || value > 40 {
		t.Fatalf("Expected value between 2 and 40, got %d", value)
	}
}

func TestDiceFullOutput(t *testing.T) {
	result := getDice(t, "input=3d6kh2&output=full")
	match, _ := regexp.MatchString(`^-?\d+ *(\[[-\d ]*\])? *(\(\[[-\d ]*\]\))?`, result)
	if !match {
		t.Fatalf("Unexpected dice output format: %s", result)
	}
}

func TestFudgeDice(t *testing.T) {
	result := getDice(t, "input=4df")
	value := readInt(t, result)
	if value < -4 || value > 4 {
		t.Fatalf("Expected value between -4 and 4, got %d", value)
	}
}

func TestFudgeDiceWithModifier(t *testing.T) {
	result := getDice(t, "input=" + url.QueryEscape("4df+10"))
	value := readInt(t, result)
	if value < 6 || value > 14 {
		t.Fatalf("Expected value between 6 and 14, got %d", value)
	}
}

func FuzzDice(f *testing.F) {
	f.Add("1d6", "sum")
	f.Add("2d20", "sum")
	f.Add("3d6kh2", "sum")
	f.Add("4d10kl3", "full")
	f.Add("38d12", "full")
	f.Add("100d100", "sum")
	f.Add("4df", "invalid")
	f.Add("4df+2", "full")
	f.Add("4df-2", "sum")

	f.Fuzz(func(t *testing.T, input string, output string) {
		resp := getResponse(t, fmt.Sprintf("/v1/dice?input=%s&output=%s", url.QueryEscape(input), url.QueryEscape(output)))
		body := readResponse(t, *resp)
		if resp.StatusCode != 200 {
			match, err := regexp.MatchString(`Bad roll format|Sides must be 1 or more|invalid input|invalid output|more dice than rolled|no result|Count must be 1 or more|Sides must be 2 or more`, string(body))
			if !match || err != nil {
				t.Fatalf("Unexpected response: %s\n%v", body, err)
			}
		} else {
					match, err := regexp.MatchString(`^-?\d+ *(\[[-\d ]*\])? *(\(\[[-\d ]*\]\))?`, string(body))
			if !match || err != nil {
				t.Fatalf("Unexpected response: %s\n%v", body, err)
			}
		}
	})
}

// ULID
func getULID(t *testing.T) string {
	return getReadResponse(t, "/v1/ulid")
}

func TestULID(t *testing.T) {
	ulid := getULID(t)
	match, _ := regexp.MatchString(`^[0-9A-Z]{26}$`, ulid)
	if !match {
		t.Fatalf("Invalid ULID format: %s", ulid)
	}
}

// NanoID
func getNanoID(t *testing.T, params string) string {
	return getReadResponse(t, "/v1/nanoid?"+params)
}

func TestDefaultNanoID(t *testing.T) {
	nanoid := getNanoID(t, "")
	if len(nanoid) != 21 {
		t.Fatalf("Expected NanoID of length 21, got %d", len(nanoid))
	}
}

func TestCustomSizeNanoID(t *testing.T) {
	nanoid := getNanoID(t, "size=30")
	if len(nanoid) != 30 {
		t.Fatalf("Expected NanoID of length 30, got %d", len(nanoid))
	}
}

func TestInvalidSizeNanoID(t *testing.T) {
	response := getResponse(t, "/v1/nanoid?size=0")
	if response.StatusCode != 400 {
		t.Fatalf("Expected status code 400, got %d", response.StatusCode)
	}
	response = getResponse(t, "/v1/nanoid?size=201")
	if response.StatusCode != 400 {
		t.Fatalf("Expected status code 400, got %d", response.StatusCode)
	}
}

// UUID
func getUUID(t *testing.T, params string) string {
	return getReadResponse(t, "/v1/uuid?"+params)
}

var uuidPattern = `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`

func TestDefaultUUID(t *testing.T) {
	uuid := getUUID(t, "")
	match, _ := regexp.MatchString(uuidPattern, uuid)
	if !match {
		t.Fatalf("Invalid UUID v4 format: %s", uuid)
	}
}

func TestUUIDv4(t *testing.T) {
	uuid := getUUID(t, "version=4")
	match, _ := regexp.MatchString(uuidPattern, uuid)
	if !match {
		t.Fatalf("Invalid UUID v4 format: %s", uuid)
	}
}

func TestUUIDv7(t *testing.T) {
	uuid := getUUID(t, "version=7")
	match, _ := regexp.MatchString(uuidPattern, uuid)
	if !match {
		t.Fatalf("Invalid UUID v7 format: %s", uuid)
	}
}

func TestBadUuidVersion(t *testing.T) {
	resp := getResponse(t, "/v1/uuid?version=9")
	if resp.StatusCode != 400 {
		t.Fatalf("Expected status code 400, got %d", resp.StatusCode)
	}
	body := readResponse(t, *resp)
	if !strings.Contains(body, "Invalid UUID version") {
		t.Fatalf("Unexpected response: %s", body)
	}
}

func hasIntHeader(t *testing.T, resp *http.Response, header string) int {
	value := resp.Header.Get(header)
	if value == "" {
		t.Fatalf("Expected %s header, got %s", header, value)
	}
	
	val, err := strconv.Atoi(value)
	if err != nil {
		t.Fatalf("Invalid %s header: %s", header, value)
	}
	
	return val
}

// Rate limit
func TestRateLimit(t *testing.T) {
	app := App(fiber.Config{
		// 1 MB. The fuzzer sometimes gets too eager and sends a lot of data.
		ReadBufferSize: 1024 * 1024,
	})

	req := httptest.NewRequest("GET", "/v1/int", nil)
	req.Header.Set("X-Forwarded-For", "10.0.0.1")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
	}
	
	limit := hasIntHeader(t, resp, "X-RateLimit-Limit")
	remaining := hasIntHeader(t, resp, "X-RateLimit-Remaining")
	reset := hasIntHeader(t, resp, "X-RateLimit-Reset")

	if limit < 0 || remaining < 0 || reset < 0 {
		t.Fatalf("Invalid rate limit headers: %d %d %d", limit, remaining, reset)
	}

	for i := 0; i < limit; i++ {
		resp, err = app.Test(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != 200 && (resp.StatusCode != 429 || !strings.Contains(readResponse(t, *resp), "Too many requests")) {
			t.Fatalf("Expected status code 200 or 429, got %d", resp.StatusCode)
		}
	}
}
