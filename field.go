package main

import (
	"fmt"
	"strconv"
)

// Deliberately increased fieldSize for more convenient checking the field
const fieldSize = 12

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

// IndicateCell indicates the cell on field
func IndicateCell(y, x byte) {
	row, _ := strconv.Atoi(string(y))
	col, _ := strconv.Atoi(string(x))

	row, col = row+1, col+1

	if Field[row][col].isBusy() {
		Field[row][col] = Cell{false, true}
	} else {
		Field[row][col] = Cell{true, false}
	}
	//PrintField()  <-- For debug
}

// PrintField prints game field to console
func PrintField() {
	fmt.Println("----------------------")
	fmt.Println("   0 1 2 3 4 5 6 7 8 9")
	for i := 1; i < fieldSize-1; i++ {
		fmt.Printf("%v: ", i-1)
		for j := 1; j < fieldSize-1; j++ {
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
