package main

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/tilschuenemann/golang-2048/gameboard"
)

// Dimensions
// Single Tile (with borders): 6x3 (8x5)
// Grid of 4x4 tiles with borders included: 29x17 (7*4+1, 4*4+1)
// Menu: 17x11

const (
	GameBoardX = 1
	GameBoardY = 1
	MenuX      = 30
	MenuY      = 1
	TileWidth  = 7
	TileHeight = 4
	GridDim    = 1
)

func main() {

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.Clear()

	var splashIsHidden bool
	drawSplashScreen(s)

	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	// setup gameboard model
	gbm := gameboard.New()

	for !gbm.HasWon || gbm.IsMergable || gbm.HasEmptyTile {
		s.Show()

		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()

		case *tcell.EventKey:
			if !splashIsHidden {
				splashIsHidden = true
				s.Clear()
			}

			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
				return
			}

			moved, ok := gbm.Move(ev.Rune())
			if !moved || !ok {
				continue
			}

			gbm.CheckIsMergable()
			gbm.AddNewTile()

			drawGameBoard(s, gbm.Gb)
			drawMenu(s, gbm.Score, string(ev.Rune()))

		}
	}
}

func drawSplashScreen(s tcell.Screen) {

	splashTextRows := []string{
		"              _                        ___   ___  _  _   ___  ",
		"             | |                      |__ \\ / _ \\| || | / _ \\ ",
		"   __ _  ___ | | __ _ _ __   __ _ ______ ) | | | | || || (_) |",
		"  / _` |/ _ \\| |/ _` | '_ \\ / _` |______/ /| | | |__   _> _ < ",
		" | (_| | (_) | | (_| | | | | (_| |     / /_| |_| |  | || (_) |",
		"  \\__, |\\___/|_|\\__,_|_| |_|\\__, |    |____|\\___/   |_| \\___/ ",
		"   __/ |                     __/ |                            ",
		"  |___/                     |___/                             "}

	for y, text := range splashTextRows {
		drawText(s, 0, y, text)
	}
}

func drawText(s tcell.Screen, x, y int, text string) {
	row := y
	col := x
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, tcell.StyleDefault)
		col++
	}
}

func drawOutline(s tcell.Screen) {
	x1 := 0
	x2 := 100
	y1 := 0
	y2 := 100

	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, tcell.StyleDefault)
		s.SetContent(col, y2, tcell.RuneHLine, nil, tcell.StyleDefault)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, tcell.StyleDefault)
		s.SetContent(x2, row, tcell.RuneVLine, nil, tcell.StyleDefault)
	}

	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, tcell.StyleDefault)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, tcell.StyleDefault)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, tcell.StyleDefault)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, tcell.StyleDefault)
	}

}

// drawTile draws a 7x4 tile, starting from x,y with value onto the screen.
func drawTile(s tcell.Screen, x, y, value int) {

	// draw box
	for i := 1; i <= 6; i++ {
		s.SetContent(x+i, y, tcell.RuneHLine, nil, tcell.StyleDefault)
		s.SetContent(x+i, y+4, tcell.RuneHLine, nil, tcell.StyleDefault)

	}
	for i := 1; i <= 3; i++ {
		s.SetContent(x, y+i, tcell.RuneVLine, nil, tcell.StyleDefault)
		s.SetContent(x+7, y+i, tcell.RuneVLine, nil, tcell.StyleDefault)
	}

	// determine tile color
	// TODO choose a nicer color palette!
	var bgColor tcell.Color
	switch value {
	case 0:
		bgColor = tcell.ColorLightBlue
	case 1:
		bgColor = tcell.ColorLightCoral
	case 2:
		bgColor = tcell.ColorLightCyan
	case 4:
		bgColor = tcell.ColorLightGoldenrodYellow
	case 8:
		bgColor = tcell.ColorLightGray
	case 16:
		bgColor = tcell.ColorLightGreen
	case 32:
		bgColor = tcell.ColorLightPink
	case 64:
		bgColor = tcell.ColorLightSalmon
	case 128:
		bgColor = tcell.ColorLightSeaGreen
	case 256:
		bgColor = tcell.ColorLightSkyBlue
	case 512:
		bgColor = tcell.ColorLightSlateGray
	case 1024:
		bgColor = tcell.ColorLightSteelBlue
	case 2048:
		bgColor = tcell.ColorLightYellow

	}
	boxStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(bgColor)
	for oy := 1; oy <= 3; oy++ {
		for ox := 1; ox <= 6; ox++ {
			s.SetContent(x+ox, y+oy, ' ', nil, boxStyle)
		}
	}

	// draw tile value
	for i, r := range []rune(fmt.Sprintf("%4d", value)) {
		s.SetContent(x+i+2, y+2, r, nil, boxStyle)
	}

}

func drawGameBoard(s tcell.Screen, gb [4][4]int) {
	for y := 0; y <= 3; y++ {
		for x := 0; x <= 3; x++ {
			drawTile(s, GameBoardX+x*TileWidth, GameBoardY+y*TileHeight, gb[y][x])
		}

	}

}

// drawMenu draws the menu onto the screen.
func drawMenu(s tcell.Screen, score int, input string) {

	controls := []string{
		"Controls:",
		"[q] to quit",
		"",
		"[u] to move up",
		"[d] to move down",
		"[l] to move left",
		"[r] to move right",
		"",
		"Current Score:",
		fmt.Sprintf("%5d", score),
		"",
	}

	if input != "" {
		controls = append(controls, "Last input:")
		controls = append(controls, input)
	}

	for i, text := range controls {
		drawText(s, MenuX, MenuY+i, text)
	}

}
