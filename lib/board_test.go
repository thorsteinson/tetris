package lib

import (
	"math/rand"
	"testing"
	"time"
)

func TestGetTileEmptyBoard(t *testing.T) {
	var i, j int

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("GetTile panicked at indices: (%v, %v)", i, j)
		}
	}()

	b := &Board{}

	for i = 0; i < BOARD_HEIGHT; i++ {
		for j = 0; j < BOARD_WIDTH; j++ {
			if b.GetTile(i, j) != EMPTY {
				t.Errorf("GetTile returned non empty tile on new board at postion: (%v %v)", i, j)
			}
		}
	}
}

// Tests the property that setting a tile and then getting it gets the
// tile that was just set
func TestSetGetTile(t *testing.T) {
	b := &Board{}
	seed := time.Now().UnixNano()
	rand.Seed(seed)

	t.Logf("Using seed: %v", seed)

	tiles := []TileColor{EMPTY, C1, C2, C3, C4, C5, C6, C7}

	var tile TileColor
	for i := 0; i < BOARD_HEIGHT; i++ {
		for j := 0; j < BOARD_WIDTH; j++ {
			// pick a random tile
			tile = tiles[rand.Intn(len(tiles))]
			b.SetTile(tile, i, j)
			if found := b.GetTile(i, j); found != tile {
				t.Errorf("Tile was set to %v but retreived value was %v", tile, found)
			}
		}
	}
}

// Fills the board with the provided tile color
func fill(b *Board, t TileColor) {
	for i := 0; i < BOARD_HEIGHT; i++ {
		for j := 0; j < BOARD_WIDTH; j++ {
			b.SetTile(t, i, j)
		}
	}
}

func TestClearBoard(t *testing.T) {
	b := &Board{}

	fill(b, C1)

	b.Clear()

	for i := 0; i < BOARD_HEIGHT; i++ {
		for j := 0; j < BOARD_WIDTH; j++ {
			if b.GetTile(i, j) != EMPTY {
				t.Errorf("Value not cleared at position: (%v, %v)", i, j)
			}
		}
	}
}

func TestEraseLine(t *testing.T) {
	b := &Board{}

	fill(b, C1)

	i := 3

	b.EraseLine(i)

	for j := 0; j < BOARD_WIDTH; j++ {
		if b.GetTile(i, j) != EMPTY {
			t.Error("Non empty value found in erased line")
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
	b.EraseLine(n)

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
