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
	currentPos := tet.ListPositions()
	projectedPos := tet.Move(dir).ListPositions()

	// Test whether a position is in a position that the tetromino is
	// currently in. It's okay to occupy a future space that' we're
	// already inside of
	isInCurrentPos := func(p Position) bool {
		for _, c := range currentPos {
			if c == p {
				return true
			}
		}

		return false
	}

	for _, p := range projectedPos {
		if isInCurrentPos(p) {
			continue
		} else if p.x < 0 ||
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
	board      *Board
	tet        ActiveTetromino
	isGameover bool
}

func NewBoardController(board *Board, tet *Tetromino) *BoardController {
	ctl := &BoardController{board: board}
	ctl.NextTet(tet)

	return ctl
}

// Returns the number of lines cleared, if any and sets the next
// tetromino as the one passed.
func (ctl *BoardController) NextTet(next *Tetromino) int {
	// This is the value of ctl.tet before it's been set. Need to do a
	// comparison with this so don't compare against a nil value when
	// checking the gameover line
	initTet := ActiveTetromino{}
	lines := ctl.board.Tetris()

	if lines == 0 && ctl.tet != initTet {
		// Check for whether they are on the gameover line
		for _, p := range ctl.tet.ListPositions() {
			if p.y == GAMEOVER_LINE {
				ctl.isGameover = true
				return 0
			}
		}
	}

	// NewActiveTet will handle setting the default position
	ctl.tet = NewActiveTet(next)

	// Do a quick gameover check
	for _, p := range ctl.tet.ListPositions() {
		if !ctl.board.IsEmpty(p.x, p.y) {
			ctl.isGameover = true
			return 0
		}
	}

	// Set new tiles
	for _, p := range ctl.tet.ListPositions() {
		ctl.board.SetTile(ShapeToTC(ctl.tet.shape), p.x, p.y)
	}

	return lines
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

// Helper method that conveniently checks whether a tile can move
// down. Need to check this at the game level, so this method helps
// clear that logic up a bit.
func (ctl *BoardController) CanMoveDown() bool {
	return ctl.tet.CanMove(DOWN, ctl.board)
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
	MOVE_FORCE_DOWN
)

// Tick will apply some sort of move and atomically update the board
// with that given move. The board before and after tick will always
// be in a consistent sensible state.
//
// Returns number of lines cleared between any movement, and whether
// the tetromino was consumed
func (ctl *BoardController) Tick(move Movement, next *Tetromino) (int, bool) {
	var lines int
	var consumed bool

	if move <= MOVE_RIGHT {
		// Movement must be a direction
		ctl.Move(Direction(move))
	} else {
		switch move {
		case MOVE_ROTATE_LEFT:
			ctl.RotLeft()
		case MOVE_ROTATE_RIGHT:
			ctl.RotRight()
		case MOVE_SLAM:
			if ctl.tet.CanMove(DOWN, ctl.board) {
				ctl.Slam()
			} else {
				// If slam is double tapped, treat it as the user
				// locking the tile in place, send the next tet over
				lines = ctl.NextTet(next)
				consumed = true
			}
			ctl.Slam()
		case MOVE_FORCE_DOWN:
			// This doesn't come from user input, but from a timer. It
			// can potentially trigger next tet if it's at the bottom
			if ctl.tet.CanMove(DOWN, ctl.board) {
				ctl.Move(DOWN)
			} else {
				lines = ctl.NextTet(next)
				consumed = true
			}
		}
	}

	return lines, consumed
}

type Game struct {
	// Keeps track of number of lines that have been cleared
	lines          int
	score          int
	level          int
	ticks          int
	controller     *BoardController
	nextTet        *Tetromino
	linesToNextLvl int
	tetSource      chan *Tetromino
}

func TetFactory(seed int64) chan *Tetromino {
	tets := make(chan *Tetromino)

	go func() {
		for shape := range ShapeGenerator(seed) {
			tets <- NewTet(shape)
		}
	}()

	return tets
}

const LINES_PER_LVL = 10
const MAX_LEVEL = 20
const DURATION_DIFF = 50 * time.Millisecond

// Create a new game with a given random seed, and hook it to some
// sort of movement channel to get inputs
func NewGame(seed int64) *Game {
	var next *Tetromino

	tets := TetFactory(seed)
	firstTet := <-tets
	next = <-tets

	game := &Game{
		controller:     NewBoardController(&Board{}, firstTet),
		level:          1,
		linesToNextLvl: LINES_PER_LVL,
		nextTet:        next,
		tetSource:      tets,
	}

	return game
}

// Applies logic for clearing lines. This modifies internal state so
// that our level is updated and score is modified
func (game *Game) ClearLines(cleared int) {
	game.linesToNextLvl -= cleared
	game.score += cleared * game.level

	if game.linesToNextLvl < 0 {
		game.level++
		game.linesToNextLvl = LINES_PER_LVL
	}
}

// Adds values to score based on the level, the number of lines
// cleared, the time, etc. Should only be called once for gameover
func (game *Game) CalcEndBonuses() int {
	var score int

	score += game.ticks
	score += game.level * 10
	if game.level == MAX_LEVEL {
		score += 1000
	}
	return score
}

// Fetches the next tetromino from internal source
func (game *Game) NextTet() {
	game.nextTet = <-game.tetSource
}

// Calculates a score that's meant to be applied between ordinairy non
// tetris ticks. It should only take the level and time into account
func (game *Game) CalcTickScore() int {
	return game.ticks * game.level
}

func (game *Game) Tick(move Movement) {
	game.ticks++ // Keeps track of the number of turns

	// Apply move to the board, get the number of lines
	cleared, consumed := game.controller.Tick(move, game.nextTet)

	if consumed {
		game.NextTet()
	}

	if cleared > 0 {
		// Tetris must have occurred
		game.ClearLines(cleared)
	} else if !game.controller.isGameover {
		game.score += game.CalcTickScore()
	} else {
		// Game is over
		game.score += game.CalcEndBonuses()
	}
}

// A value that represents a point in time for a given game. This can
// be produced by a game, and sent to something else to draw it or
// something else. It's intentionally a single large value
type GameSnapshot struct {
	Score      int
	Level      int
	Ticks      int
	Board      Board
	CurrentTet Tetromino
	NextTet    Tetromino
	Position   Position
}

func (game *Game) Snap() GameSnapshot {
	return GameSnapshot{
		Score:      game.score,
		Level:      game.level,
		Ticks:      game.ticks,
		Board:      *game.controller.board,
		CurrentTet: *game.controller.tet.Tetromino,
		NextTet:    *game.nextTet,
		Position:   game.controller.tet.Position,
	}
}

// Listens for incoming movements on a channel and applies them until
// the game is over. If called with the debug flag then the timer is
// disabled and movement is simply free form
func (game *Game) Listen(moves <-chan Movement, snaps chan<- GameSnapshot, debug bool) {
	if debug {
		for !game.controller.isGameover {
			for move := range moves {
				game.Tick(move)
				snaps <- game.Snap()
			}
		}
	} else {
		timer := NewResetTimer(DEFAULT_DURATION)

		var move Movement
		for !game.controller.isGameover {
			// Update the timer duration, this will progressively
			// speed the game up to a minimum of 50ms between moves at
			// lvl 20
			timer.duration = DEFAULT_DURATION - DURATION_DIFF*time.Duration(game.level)

			select {
			case <-timer.out:
				game.Tick(MOVE_FORCE_DOWN)
				snaps <- game.Snap()
				continue
			default:
			}

			select {
			case <-timer.out:
				move = MOVE_FORCE_DOWN
			case move = <-moves:
				if game.controller.CanMoveDown() && move == MOVE_DOWN || move == MOVE_SLAM {
					timer.Reset()
				}
			}

			game.Tick(move)
			snaps <- game.Snap()

		}
	}
}
