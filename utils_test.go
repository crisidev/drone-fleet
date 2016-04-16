package main

import "testing"

func TestDeleteEmptyString(t *testing.T) {
	slice := []string{"", "a", "b"}
	slice = DeleteEmptyString(slice)
	if len(slice) != 2 {
		t.Fatalf("expected slice lenght == 2, got %d", len(slice))
	}
	slice = []string{"a", "b"}
	slice = DeleteEmptyString(slice)
	if len(slice) != 2 {
		t.Fatalf("expected slice lenght == 2, got %d", len(slice))
	}
}
