package main

import (
	"fmt"
	"strconv"
)

// Deliberately increased fieldSize for more convenient checking the field
const fieldSize = 12

// Field is a game field
type Field [fieldSize][fieldSize]bool

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
	//f.print() // <-- For debug
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

// isDestroyed returns true if player hitted the last part of ship, false if doesn't
func (f *Field) isDestroyed(y, x byte, gameField *Field) bool {
	row, _ := strconv.Atoi(string(y))
	col, _ := strconv.Atoi(string(x))

	direction, k, l := f.GetOrientation(row, col)
	row, col = row+1, col+1
	if direction == true {
		for i := row - k; i <= row+l; i++ {
			if gameField[i][col] == false {
				return false
			}
		}
	} else {
		for i := col - k; i <= col+l; i++ {
			if gameField[row][i] == false {
				return false
			}
		}
	}
	return true
}

// GetStrickenShips returns StrickenShips
func (f *Field) GetStrickenShips(msg []byte, un string) StrickenShips {
	if !f.isHitted(msg[1], msg[3]) {
		return StrickenShips{Ambient: []string{string(msg)}}
	}
	if !f.isDestroyed(msg[1], msg[3], shots[un]) {
		return StrickenShips{Hitted: string(msg)}
	}

	row, _ := strconv.Atoi(string(msg[1]))
	col, _ := strconv.Atoi(string(msg[3]))

	direction, k, l := f.GetOrientation(row, col)
	s := StrickenShips{Hitted: string(msg)}

	if direction == true {
		for i := row - k - 1; i <= row+l+1; i++ {
			s.Ambient = append(s.Ambient, fmt.Sprintf("e%v-%v", i, col-1))
			s.Ambient = append(s.Ambient, fmt.Sprintf("e%v-%v", i, col+1))
		}
		s.Ambient = append(s.Ambient, fmt.Sprintf("e%v-%v", row-k-1, col))
		s.Ambient = append(s.Ambient, fmt.Sprintf("e%v-%v", row+l+1, col))
		for i := row - k; i <= row+l; i++ {
			f[i+1][col+1] = false
		}
	} else {
		for i := col - k - 1; i <= col+l+1; i++ {
			s.Ambient = append(s.Ambient, fmt.Sprintf("e%v-%v", row-1, i))
			s.Ambient = append(s.Ambient, fmt.Sprintf("e%v-%v", row+1, i))
		}
		s.Ambient = append(s.Ambient, fmt.Sprintf("e%v-%v", row, col-k-1))
		s.Ambient = append(s.Ambient, fmt.Sprintf("e%v-%v", row, col+l+1))
		for i := col - k; i <= col+l; i++ {
			f[row+1][i+1] = false
		}
	}
	return s
}

// GetOrientation returns orientation of ship: false if the ship is horizontal and
// true if the ship is vertical; i,j - positions of shift to left and right
func (f *Field) GetOrientation(y, x int) (bool, int, int) {
	row, col := y+1, x+1

	var i, j int
	if f[row][col-1] == true || f[row][col+1] == true {
		for i = 1; f[row][col-i] == true; i++ {
		}
		for j = 1; f[row][col+j] == true; j++ {
		}
		return false, i - 1, j - 1
	}
	for i = 1; f[row-i][col] == true; i++ {
	}
	for j = 1; f[row+j][col] == true; j++ {
	}
	return true, i - 1, j - 1
}

// isHitted returns true if player hit the ship, false if doesn't
func (f *Field) isHitted(y, x byte) bool {
	row, _ := strconv.Atoi(string(y))
	col, _ := strconv.Atoi(string(x))

	return f[row+1][col+1]
}

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
