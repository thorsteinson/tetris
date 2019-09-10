package lib

import (
	"time"
)

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

const DEFAULT_DURATION = time.Second

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
	isGameover  bool
	timer       *ResetTimer
}

func NewBoardController(
	board *Board,
	source <-chan *Tetromino,
	lines chan<- int,
	gameover chan<- struct{},
	dur time.Duration,
) *BoardController {

	var timerDur time.Duration
	// Compare the zero value
	if dur != timerDur {
		timerDur = dur
	} else {
		timerDur = DEFAULT_DURATION
	}

	ctl := &BoardController{
		board:       board,
		tetSource:   source,
		lineCounter: lines,
		gameover:    gameover,
		timer:       NewResetTimer(timerDur),
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
				close(ctl.gameover)

				ctl.isGameover = true
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
			close(ctl.gameover)

			ctl.isGameover = true
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

type Movement int

const (
	MOVE_UP Movement = iota
	MOVE_DOWN
	MOVE_LEFT
	MOVE_RIGHT
	MOVE_SLAM
	MOVE_ROTATE_LEFT
	MOVE_ROTATE_RIGHT
)

// The Listen method connects to a movement channel which provides an
// input movement for the tetromino and moves it around. This could be
// as simple as a list of dedicated movements, or it could be tied to
// a real world input source to get interactive movement.
//
// There's a 2nd channel provided which is meant for locking the peice
// in place and triggering the next tetromino. In a normal game this
// can be constructed from a ticker that forces the tetromino to move
// down at a regular rate.
//
// This function will block until the game finishes. It will return
// immediately after the gameover signal has been fired
func (ctl *BoardController) Listen(moves <-chan Movement) {
	var dir Direction

	// Repeat this loop until the game finishes
	for !ctl.isGameover {

		select {
		case move := <-moves:
			if move <= MOVE_RIGHT {
				dir = Direction(move)
				if dir == DOWN && ctl.tet.CanMove(DOWN, ctl.board) {
					// A downwards movement should reset our
					// timer. But only if we're not already blocked.
					ctl.timer.Reset()
				}
				ctl.Move(dir)
			} else {
				switch move {
				case MOVE_SLAM:
					ctl.Slam()

					ctl.timer.Reset()

				case MOVE_ROTATE_LEFT:
					ctl.RotLeft()
				case MOVE_ROTATE_RIGHT:
					ctl.RotRight()
				}
			}

		case <-ctl.timer.out:
			if !ctl.tet.CanMove(DOWN, ctl.board) {
				ctl.NextTet()
			} else {
				ctl.Move(DOWN)
			}
		}
	}
}

type Game struct {
	// Keeps track of number of lines that have been cleared
	lines          int
	score          int
	controller     *BoardController
	level          int
	nextTet        *Tetromino
	moves          <-chan Movement
	linesToNextLvl int
}

// Create a new game with a given random seed, and hook it to some
// sort of movement channel to get inputs
func NewGame(seed int64, moves <-chan Movement) *Game {

	const LINES_PER_LVL = 10
	const MAX_LEVEL = 20
	const DURATION_DIFF = 50 * time.Millisecond

	// 	board *Board,
	// 	source <-chan *Tetromino,
	// 	lines chan<- int,
	// 	gameover chan<- struct{},
	// 	dur time.Duration,

	var next *Tetromino

	// Create infinite stream of tetrominos
	tets := make(chan *Tetromino)
	go func() {
		for s := range ShapeGenerator(seed) {
			nextTet := NewTet(s)
			// Update the pointer, so it points to the value that will
			// be next be consumed when the board calls NextTet(),
			// this gives us our preview
			next = nextTet
			// Push the value to the channel
			tets <- nextTet
		}
	}()

	lineC := make(chan int)

	gameoverC := make(chan struct{})

	game := &Game{
		controller: NewBoardController(
			&Board{},
			tets,
			lineC,
			gameoverC,
			0,
		),
		score:          0,
		level:          1,
		linesToNextLvl: LINES_PER_LVL,
		nextTet:        next,
	}

	// Increase line counter as it updates, and modify level if we
	// reach a threshold
	go func() {
		for cleared := range lineC {
			game.score += 100 * cleared

			game.lines += cleared
			game.linesToNextLvl -= cleared
			if game.linesToNextLvl < 0 {
				// increment the lvl
				if game.level < MAX_LEVEL {
					game.level++

					// Modify the speed of the game, decrease duration
					// between turns thus speeding it up.
					game.controller.timer.duration -= DURATION_DIFF
				}

				// Reset level counter
				game.linesToNextLvl = LINES_PER_LVL

				// Apply score bonus
				game.score += game.level * 1000

			}
		}
	}()

	return game
}

// Runs a game to completion
func (game *Game) Run() {
	game.controller.Listen(game.moves)
}
