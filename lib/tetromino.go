package lib

type Shape int

const (
	TET_SQUARE Shape = iota
	TET_S
	TET_Z
	TET_L
	TET_T
	TET_J
	TET_LINE
)

type TetGrid struct {
	grid []bool
	size int
}

var squareTet = TetGrid{
	grid: {
		true, true,
		true, true,
	},
	size: 2,
}

var sTet = TetGrid{
	grid: {
		false, false, false,
		false, true, true,
		true, true, false,
	},
	size: 3,
}

var zTet = TetGrid{
	grid: {
		false, false, false,
		true, true, false,
		false, true, true,
	},
	size: 3,
}

var lTet = TetGrid{
	grid: {
		false, true, false,
		false, true, false,
		false, true, true,
	},
	size: 3,
}

var jTet = TetGrid{
	grid: {
		false, false, false,
		true, true, true,
		false, true, false,
	},
	size: 3,
}

var jTet = TetGrid{
	grid: {
		false, true, false,
		false, true, false,
		true, true, false,
	},
	size: 3,
}

var lineTet = TetGrid{
	grid: {
		false, true, false, false,
		false, true, false, false,
		false, true, false, false,
		false, true, false, false,
	},
	size: 4,
}

// A tetromino is a given shape, that has some sort of mask which
// determines it's direction. The mask is a pointer to a slice,
// because we can do the rotations in advance, and just change to the
// right mask when we rotate it
type Tetromino struct {
	mask  *[]bool
	size  int
	shape Shape
}

// Pivots a square grid, returing the values of that grid but rotated.
func pivot(grid []bool, size int) []bool {
	switch size {
	case 2:
		return []bool{
			grid[1], grid[3],
			grid[0], grid[2],
		}
	case 3:
		return []bool{
			grid[2], grid[5], grid[8],
			grid[1], grid[4], grid[7],
			grid[0], grid[3], grid[6],
		}
	case 4:
		return []bool{
			grid[3], grid[7], grid[11], grid[15],
			grid[2], grid[6], grid[10], grid[14],
			grid[1], grid[5], grid[9], grid[13],
			grid[0], grid[4], grid[8], grid[12],
		}
	default:
		panic("Cannot rotate grid with provided size")
	}
}
