package main

import (
	"fmt"
	"strconv"
)

const fieldSize = 10

// Cell is a struct with 2 filds:
// busy: true if ship is already standing at this cell
// access: true if ship can be placed at this cell
type Cell struct {
	busy   bool
	access bool
}

func (c Cell) isBusy() bool {
	return c.busy
}

func (c Cell) isAccessible() bool {
	return c.access
}

// Field is a game field
var Field [fieldSize][fieldSize]Cell

// fieldInit initializes the game Field
func fieldInit() {
	for i := 0; i < fieldSize; i++ {
		for j := 0; j < fieldSize; j++ {
			Field[i][j] = Cell{busy: false, access: true}
		}
	}
}

// ActWithCell @TODO
func ActWithCell(y byte, x byte) {
	row, _ := strconv.Atoi(string(y))
	col, _ := strconv.Atoi(string(x))

	if Field[row][col].isBusy() {
		Field[row][col] = Cell{false, true}
	} else {
		Field[row][col] = Cell{true, false}
	}
	PrintField()
}

// PrintField prints game field
func PrintField() {
	fmt.Println("----------------------")
	fmt.Println("   0 1 2 3 4 5 6 7 8 9")
	for i := 0; i < fieldSize; i++ {
		fmt.Printf("%v: ", i)
		for j := 0; j < fieldSize; j++ {
			if Field[i][j].isBusy() {
				fmt.Print("X ")
			} else {
				fmt.Print("O ")
			}
		}
		fmt.Println()
	}
	fmt.Println("----------------------")
}
