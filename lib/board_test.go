package lib

import (
	"math/rand"
	"testing"
)

func TestGetTileEmptyBoard(t *testing.T) {
	var i, j int

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("GetTile panicked at indices: (%v, %v)", i, j)
		}
	}()

	b := &Board{}

	for y := 0; y < BOARD_HEIGHT; y++ {
		for x := 0; x < BOARD_WIDTH; x++ {
			if b.GetTile(i, j) != EMPTY {
				t.Errorf("GetTile returned non empty tile on new board at postion: (%v %v)", x, y)
			}
		}
	}
}

// Tests the property that setting a tile and then getting it gets the
// tile that was just set
func TestSetGetTile(t *testing.T) {
	b := &Board{}

	tiles := []TileColor{EMPTY, C1, C2, C3, C4, C5, C6, C7}

	var tile TileColor
	for y := 0; y < BOARD_HEIGHT; y++ {
		for x := 0; x < BOARD_WIDTH; x++ {
			// pick a random tile
			tile = tiles[rand.Intn(len(tiles))]
			t.Logf("Testing position: (%v, %v)", x, y)
			b.SetTile(tile, x, y)
			if found := b.GetTile(x, y); found != tile {
				t.Errorf("Tile was set to %v but retreived value was %v", tile, found)
			}
		}
	}
}

// Fills the board with the provided tile color
func fill(b *Board, t TileColor) {
	for y := 0; y < BOARD_HEIGHT; y++ {
		for x := 0; x < BOARD_WIDTH; x++ {
			b.SetTile(t, x, y)
		}
	}
}

func TestClearBoard(t *testing.T) {
	b := &Board{}

	fill(b, C1)

	b.Clear()

	for y := 0; y < BOARD_HEIGHT; y++ {
		for x := 0; x < BOARD_WIDTH; x++ {
			if b.GetTile(x, y) != EMPTY {
				t.Errorf("Value not cleared at position: (%v, %v)", x, y)
			}
		}
	}
}

func TestTetris(t *testing.T) {
	b := &Board{}

	// Sets teh board so that it's totally full, with just a few holes
	// that create a vertical line
	fill(b, C1)
	b.SetTile(EMPTY, 5, 5)
	b.SetTile(EMPTY, 5, 6)
	b.SetTile(EMPTY, 5, 7)

	linesCleared := b.Tetris()
	expLinesCleared := BOARD_HEIGHT - 3

	if linesCleared != expLinesCleared {
		t.Errorf("Incorrect number of lines cleared reported. Expected %v, found %v",
			expLinesCleared, linesCleared)
	}

	// Only the first 3 lines should have values present
	for y := 0; y < 3; y++ {
		var lineHasValue bool
		for x := 0; x < BOARD_WIDTH; x++ {
			if !b.IsEmpty(x, y) {
				lineHasValue = true
			}
		}

		if !lineHasValue {
			t.Error("One of the first 3 lines has no values present")
		}
	}

	// Ensure that all lines above 2 are totally empty due to tetris
	for y := 3; y < BOARD_HEIGHT; y++ {
		var lineNonEmpty bool
		for x := 0; x < BOARD_WIDTH; x++ {
			if !b.IsEmpty(x, y) {
				lineNonEmpty = true
			}
		}

		if lineNonEmpty {
			t.Errorf("Expected line to be empty but found value at line: %v", y)
		}
	}
}

func TestFullLines(t *testing.T) {
	b := &Board{}

	fill(b, C1)

	if len(b.FullLines()) != BOARD_HEIGHT {
		t.Error("Not all lines detected in totally filled board")
	}

	n := 10
	for x := 0; x < BOARD_WIDTH; x++ {
		b.SetTile(EMPTY, x, n)
	}

	fullLines := b.FullLines()

	// Every line EXCEPT 10 should now be full
	for _, lineIdx := range fullLines {
		if lineIdx == n {
			t.Error("Erased line detected as a full line")
		}
	}
}

func TestIsEmpty(t *testing.T) {
	b := &Board{}

	if !b.IsEmpty(0, 0) {
		t.Error("Position should be empty initially")
	}

	b.SetTile(1, 0, 0)

	if b.IsEmpty(0, 0) {
		t.Error("Position is not empty after setting")
	}
}
