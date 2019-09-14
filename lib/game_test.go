package lib

import (
	"math/rand"
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
			dir:            UP,
			tetPos:         Position{-100, 0},
			shape:          TET_LINE,
			expected:       false,
			name:           "Below X Range",
		},
		{
			boardPositions: []Position{},
			dir:            UP,
			tetPos:         Position{-100, 0},
			shape:          TET_LINE,
			expected:       false,
			name:           "Above X Range",
		},
		{
			boardPositions: []Position{},
			dir:            UP,
			tetPos:         Position{0, 100},
			shape:          TET_LINE,
			expected:       false,
			name:           "Above Y Range",
		},
		{
			boardPositions: []Position{},
			dir:            UP,
			tetPos:         Position{0, -100},
			shape:          TET_LINE,
			expected:       false,
			name:           "Below Y Range",
		},
		{
			boardPositions: []Position{{4, 4}, {5, 4}, {6, 4}, {7, 4}},
			dir:            UP,
			tetPos:         Position{3, 8},
			shape:          TET_LINE,
			expected:       true,
			name:           "Above a horizontal line",
		},
		{
			boardPositions: []Position{{4, 4}, {5, 4}, {6, 4}, {7, 4}},
			dir:            LEFT,
			tetPos:         Position{3, 8},
			shape:          TET_LINE,
			expected:       true,
			name:           "Left of a horizontal line",
		},
		{
			boardPositions: []Position{{4, 4}, {5, 4}, {6, 4}, {7, 4}},
			dir:            RIGHT,
			tetPos:         Position{3, 8},
			shape:          TET_LINE,
			expected:       true,
			name:           "Right of a horizontal line",
		},
		{
			boardPositions: []Position{{4, 4}, {5, 4}, {6, 4}, {7, 4}},
			dir:            DOWN,
			tetPos:         Position{3, 8},
			shape:          TET_LINE,
			expected:       false,
			name:           "Intersecting a horizontal line",
		},
		{
			boardPositions: []Position{},
			dir:            DOWN,
			tetPos:         Position{0, 3},
			shape:          TET_LINE,
			expected:       false,
			name:           "Going under the board",
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

func TestBoardControllerNextTet(t *testing.T) {
	board := &Board{}

	ctl := NewBoardController(board, NewTet(TET_LINE))

	// Check that the active tet has a line shape
	if ctl.tet.shape != TET_LINE {
		t.Error("Tetromino has non line shape")
	}

	// Check that none of the tiles are empty, they should have some value
	for _, p := range ctl.tet.ListPositions() {
		if ctl.board.IsEmpty(p.x, p.y) {
			t.Errorf("Unexpected empty tile: %v", p)
		}
	}

	// Look at the positions of the tetromino and compare that to the
	// board. If these aren't set, we've got a problem
	expectedTC := ShapeToTC(ctl.tet.shape)
	for _, p := range ctl.tet.ListPositions() {
		if tc := ctl.board.GetTile(p.x, p.y); tc != expectedTC {
			t.Errorf("Incorrect TileColor: expected %v, found %v", expectedTC, tc)
		}
	}

	// Now slam the tet to the bottom, and get the next one
	ctl.Slam()
	ctl.NextTet(NewTet(TET_SQUARE))

	// Check that there are now a total of 8 tiles in the board
	var tileCount int
	for _, tc := range ctl.board.tiles {
		if tc != EMPTY {
			tileCount++
		}
	}

	if tileCount != 8 {
		t.Error("Unexpected number of tiles found")
	}
}

func TestBoardControllerMove(t *testing.T) {
	const TRIALS = 20
	const MOVEMENTS = 100

	endPositions := make(map[Position]bool)

	trial := func() {
		board := &Board{}

		ctl := NewBoardController(board, NewTet(TET_LINE))

		// Move randomly 100 times. If our movement code is safe, then it
		// should end up just fine without crashing. It's also highly
		// likely it doesn't end up in the same spot it started
		var dir Direction
		for i := 0; i < MOVEMENTS; i++ {
			dir = Direction(rand.Intn(4))
			ctl.Move(dir)
		}

		// Return if the starting and ending position are the SAME
		endPositions[ctl.tet.Position] = true
	}

	for i := 0; i < TRIALS; i++ {
		trial()
	}

	if len(endPositions) < 5 {
		t.Errorf("After %v trials the tetromino ended up in %v unique locations. This is too low",
			TRIALS,
			len(endPositions))
	}
}

func TestSlam(t *testing.T) {
	type testCase struct {
		boardPositions []Position
		shape          Shape
		expectedYMin   int
	}

	tests := []testCase{
		{
			boardPositions: []Position{},
			shape:          TET_LINE,
			expectedYMin:   0,
		},
		{
			boardPositions: []Position{{2, 3}, {3, 3}, {4, 3}, {5, 3}, {6, 3}, {7, 3}},
			shape:          TET_LINE,
			expectedYMin:   4,
		},
	}

	for _, test := range tests {
		// Setup the board controller
		board := &Board{}
		ctl := NewBoardController(board, NewTet(test.shape))

		for _, p := range test.boardPositions {
			ctl.board.SetTile(ShapeToTC(test.shape), p.x, p.y)
		}

		// Slam the tetromino
		ctl.Slam()

		// Check that the minimum y-height of the tile matches what we expect
		var yMin = BOARD_HEIGHT - 1
		for _, p := range ctl.tet.ListPositions() {
			if p.y < yMin {
				yMin = p.y
			}
		}

		if yMin != test.expectedYMin {
			t.Errorf("Expected Y-min %v, found %v", test.expectedYMin, yMin)
		}
	}
}

// Ensure the tetris effect is applied when NextTet is called
func TestNextTetTetris(t *testing.T) {
	board := &Board{}
	ctl := NewBoardController(board, NewTet(TET_T))

	// Set all but one tile in the bottom of a board to non empty
	for x := 0; x < BOARD_WIDTH; x++ {
		board.SetTile(C1, x, 0)
	}
	// Clear the middle tile
	board.SetTile(EMPTY, 5, 0)

	// Slam that T-Piece right down into the center, filling that hole
	ctl.Slam()
	ctl.NextTet(NewTet(TET_T))

	// There should only be 3 tiles in the entire board
	var tileCount int
	var tile TileColor
	for y := 0; y < BOARD_HEIGHT; y++ {
		for x := 0; x < BOARD_WIDTH; x++ {
			tile = board.GetTile(x, y)
			if tile != EMPTY {
				t.Logf("Tile with color %v found at postion (%v, %v)", tile, x, y)
				tileCount++
			}
		}
	}

	// There should be three tiles remaining, after the tetris, and
	// the tiles for the next tetromino, so a total of 7
	const EXP_COUNT = 7
	if tileCount != EXP_COUNT {
		t.Errorf("Unexpected number of tiles. Expected %v, found: %v", EXP_COUNT, tileCount)
	}

	// These specified positions should have the TET_T color
	// associated with them
	ps := []Position{{4, 0}, {5, 0}, {6, 0}}
	expectedColor := ShapeToTC(TET_T)
	for _, p := range ps {
		if color := board.GetTile(p.x, p.y); color != expectedColor {
			t.Errorf("Unexpected color at position %v, found %v, expected %v", p, color, expectedColor)
		}
	}
}

func TestBoardControllerRotation(t *testing.T) {
	// Setup
	board := &Board{}
	ctl := NewBoardController(board, NewTet(TET_LINE))

	// Move our line peice all the way to the right, as far as it will
	// go
	for ctl.tet.CanMove(RIGHT, ctl.board) {
		ctl.Move(RIGHT)
	}

	// Now attempt a rotation. Since this is on the edge of the board,
	// it should force our position to change
	pInit := ctl.tet.Position

	ctl.RotLeft()

	pEnd := ctl.tet.Position

	if pInit == pEnd {
		t.Errorf("Start and end positions are same: %v, expected them to be different", pInit)
	}

	// Move it towards the center and set it vertical
	ctl.Move(LEFT)
	ctl.Move(LEFT)
	ctl.Move(LEFT)
	ctl.Move(LEFT)
	ctl.Move(LEFT)
	ctl.RotRight()

	// Set the tiles to the left and right of the tetromino, to ensure
	// it can't rotate while boxed
	for _, p := range ctl.tet.ListPositions() {
		ctl.board.SetTile(C1, p.x+1, p.y)
		ctl.board.SetTile(C1, p.x-1, p.y)
	}

	initPositions := ctl.tet.ListPositions()
	ctl.RotLeft()
	for i, p := range ctl.tet.ListPositions() {
		if initPositions[i] != p {
			t.Errorf("Tetromino rotated when surrounded by tiles. Expected %v, found %v", initPositions[i], p)
		}
	}
}

// Simulates a ton of random movements, just to ensure that nothing we
// do could possibly end up with us out of bounds or breaking the game
func TestRandomWalkStressTest(t *testing.T) {
	// Setup
	board := &Board{}
	source := make(chan *Tetromino, 10)
	for _, s := range shapes {
		source <- NewTet(s)
	}
	// We need to queue up an extra shape, or we'll deadlock
	source <- NewTet(TET_SQUARE)
	ctl := NewBoardController(board, NewTet(TET_SQUARE))

	const MOVEMENTS = 200

	var movement int
	for _, s := range shapes {
		t.Logf("Testing random walk stress test with shape: %v", s)

		for i := 0; i < MOVEMENTS; i++ {
			t.Logf("Movement %v of %v", i+1, MOVEMENTS)
			movement = rand.Intn(6)
			t.Logf("Next Movement: %v", movement)
			if movement < 4 {
				ctl.Move(Direction(movement))
			} else if movement == 4 {
				t.Logf("Current position: %v", ctl.tet.Position)
				t.Log("Current tile positions")
				t.Logf("%v", ctl.tet.ListPositions())
				ctl.RotLeft()
			} else {
				ctl.RotRight()
			}
		}

		ctl.Slam()
		ctl.NextTet(NewTet(s))
	}
}

func TestGameOver(t *testing.T) {
	// Setup
	board := &Board{}
	// We need to queue up an extra shape, or we'll deadlock
	ctl := NewBoardController(board, NewTet(TET_SQUARE))

	if ctl.isGameover {
		t.Error("Gameover signal detected too early")
	}

	ctl.NextTet(NewTet(TET_SQUARE))

	if !ctl.isGameover {
		t.Error("Gameover signal not found")
	}
}

func TestNaturalGameover(t *testing.T) {
	// Setup
	board := &Board{}
	ctl := NewBoardController(board, NewTet(TET_LINE))

	for i := 0; i < 4; i++ {
		ctl.Slam()
		ctl.NextTet(NewTet(TET_LINE))
		if ctl.isGameover {
			t.Error("Gameover detected early")
		}
	}

	ctl.Slam()
	ctl.NextTet(NewTet(TET_LINE))

	// After 6 line peices, a game over should activate
	if !ctl.isGameover {
		t.Error("No gameover detected")
	}
}

func TestBoardForceTick(t *testing.T) {
	board := &Board{}
	ctl := NewBoardController(board, NewTet(TET_LINE))

	// If we send force down movements, (simulates doing nothing), the
	// game should eventually end
	for counter := 0; counter < 200; counter++ {
		ctl.Tick(MOVE_FORCE_DOWN, NewTet(TET_LINE))
	}

	if !ctl.isGameover {
		t.Error("Game should have ended when doing nothing")
	}
}

func TestGameClearLines(t *testing.T) {
	game := NewGame(0)

	if game.Level() != 1 {
		t.Error("Default level is incorrect")
	}

	game.ClearLines(5)

	if game.lines != 5 {
		t.Error("Lines aren't getting properly added")
	}

	game.ClearLines(100)

	if game.Level() != MAX_LEVEL {
		t.Error("Max level wasn't reached")
	}
}

func TestNextTet(t *testing.T) {
	game := NewGame(0)

	tet := game.nextTet

	for i := 0; i < 10; i++ {
		game.NextTet()

		if tet == game.nextTet {
			t.Error("Next tetromino did not change")
		}

		tet = game.nextTet
	}
}

// Simulates a game where the user does absolutely nothing at all
func TestGameSimNothing(t *testing.T) {
	game := NewGame(0)

	for !game.controller.isGameover {
		game.Tick(MOVE_FORCE_DOWN)
	}
}

// Simulates a game with totally random movements
func TestGameRandom(t *testing.T) {
	game := NewGame(0)

	var move Movement
	for !game.controller.isGameover {
		for i := 0; i < 10; i++ {
			move = Movement(rand.Intn(int(MOVE_FORCE_DOWN)))
			game.Tick(move)
		}

		tet := game.nextTet
		for tet == game.nextTet {
			game.Tick(MOVE_FORCE_DOWN)
		}
	}

	if game.score < 0 {
		t.Error("Score is negative somehow")
	}
}

func TestGameListen(t *testing.T) {
	game := NewGame(0)

	// A channel that only sends slams, so we finish our game quickly
	slamChan := make(chan Movement)
	go func() {
		for {
			slamChan <- MOVE_SLAM
		}
	}()

	snapChan := make(chan GameSnapshot)
	go func() {
		for snap := range snapChan {
			t.Log("Game board:")
			t.Logf("%v", &snap.Board)
		}
	}()

	game.Play(slamChan, snapChan, false)

	if !game.controller.isGameover {
		t.Error("Expected gameover, but game is still active")
	}
}
