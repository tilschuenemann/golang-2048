package gameboard

import (
	"fmt"
	"math/rand"
)

type GameBoardModel struct {
	Gb           [4][4]int
	Score        int
	HasWon       bool
	IsMergable   bool
	HasEmptyTile bool
}

// New creates a new GameBoardModel.
func New() *GameBoardModel {

	gbm := GameBoardModel{
		Score:        0,
		HasWon:       false,
		IsMergable:   true,
		HasEmptyTile: true,
	}

	for y := 0; y <= 3; y++ {
		for x := 0; x <= 3; x++ {
			gbm.Gb[y][x] = 0
		}
	}

	for i := 0; i <= 1; i++ {
		x := rand.Intn(3)
		y := rand.Intn(3)
		gbm.Gb[x][y] = 2 << rand.Intn(2)
	}

	return &gbm
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

func (gbm *GameBoardModel) Move(input rune) (bool, bool) {

	validInput := false
	validRunes := []rune{'u', 'd', 'l', 'r'}
	for _, r := range validRunes {
		if input == r {
			validInput = true
			break
		}
	}

	if !validInput {
		return false, validInput
	}

	var tilt bool
	var reverse bool

	if input == 'u' || input == 'd' {
		tilt = true
	}

	if input == 'l' || input == 'd' {
		reverse = true
	}

	if tilt {
		gbm.tilt(false)
	}
	shifted := gbm.shiftHorizontal(reverse)
	merged := gbm.mergeAdjacentHorizontal(reverse)
	shifted_again := gbm.shiftHorizontal(reverse)

	if tilt {
		gbm.tilt(true)
	}

	return shifted || merged || shifted_again, validInput

}

func (gbm *GameBoardModel) shiftHorizontal(reverse bool) bool {

	var validMove bool

	end := 3
	stop := 0
	op := -1

	if reverse {
		end = 0
		stop = 3
		op = 1
	}

	for y := 0; y <= 3; y++ {
		// shift 0 to left
		for i := 0; i <= 2; i++ {
			for x := end; x != stop; x = x + op {
				if gbm.Gb[y][x] == 0 && x != stop {
					gbm.Gb[y][x] = gbm.Gb[y][x+op]
					gbm.Gb[y][x+op] = 0

					validMove = true
				}
			}
		}

	}

	return validMove
}

// tilt tilts the gameboard to the right.
func (gbm *GameBoardModel) tilt(reverse bool) {
	var tGB [4][4]int

	for y := 0; y <= 3; y++ {
		for x := 0; x <= 3; x++ {
			if reverse {
				tGB[x][y] = gbm.Gb[y][3-x]
			} else {
				tGB[x][3-y] = gbm.Gb[y][x]
			}
		}
	}
	gbm.Gb = tGB
}

func (gbm *GameBoardModel) mergeAdjacentHorizontal(reverse bool) bool {

	var validMove bool

	end := 3
	stop := 0
	op := -1

	if reverse {
		end = 0
		stop = 3
		op = 1
	}

	for y := 0; y <= 3; y++ {
		for x := end; x != stop; x = x + op {
			if gbm.Gb[y][x] == gbm.Gb[y][x+op] {
				gbm.Score += gbm.Gb[y][x+op]

				gbm.Gb[y][x+op] = gbm.Gb[y][x+op] << 1
				gbm.Gb[y][x] = 0

				if gbm.Gb[y][x+op] == 2048 {
					gbm.HasWon = true
				}

				validMove = true

			}
		}
	}

	return validMove
}

func (gbm *GameBoardModel) AddNewTile() {
	var emptyTiles []*int
	value := 1 << 2
	for y := 0; y <= 3; y++ {
		for x := 0; x <= 3; x++ {
			if gbm.Gb[y][x] == 0 {
				emptyTiles = append(emptyTiles, &gbm.Gb[y][x])
			}
		}
	}

	if len(emptyTiles) == 0 {
		gbm.HasEmptyTile = false
		return
	}

	*emptyTiles[rand.Intn(len(emptyTiles))] = value
	gbm.HasEmptyTile = true
}

// Print prints the current game.
func (gbm *GameBoardModel) Print() {
	for y := 0; y <= 3; y++ {
		for x := 0; x <= 3; x++ {
			fmt.Printf("| %4d ", gbm.Gb[y][x])
		}

		if y == 3 {
			fmt.Printf("| Score: %6d | ", gbm.Score)
		} else {
			fmt.Print("|               |")
		}
		fmt.Println()
	}
}

// CheckHas2048 checks if the gameboard has a 2048 tile.
// func (gbm *GameBoardModel) CheckHas2048() {
// 	for y := 0; y <= 3; y++ {
// 		for x := 0; x <= 3; x++ {
// 			if gbm.Gb[y][x] == 2048 {
// 				gbm.HasWon = true
// 				return
// 			}

// 		}
// 	}

// }

// checkIsMergable checks the gameboard for adjacent tiles that can be merged.
func (gbm *GameBoardModel) CheckIsMergable() {
	for y := 0; y <= 2; y++ {
		for x := 0; x <= 2; x++ {
			if gbm.Gb[y][x] == gbm.Gb[y][x+1] || gbm.Gb[y][x] == gbm.Gb[y+1][x] {
				gbm.IsMergable = true
				return
			}
		}
	}

	if gbm.Gb[3][3] == gbm.Gb[2][3] || gbm.Gb[3][3] == gbm.Gb[3][2] {
		gbm.IsMergable = true
		return
	}
	gbm.IsMergable = false

}
