package epazote

import (
	"testing"
)

func TestColorRed(t *testing.T) {

	color := Red("@")

	if color != "\x1b[0;31m@\x1b[0;00m" {
		t.Errorf("Expected red got: %s", color)
	}
}

func TestColorGreen(t *testing.T) {

	color := Green("@")

	if color != "\x1b[0;32m@\x1b[0;00m" {
		t.Errorf("Expected green got: %s", color)
	}
}
