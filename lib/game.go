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
