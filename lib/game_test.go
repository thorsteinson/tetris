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

func TestListPositions(t *testing.T) {
	type testCase struct {
		expPositions []Position
		shape        Shape
	}

	cases := []testCase{
		testCase{
			expPositions: []Position{
				Position{STARTING_X + 0, STARTING_Y + 0},
				Position{STARTING_X + 0, STARTING_Y - 1},
				Position{STARTING_X + 1, STARTING_Y + 0},
				Position{STARTING_X + 1, STARTING_Y - 1},
			},
			shape: TET_SQUARE,
		},
		testCase{
			expPositions: []Position{
				Position{STARTING_X + 0, STARTING_Y - 2},
				Position{STARTING_X + 1, STARTING_Y - 2},
				Position{STARTING_X + 1, STARTING_Y - 1},
				Position{STARTING_X + 2, STARTING_Y - 1},
			},
			shape: TET_S,
		},
		testCase{
			expPositions: []Position{
				Position{STARTING_X + 0, STARTING_Y - 1},
				Position{STARTING_X + 1, STARTING_Y - 1},
				Position{STARTING_X + 1, STARTING_Y - 2},
				Position{STARTING_X + 2, STARTING_Y - 2},
			},
			shape: TET_Z,
		},
		testCase{
			expPositions: []Position{
				Position{STARTING_X + 1, STARTING_Y + 0},
				Position{STARTING_X + 1, STARTING_Y - 1},
				Position{STARTING_X + 1, STARTING_Y - 2},
				Position{STARTING_X + 1, STARTING_Y - 3},
			},
			shape: TET_LINE,
		},
	}

	var expectedSet map[Position]bool
	var tet ActiveTetromino
	var foundPositions []Position

	for _, test := range cases {
		// Create the expected set
		expectedSet = make(map[Position]bool)
		for _, p := range test.expPositions {
			expectedSet[p] = true
		}

		tet = NewActiveTet(NewTet(test.shape))
		foundPositions = tet.ListPositions()

		t.Logf("Expected Positions for shape %v: %v", test.shape, test.expPositions)

		if len(foundPositions) != len(expectedSet) {
			t.Errorf("Number of found positions is not 4.")
		}

		for _, pos := range foundPositions {
			if !expectedSet[pos] {
				t.Errorf("Unexpected position: %v", pos)
			}
		}
	}
}
