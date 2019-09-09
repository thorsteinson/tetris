package lib

// A tetromino that also has positional information so it can move
// around a board
type ActiveTetromino struct {
	*Tetromino
	Position
}

type Direction int
type Position struct {
	x, y int
}

const (
	UP Direction = iota
	DOWN
	LEFT
	RIGHT
)

const STARTING_X = 4
const STARTING_Y = 23

const GAMEOVER_LINE = 20

func NewActiveTet(t *Tetromino) ActiveTetromino {
	return ActiveTetromino{
		Tetromino: t,
		Position: Position{
			x: STARTING_X,
			y: STARTING_Y,
		},
	}
}

func (tet ActiveTetromino) Move(dir Direction) ActiveTetromino {
	switch dir {
	case UP:
		tet.y += 1
	case DOWN:
		tet.y -= 1
	case LEFT:
		tet.x -= 1
	case RIGHT:
		tet.x += 1
	default:
		panic("Illegal direction passed")
	}

	return tet
}

func (tet ActiveTetromino) GetPos() (int, int) {
	return tet.x, tet.y
}

// Returns a list of positions for use in a board, where the tiles appear
func (tet ActiveTetromino) ListPositions() []Position {
	x := tet.x
	y := tet.y

	ps := []Position{}
	mask := tet.Tetromino.GetMask()

	for dy := 0; dy < tet.size; dy++ {
		for dx := 0; dx < tet.size; dx++ {
			if mask[dy*tet.size+dx] {
				// dy is inverted because we're counting from the top,
				// but our grid is based in the reverse direction
				ps = append(ps, Position{x + dx, y - dy})
			}
		}
	}

	return ps
}

// Returns true if the tetromino can be moved in the given direction
// without intersecting any tiles in the board, and within the
// boundaries of the board
func (tet ActiveTetromino) CanMove(dir Direction, board *Board) bool {
	projectedPos := tet.Move(dir).ListPositions()

	for _, p := range projectedPos {
		if p.x < 0 ||
			p.x >= BOARD_WIDTH ||
			p.y < 0 ||
			p.y >= BOARD_HEIGHT ||
			!board.IsEmpty(p.x, p.y) {
			return false
		}
	}

	return true
}

// A BoardController is an entity that manages the state of a board
// and an active tetromino. It moves the tetromino around with respect
// to the board, and can glue the tetromino to the board as one would
// expect with tetris.
type BoardController struct {
	board       *Board
	tet         ActiveTetromino
	tetSource   <-chan *Tetromino
	lineCounter chan<- int
	gameover    chan<- struct{}
}

func NewBoardController(
	board *Board,
	source <-chan *Tetromino,
	lines chan<- int,
	gameover chan<- struct{},
) *BoardController {
	ctl := &BoardController{
		board:       board,
		tetSource:   source,
		lineCounter: lines,
		gameover:    gameover,
	}
	ctl.NextTet()

	return ctl
}

// Recieves next tetromino from channel, and changes the active
// tetromino to the next one received. This implicitly locks the tiles
// of whatever the previous tetromino was and alters the state of the
// board by applying tetris
func (ctl *BoardController) NextTet() {
	// This is the value of ctl.tet before it's been set. Need to do a
	// comparison with this so don't compare against a nil value when
	// checking the gameover line
	initTet := ActiveTetromino{}

	if lines := ctl.board.Tetris(); lines > 0 {
		ctl.lineCounter <- lines
	} else if ctl.tet != initTet {
		// Check for whether they are on the gameover line
		for _, p := range ctl.tet.ListPositions() {
			if p.y == GAMEOVER_LINE {
				var gameoverSignal struct{}
				ctl.gameover <- gameoverSignal
				return
			}
		}
	}

	// NewActiveTet will handle setting the default position
	ctl.tet = NewActiveTet(<-ctl.tetSource)

	// Check for a collision, and send the gameover signal if one is
	// detected.
	for _, p := range ctl.tet.ListPositions() {
		if !ctl.board.IsEmpty(p.x, p.y) {
			var gameoverSignal struct{}
			ctl.gameover <- gameoverSignal
			return
		}
	}

	// Set new tiles
	for _, p := range ctl.tet.ListPositions() {
		ctl.board.SetTile(ShapeToTC(ctl.tet.shape), p.x, p.y)
	}
}

// This is a helper function that let's us pass a function, and
// inbetween, we'll unset and then set the tiles. This ensures we
// don't leave any copies around. If no state changes, the tiles
// should be put where they were before
func (ctl *BoardController) updateTiles(callback func() ActiveTetromino) {
	// Unset the tiles
	for _, p := range ctl.tet.ListPositions() {
		ctl.board.SetTile(EMPTY, p.x, p.y)
	}

	ctl.tet = callback()

	// Set the new tiles
	for _, p := range ctl.tet.ListPositions() {
		ctl.board.SetTile(ShapeToTC(ctl.tet.shape), p.x, p.y)
	}
}

// Move will idempotently move the active tetris piece. If it can't be
// moved, then it won't be moved.
func (ctl *BoardController) Move(dir Direction) {
	if ctl.tet.CanMove(dir, ctl.board) {
		ctl.updateTiles(func() ActiveTetromino {
			return ctl.tet.Move(dir)
		})
	}
}

// Core rotation
func (ctl *BoardController) rotate(isLeft bool) {
	var rotationFunc, rotationInverse func()

	if isLeft {
		rotationFunc = ctl.tet.RotLeft
		rotationInverse = ctl.tet.RotRight
	} else {
		rotationFunc = ctl.tet.RotRight
		rotationInverse = ctl.tet.RotLeft
	}

	ctl.updateTiles(func() ActiveTetromino {
		// Apply the rotation
		rotationFunc()

		// Consider all edge cases
		// 1. Pushing to the left on the X-Axis
		// 2. Pushing up on the Y-Axis
		// 3. Pushing to the right on the X-Axis
		// 4. Pushing down on th Y- axis
		minX := 0
		minY := 0
		maxX := BOARD_WIDTH - 1
		maxY := BOARD_HEIGHT - 1
		for _, p := range ctl.tet.ListPositions() {
			// Find minimum and maximum x and y values
			if p.x > maxX {
				maxX = p.x
			} else if p.x < minX {
				minX = p.x
			}

			if p.y < minY {
				minY = p.y
			} else if p.y > maxY {
				maxY = p.y
			}
		}

		// Shift in needed directions so it's in bounds, then check for
		// any collisions
		var deltaX, deltaY int
		var yDir, xDir Direction

		// The direction and delta we apply depends on which threshold
		// was crossed.
		if minX < 0 {
			deltaX = 0 - minX
			xDir = RIGHT
		} else {
			deltaX = maxX - (BOARD_WIDTH - 1)
			xDir = LEFT
		}

		if minY < 0 {
			deltaY = 0 - minY
			yDir = UP
		} else {
			deltaY = maxY - (BOARD_HEIGHT - 1)
			yDir = DOWN
		}

		projectedTet := ctl.tet
		for i := 0; i < deltaX; i++ {
			projectedTet = projectedTet.Move(xDir)
		}
		for i := 0; i < deltaY; i++ {
			projectedTet = projectedTet.Move(yDir)
		}

		// Check for internal collisions
		var colliding bool
		for _, p := range projectedTet.ListPositions() {
			if !ctl.board.IsEmpty(p.x, p.y) {
				colliding = true
				break
			}
		}

		if colliding {
			// Undo the rotation, the operation is idempotent
			rotationInverse()
		}
		return projectedTet
	})
}

// Attempting to rotate left or right will rotate in place if
// possible, and possibly shift the tetromino along the x-axis to make
// it fit. So it will rotate AND possibly move the tetromino
func (ctl *BoardController) RotLeft() {
	ctl.rotate(true)
}

func (ctl *BoardController) RotRight() {
	ctl.rotate(false)
}

// Slam will have a tetromino fall all the way to the bottom of the
// board, or until it reaches something along it's path to the bottom.
func (ctl *BoardController) Slam() {
	ctl.updateTiles(func() ActiveTetromino {
		// Move as far down as possible
		projectedTet := ctl.tet
		for projectedTet.CanMove(DOWN, ctl.board) {
			projectedTet = projectedTet.Move(DOWN)
		}

		return projectedTet
	})
}

type Game struct {
	// Keeps track of number of lines that have been cleared
	lines int
	score int
	// A board that contains all of the static tiles in the game
	staticBoard Board
	// A board that only has tiles for the current tetromino in play
	playerBoard Board
	// A source of random tetrominos
	tetrominos chan *Tetromino
	// The next tetromino that will be put into player after the
	// current one is finishes
	nextTet *Tetromino
	// The current tetromino in play
	currentTet ActiveTetromino
}
