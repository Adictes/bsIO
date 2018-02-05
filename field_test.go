package main

import (
	"testing"
)

func TestGetAvailableShips(t *testing.T) {
	var f Field

	expected := Ships{4, 3, 2, 1}
	if got := f.GetAvailableShips(); got != expected {
		t.Errorf("Got: %v , expected: %v", got, expected)
	}

	f = Field{
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, true, false, false, true, false, false, true, true, true, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
	}

	expected = Ships{2, 3, 1, 1}
	if got := f.GetAvailableShips(); got != expected {
		t.Errorf("Got: %v , expected: %v", got, expected)
	}

	f = Field{
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, true, true, true, true, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, true, false},
		{false, false, false, false, false, false, false, false, false, false, true, false},
		{false, false, false, false, false, false, false, false, false, false, true, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
	}

	expected = Ships{4, 3, 1, 0}
	if got := f.GetAvailableShips(); got != expected {
		t.Errorf("Got: %v , expected: %v", got, expected)
	}

	f = Field{
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
		{false, true, true, false, false, false, false, false, false, false, true, false},
		{false, false, false, false, true, false, false, true, false, false, false, false},
		{false, false, false, false, true, false, false, true, false, false, false, false},
		{false, false, false, false, true, false, false, true, false, true, true, false},
		{false, false, false, false, true, false, false, false, false, false, false, false},
		{false, true, false, false, false, false, false, false, false, false, true, false},
		{false, false, false, true, true, true, false, false, true, false, false, false},
		{false, false, false, false, false, false, false, false, true, false, false, false},
		{false, false, true, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false},
	}

	expected = Ships{0, 0, 0, 0}
	if got := f.GetAvailableShips(); got != expected {
		t.Errorf("Got: %v , expected: %v", got, expected)
	}
}
