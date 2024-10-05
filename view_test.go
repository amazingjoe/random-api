package main

import (
	"testing"
)

func TestIndex(t *testing.T) {

	resp := getResponse(t, "/")
	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
	}
}
