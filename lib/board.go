package lib

const (
	BOARD_WIDTH  = 10
	BOARD_HEIGHT = 40
	BOARD_SIZE   = BOARD_WIDTH * BOARD_HEIGHT
)

type TileColor int

// A board is basically a grid of tiles that have an associated
// color. We say that 0 is EMPTY, or has no color at all. A nice side
// effect is that a new board basically consists entirely of the EMPTY
// tile color. The other colors can serve as an index into some
// palette to determine the actual color.
const (
	EMPTY TileColor = iota
	C1
	C2
	C3
	C4
	C5
	C6
	C7
)

type Board struct {
	tiles [BOARD_SIZE]TileColor
}

// GetTile returns the index of the provided tile. The bottom left
// point is considered 0,0
func (b *Board) GetTile(x, y int) TileColor {
	if x < 0 || x >= BOARD_WIDTH {
		panic("Tile outside of board width")
	}
	if  y < 0 || y >= BOARD_HEIGHT {
		panic("Tile outside of board height")
	}

	return b.tiles[coordToTileIdx(x, y)]
}

// Helper function which converts coordinates for us
func coordToTileIdx(x, y int) int {
	y = BOARD_HEIGHT - y - 1
	return y*BOARD_WIDTH + x
}

// A helper for testing that a given tile color is in the valid range
// of values we've set
func invalidTile(t TileColor) bool {
	return t > C7 || t < EMPTY
}

func (b *Board) SetTile(t TileColor, x, y int) {
	if invalidTile(t) {
		panic("Invalid tile value passed")
	}

	b.tiles[coordToTileIdx(x, y)] = t
}

// Clear completely resets the board with a new one that's empty
func (b *Board) Clear() {
	b.tiles = [BOARD_SIZE]TileColor{}
}

// EraseLine clears the provided line and sets it back to empty
func (b *Board) EraseLine(y int) {
	for x := 0; x< BOARD_WIDTH; x++ {
		b.SetTile(EMPTY, x, y)
	}
}

func (b *Board) FullLines() []int {
	lines := []int{}

	for y := 0; y < BOARD_HEIGHT; y++ {
		var isEmpty bool
		for x := 0; x < BOARD_WIDTH; x++ {
			if b.GetTile(x, y) == EMPTY {
				isEmpty = true
			}
		}

		if !isEmpty {
			lines = append(lines, y)
		}
	}

	return lines
}

func (b *Board) IsEmpty(x, y int) bool {
	return b.GetTile(x, y) == EMPTY
}
