package gameboard

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

type State int

const (
	Ongoing State = iota
	Won
	Lost
)

type GameBoardModel struct {
	Gb    [4][4]int
	Score int
}

// New creates a new GameBoardModel.
func New() *GameBoardModel {

	gbm := GameBoardModel{
		Score: 0,
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

func (gbm *GameBoardModel) MoveWrapper(input tcell.Key) {

	if gbm.GetState() != Ongoing {
		return
	}
	moved, ok := gbm.Move(input)
	if !ok {
		return
	}

	if !moved {
		return
	}

	gbm.AddTile()

}

func (gbm *GameBoardModel) Move(input tcell.Key) (bool, bool) {
	var tilt, reverse bool
	switch input {
	case tcell.KeyUp:
		tilt = true
	case tcell.KeyDown:
		tilt = true
		reverse = true
	case tcell.KeyLeft:
		reverse = true
	case tcell.KeyRight:
		tilt = false
		reverse = false
	default:
		return false, false
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

	return shifted || merged || shifted_again, true

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
				validMove = true

			}
		}
	}

	return validMove
}

// AddTile adds a new tile to the gameboard, if possible.
func (gbm *GameBoardModel) AddTile() {
	var emptyTiles []*int
	value := 2 << rand.Intn(2)
	for y := 0; y <= 3; y++ {
		for x := 0; x <= 3; x++ {
			if gbm.Gb[y][x] == 0 {
				emptyTiles = append(emptyTiles, &gbm.Gb[y][x])
			}
		}
	}

	if len(emptyTiles) == 0 {
		return
	}

	*emptyTiles[rand.Intn(len(emptyTiles))] = value
}

func (gbm *GameBoardModel) Has2048Tile() bool {
	for y := 0; y <= 3; y++ {
		for x := 0; x <= 3; x++ {
			if gbm.Gb[y][x] == 2048 {
				return true
			}
		}
	}
	return false
}

func (gbm *GameBoardModel) HasEmptyTile() bool {
	for y := 0; y <= 3; y++ {
		for x := 0; x <= 3; x++ {
			if gbm.Gb[y][x] == 0 {
				return true
			}
		}
	}
	return false
}

// IsMergable checks the gameboard for adjacent tiles that can be merged.
func (gbm *GameBoardModel) IsMergable() bool {
	for y := 0; y <= 2; y++ {
		for x := 0; x <= 2; x++ {
			if gbm.Gb[y][x] != 0 && (gbm.Gb[y][x] == gbm.Gb[y][x+1] || gbm.Gb[y][x] == gbm.Gb[y+1][x]) {
				return true

			}
		}
	}

	if gbm.Gb[3][3] != 0 && (gbm.Gb[3][3] == gbm.Gb[2][3] || gbm.Gb[3][3] == gbm.Gb[3][2]) {
		return true
	}
	return false

}

// GetState returns the current State.
func (gbm *GameBoardModel) GetState() State {

	has2048 := gbm.Has2048Tile()
	if has2048 {
		return Won
	}

	hasEmptyTile := gbm.HasEmptyTile()
	isMergable := gbm.IsMergable()

	if !hasEmptyTile && !isMergable {
		return Lost
	}

	return Ongoing

}
