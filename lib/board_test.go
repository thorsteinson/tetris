package lib

import (
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

	for i = 0; i < BOARD_HEIGHT; i++ {
		for j = 0; j < BOARD_WIDTH; j++ {
			if b.GetTile(i, j) != EMPTY {
				t.Errorf("GetTile returned non empty tile on new board at postion: (%v %v)", i, j)
			}
		}
	}
}
