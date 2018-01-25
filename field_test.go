package main

import (
	"testing"
)

func TestGetAvailableShips(t *testing.T) {
	var f Field
	f.Init()

	expected := Ships{4, 3, 2, 1}
	if got := f.GetAvailableShips(); got != expected {
		t.Errorf("Got: %v , expected: %v", got, expected)
	}

	f.Init()

	f[1] = [fieldSize]Cell{
		{busy: false, access: true},
		{busy: true, access: false},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: true, access: false},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: true, access: false},
		{busy: true, access: false},
		{busy: true, access: false},
		{busy: false, access: true},
		{busy: false, access: true},
	}

	expected = Ships{2, 3, 1, 1}
	if got := f.GetAvailableShips(); got != expected {
		t.Errorf("Got: %v , expected: %v", got, expected)
	}

	f.Init()

	f[1] = [fieldSize]Cell{
		{busy: false, access: true},
		{busy: true, access: false},
		{busy: true, access: false},
		{busy: true, access: false},
		{busy: true, access: false},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
	}
	f[8] = [fieldSize]Cell{
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: true, access: false},
		{busy: false, access: true},
	}
	f[9] = [fieldSize]Cell{
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: true, access: false},
		{busy: false, access: true},
	}
	f[10] = [fieldSize]Cell{
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: true, access: false},
		{busy: false, access: true},
	}

	expected = Ships{4, 3, 1, 0}
	if got := f.GetAvailableShips(); got != expected {
		t.Errorf("Got: %v , expected: %v", got, expected)
	}

	f.Init()

	f[2] = [fieldSize]Cell{
		{busy: false, access: true},
		{busy: true, access: false},
		{busy: true, access: false},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: true, access: false},
		{busy: false, access: true},
	}
	f[3] = [fieldSize]Cell{
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: true, access: false},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: true, access: false},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
	}
	f[4] = [fieldSize]Cell{
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: true, access: false},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: true, access: false},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
	}
	f[5] = [fieldSize]Cell{
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: true, access: false},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: true, access: false},
		{busy: false, access: true},
		{busy: true, access: false},
		{busy: true, access: false},
		{busy: false, access: true},
	}
	f[6] = [fieldSize]Cell{
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: true, access: false},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
	}
	f[7] = [fieldSize]Cell{
		{busy: false, access: true},
		{busy: true, access: false},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: true, access: false},
		{busy: false, access: true},
	}
	f[8] = [fieldSize]Cell{
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: true, access: false},
		{busy: true, access: false},
		{busy: true, access: false},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: true, access: false},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
	}
	f[9] = [fieldSize]Cell{
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: true, access: false},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
	}
	f[10] = [fieldSize]Cell{
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: true, access: false},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
		{busy: false, access: true},
	}

	expected = Ships{0, 0, 0, 0}
	if got := f.GetAvailableShips(); got != expected {
		t.Errorf("Got: %v , expected: %v", got, expected)
	}

}
