package gameboard_test

import (
	"testing"

	"github.com/tilschuenemann/golang-2048/gameboard"
)

func TestNewGameboard(t *testing.T) {

	gbm := gameboard.New()

	if gbm.Score != 0 {
		t.Error("gbm score != 0")
	}

	non_zero_tile_count := 0
	for y := 0; y <= 3; y++ {
		for x := 0; x <= 3; x++ {
			if gbm.Gb[y][x] != 0 {
				non_zero_tile_count++
			}
		}
	}
	if non_zero_tile_count != 2 {
		t.Error("gbm.Gb has more than two non-zero tiles!")
	}

}

func TestState(t *testing.T) {

	gbm := gameboard.New()
	gbm.Gb = [4][4]int{{2, 4, 2, 4}, {8, 16, 8, 16}, {2, 4, 2, 4}, {8, 16, 8, 16}}

	isMergable := gbm.IsMergable()
	if isMergable {
		t.Error("gbm.Gb is mergable when it shouldn't be!")
	}

	gbm.HasEmptyTile()
}

func TestIsMergable(t *testing.T) {
	testcases := []struct {
		description string
		board       [4][4]int
		result      bool
	}{
		{"empty gb", [4][4]int{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}}, false},
		{"single tile", [4][4]int{{2, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}}, false},
		{"vertical merge", [4][4]int{{2, 2, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}}, true},
		{"horizontal merge", [4][4]int{{2, 0, 0, 0}, {2, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}}, true},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {

			gbm := gameboard.New()
			gbm.Gb = tc.board
			isMergable := gbm.IsMergable()

			if isMergable != tc.result {
				t.Errorf("got %t want %t", isMergable, tc.result)
			}
		})
	}
}

func TestEmptyTiles(t *testing.T) {

	testcases := []struct {
		description string
		board       [4][4]int
		result      bool
	}{
		{"full gb", [4][4]int{{2, 4, 2, 4}, {8, 16, 8, 16}, {2, 4, 2, 4}, {8, 16, 8, 16}}, false},
		{"empty gb", [4][4]int{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}}, true},
		{"random gb with empty", [4][4]int{{2, 4, 8, 0}, {2, 4, 8, 0}, {2, 4, 8, 0}, {2, 4, 8, 0}}, true},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			gbm := gameboard.New()
			gbm.Gb = tc.board
			hasEmptyTile := gbm.HasEmptyTile()

			if hasEmptyTile != tc.result {
				t.Errorf("got %t want %t", hasEmptyTile, tc.result)
			}
		})
	}

}

func Test2048(t *testing.T) {

	testcases := []struct {
		description string
		board       [4][4]int
		result      bool
	}{
		{"empty gb", [4][4]int{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}}, false},
		{"2048", [4][4]int{{2048, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}}, true},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			gbm := gameboard.New()
			gbm.Gb = tc.board
			has2048 := gbm.Has2048Tile()

			if has2048 != tc.result {
				t.Errorf("got %t want %t", has2048, tc.result)
			}
		})
	}

}

func TestAddTile(t *testing.T) {
	testcases := []struct {
		description string
		board       [4][4]int
		result      bool
	}{
		{"empty gb", [4][4]int{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}}, true},
		{"full gb", [4][4]int{{2, 2, 2, 2}, {2, 2, 2, 2}, {2, 2, 2, 2}, {2, 2, 2, 2}}, false},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			gbm := gameboard.New()
			gbm.Gb = tc.board
			gbm.AddTile()

			if (tc.board == gbm.Gb) == tc.result {
				t.Error("No new tile added!")
			}

		})
	}

}
