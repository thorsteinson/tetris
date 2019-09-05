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

func TestClearBoard(t *testing.T) {
	b := &Board{}

	// Set every tile to a value
	for i := 0; i < BOARD_HEIGHT; i++ {
		for j := 0; j < BOARD_WIDTH; j++ {
			b.SetTile(C1, i, j)
		}
	}

	b.Clear()

	for i := 0; i < BOARD_HEIGHT; i++ {
		for j := 0; j < BOARD_WIDTH; j++ {
			if b.GetTile(i, j) != EMPTY {
				t.Errorf("Value not cleared at position: (%v, %v)", i, j)
			}
		}
	}
}
