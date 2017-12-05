package main

import (
	"strings"
	"testing"
)

func TestStringIndex(t *testing.T) {
	const s, sep, want = "chicken", "ken", 5
	got := strings.Index(s, sep)
	if got != want {
		t.Errorf("Index(%q,%q) = %v; want %v", s, sep, got, want)
	}
}

func TestLongRunning(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
}