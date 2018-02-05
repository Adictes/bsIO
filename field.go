package main

import (
	"fmt"
	"strconv"
)

// Deliberately increased fieldSize for more convenient checking the field
const fieldSize = 12

// Field is a game field
type Field [fieldSize][fieldSize]bool

/* Cell is a struct with 2 filds:
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

// Init initializes the game field
func (f *Field) Init() {
	for i := 0; i < fieldSize; i++ {
		for j := 0; j < fieldSize; j++ {
			f[i][j] = Cell{busy: false, access: true}
		}
	}
}*/

// IndicateCell indicates the cell on field
func (f *Field) IndicateCell(y, x byte) {
	row, _ := strconv.Atoi(string(y))
	col, _ := strconv.Atoi(string(x))

	row, col = row+1, col+1

	if f[row][col] == false {
		f[row][col] = true
	} else {
		f[row][col] = false
	}
	fmt.Println(f.CheckPositionOfShips())
	f.print() // <-- For debug
}

// Ships is keeping number of available ships
type Ships struct {
	SingleDecker int
	TwoDecker    int
	ThreeDecker  int
	FourDecker   int
}

// GetAvailableShips returns available ships as struct
func (f *Field) GetAvailableShips() Ships {
	ships := Ships{4, 3, 2, 1}
	var seenCells [12][12]bool
	var shipLength = 0
	// Пройдемся сначала по горизонтали
	for i := 1; i < fieldSize-1; i++ {
		for j := 1; j < fieldSize; j++ {
			if f[i][j] == true {
				if f[i][j-1] == true || f[i][j+1] == true {
					shipLength++
					seenCells[i][j] = true
				}
			} else if shipLength != 0 {
				ships.shrink(shipLength)
				shipLength = 0
			}
		}
	}
	// Теперь по вертикали
	for j := 1; j < fieldSize-1; j++ {
		for i := 1; i < fieldSize; i++ {
			if seenCells[i][j] == true {
				continue
			}
			if f[i][j] == true {
				shipLength++
			} else if shipLength != 0 {
				ships.shrink(shipLength)
				shipLength = 0
			}
		}
	}
	return ships
}

// shrink check length of ship and remove one
func (s *Ships) shrink(length int) {
	if length == 1 {
		s.SingleDecker--
	} else if length == 2 {
		s.TwoDecker--
	} else if length == 3 {
		s.ThreeDecker--
	} else if length == 4 {
		s.FourDecker--
	}
}

// // GetNotAccessibleCells returns not accessible cells, by searching through all field
// func (f *Field) GetNotAccessibleCells() (coords []string) {
// 	for i := 1; i < fieldSize-1; i++ {
// 		for j := 1; j < fieldSize-1; j++ {
// 			if f[i][j] == false {
// 				coords = append(coords, strconv.Itoa(i-1)+"-"+strconv.Itoa(j-1))
// 			}
// 		}
// 	}
// 	return coords
// }

// CheckPositionOfShips checks correctness of ships setting
func (f *Field) CheckPositionOfShips() bool {
	if (f.GetAvailableShips() != Ships{0, 0, 0, 0}) {
		return false
	}

	var seenCells [12][12]bool
	var shipLength = 0
	// Пройдемся сначала по горизонтали
	for i := 1; i < fieldSize-1; i++ {
		for j := 1; j < fieldSize; j++ {
			if f[i][j] == true {
				if f[i][j-1] == true || f[i][j+1] == true {
					shipLength++
					seenCells[i][j] = true
				}
			} else if shipLength > 4 {
				return false
			} else if shipLength != 0 {
				for k := j - shipLength - 1; k <= j; k++ {
					if f[i-1][k] == true || f[i+1][k] == true {
						return false
					}
				}
				shipLength = 0
			}
		}
	}
	// Теперь по вертикали
	for j := 1; j < fieldSize-1; j++ {
		for i := 1; i < fieldSize; i++ {
			if seenCells[i][j] == true {
				continue
			}
			if f[i][j] == true {
				shipLength++
			} else if shipLength > 4 {
				return false
			} else if shipLength != 0 {
				for k := i - shipLength - 1; k <= i; k++ {
					if f[k][j-1] == true || f[k][j+1] == true {
						return false
					}
				}
				shipLength = 0
			}
		}
	}
	return true
}

// print prints game field to console
func (f *Field) print() {
	fmt.Println("----------------------")
	fmt.Println("   0 1 2 3 4 5 6 7 8 9")
	for i := 1; i < fieldSize-1; i++ {
		fmt.Printf("%v: ", i-1)
		for j := 1; j < fieldSize-1; j++ {
			if f[i][j] == true {
				fmt.Print("X ")
			} else {
				fmt.Print("O ")
			}
		}
		fmt.Println()
	}
	fmt.Println("----------------------")
}
