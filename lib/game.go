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

const STARTING_X = 5
const STARTING_Y = 23

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
	board     *Board
	tet       ActiveTetromino
	tetSource <-chan *Tetromino
}

func NewBoardController(board *Board, source <-chan *Tetromino) *BoardController {
	ctl := &BoardController{
		board:     board,
		tetSource: source,
	}
	ctl.NextTet()

	return ctl
}

// Recieves next tetromino from channel, and changes the active
// tetromino to the next one received.
func (ctl *BoardController) NextTet() {
	// NewActiveTet will handle setting the default position
	ctl.tet = NewActiveTet(<-ctl.tetSource)

	// Set new tiles
	for _, p := range ctl.tet.ListPositions() {
		ctl.board.SetTile(ShapeToTC(ctl.tet.shape), p.x, p.y)
	}
}

// Move will idempotently move the active tetris piece. If it can't be
// moved, then it won't be moved.
func (ctl *BoardController) Move(dir Direction) {
	if ctl.tet.CanMove(dir, ctl.board) {
		// Erase tiles from board
		for _, p := range ctl.tet.ListPositions() {
			ctl.board.SetTile(EMPTY, p.x, p.y)
		}

		ctl.tet = ctl.tet.Move(dir)

		// Set the new tiles
		for _, p := range ctl.tet.ListPositions() {
			ctl.board.SetTile(TileColor(ctl.tet.shape), p.x, p.y)
		}
	}
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
