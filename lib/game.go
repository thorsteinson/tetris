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
