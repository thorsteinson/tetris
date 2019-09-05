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
func (b *Board) GetTile(i, j int) TileColor {
	i = BOARD_HEIGHT - i - 1 // Reverse notion of up and down
	return b.tiles[i*BOARD_WIDTH+j]
}
