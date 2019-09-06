package lib

import (
	"testing"
)

func TestActiveTetrominoMove(t *testing.T) {
	x, _ := NewActiveTet(nil).Move(LEFT).GetPos()
	if x != STARTING_X-1 {
		t.Error("Left movement didn't shift tetromino")
	}

	x, _ = NewActiveTet(nil).Move(RIGHT).GetPos()
	if x != STARTING_X+1 {
		t.Error("Right movement didn't shift tetromino")
	}

	_, y := NewActiveTet(nil).Move(UP).GetPos()
	if y != STARTING_Y+1 {
		t.Error("Up movement didn't shift tetromino")
	}

	_, y = NewActiveTet(nil).Move(DOWN).GetPos()
	if y != STARTING_Y-1 {
		t.Error("Down movement didn't shift tetromino")
	}
}

func TestListCoodinates(t *testing.T) {
	tet := NewActiveTet(NewTet(TET_LINE))

	// The line piece is offset by 1 in it's grid, and is resting
	// vertically.

	expectedSet := map[Position]bool{
		Position{STARTING_X + 1, STARTING_Y}:     true,
		Position{STARTING_X + 1, STARTING_Y - 1}: true,
		Position{STARTING_X + 1, STARTING_Y - 2}: true,
		Position{STARTING_X + 1, STARTING_Y - 3}: true,
	}
	foundSet := make(map[Position]bool)

	coords := tet.ListCoords()

	t.Logf("Following coordinates were listed: %v", coords)

	for _, c := range coords {
		foundSet[c] = true
	}

	for expected := range expectedSet {
		if !foundSet[expected] {
			t.Errorf("Expected %v, but it wasn't found", expected)
		}
	}
}
