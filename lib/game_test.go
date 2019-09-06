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

func TestCanMove(t *testing.T) {
	type testCase struct {
		// Seed the board with these positions
		boardPositions []Position
		dir            Direction // Direction to attempt to move
		tetPos         Position  // Position of the active tet to start from
		shape          Shape     // The shape of the tetromino
		expected       bool      // Whether this move should be allowed
		name           string    // Name of the test
	}

	tests := []testCase{
		{
			boardPositions: []Position{},
			dir: UP,
			tetPos: Position{-100, 0},
			shape: TET_LINE,
			expected: false,
			name: "Below X Range",
		},
		{
			boardPositions: []Position{},
			dir: UP,
			tetPos: Position{-100, 0},
			shape: TET_LINE,
			expected: false,
			name: "Above X Range",
		},
		{
			boardPositions: []Position{},
			dir: UP,
			tetPos: Position{0, 100},
			shape: TET_LINE,
			expected: false,
			name: "Above Y Range",
		},
		{
			boardPositions: []Position{},
			dir: UP,
			tetPos: Position{0, -100},
			shape: TET_LINE,
			expected: false,
			name: "Below Y Range",
		},
		{
			boardPositions: []Position{{4,4}, {5,4}, {6,4}, {7,4}},
			dir: UP,
			tetPos: Position{3, 8},
			shape: TET_LINE,
			expected: true,
			name: "Above a horizontal line",
		},
		{
			boardPositions: []Position{{4,4}, {5,4}, {6,4}, {7,4}},
			dir: LEFT,
			tetPos: Position{3, 8},
			shape: TET_LINE,
			expected: true,
			name: "Left of a horizontal line",
		},
		{
			boardPositions: []Position{{4,4}, {5,4}, {6,4}, {7,4}},
			dir: RIGHT,
			tetPos: Position{3, 8},
			shape: TET_LINE,
			expected: true,
			name: "Right of a horizontal line",
		},
		{
			boardPositions: []Position{{4,4}, {5,4}, {6,4}, {7,4}},
			dir: DOWN,
			tetPos: Position{3, 8},
			shape: TET_LINE,
			expected: false,
			name: "Intersecting a horizontal line",
		},
	}

	var board *Board = &Board{}
	var tet ActiveTetromino
	for _, test := range tests {
		for _, p := range test.boardPositions {
			board.SetTile(C1, p.x, p.y)
		}

		// Create tetronimo and alter position
		tet = NewActiveTet(NewTet(test.shape))
		tet.x = test.tetPos.x
		tet.y = test.tetPos.y

		if tet.CanMove(test.dir, board) != test.expected {
			t.Errorf("Failed test: %v", test.name)
		}

		// Reset the board for the next test
		board.Clear()
	}
}
