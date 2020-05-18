package main

import "testing"

func TestHello(t *testing.T) {
	got := Hello()
	expected := "Hello, world\n"

	if got != expected {
		t.Errorf("Expected:\n%q\nGot:\n%q", expected, got)
	}
}
