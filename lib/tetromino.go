package lib

import (
	"math/rand"
)

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

// Handy array we can use to refer to all shapes and iterate over
var shapes = []Shape{TET_SQUARE, TET_S, TET_Z, TET_L, TET_T, TET_J, TET_LINE}

// Reference to the greatest shape. Useful for iteration and for randomness
const maxShape = TET_LINE

type TetGrid struct {
	grid []bool
	size int
}

var squareGrid = TetGrid{
	grid: []bool{
		true, true,
		true, true,
	},
	size: 2,
}

var sGrid = TetGrid{
	grid: []bool{
		false, false, false,
		false, true, true,
		true, true, false,
	},
	size: 3,
}

var zGrid = TetGrid{
	grid: []bool{
		false, false, false,
		true, true, false,
		false, true, true,
	},
	size: 3,
}

var lGrid = TetGrid{
	grid: []bool{
		false, true, false,
		false, true, false,
		false, true, true,
	},
	size: 3,
}

var tGrid = TetGrid{
	grid: []bool{
		false, false, false,
		true, true, true,
		false, true, false,
	},
	size: 3,
}

var jGrid = TetGrid{
	grid: []bool{
		false, true, false,
		false, true, false,
		true, true, false,
	},
	size: 3,
}

var lineGrid = TetGrid{
	grid: []bool{
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
	mask        *[]bool
	size        int
	shape       Shape
	rotationIdx int
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

var rotations [][]*[]bool

func init() {
	// Set the values for the grids, so rotations are just a matter of
	// modifying an index for a lookup instead of actually doing a
	// rotation
	grids := []TetGrid{squareGrid, sGrid, zGrid, lGrid, tGrid, jGrid, lineGrid}
	for _, grid := range grids {
		rotationSet := []*[]bool{}

		var rotation []bool = grid.grid
		rotationSet = append(rotationSet, &rotation)

		for i := 0; i < 3; i++ {
			rotation := pivot(rotation, grid.size)
			rotationSet = append(rotationSet, &rotation)
		}

		rotations = append(rotations, rotationSet)
	}
}

func NewTet(s Shape) *Tetromino {
	lookupSize := func() int {
		switch s {
		case TET_SQUARE:
			return 2
		case TET_J:
			return 3
		case TET_L:
			return 3
		case TET_T:
			return 3
		case TET_S:
			return 3
		case TET_Z:
			return 3
		case TET_LINE:
			return 4
		default:
			panic("Invalid shape passed in")
		}
	}

	return &Tetromino{
		mask:  rotations[s][0],
		size:  lookupSize(),
		shape: s,
	}
}

func (tet *Tetromino) RotLeft() {
	tet.rotationIdx++
	// Reset to start if needed
	tet.rotationIdx = tet.rotationIdx % 4
	tet.mask = rotations[tet.shape][tet.rotationIdx]
}

func (tet *Tetromino) RotRight() {
	tet.rotationIdx--
	// Reset to end if needed
	if tet.rotationIdx < 0 {
		tet.rotationIdx = 3
	}
	tet.mask = rotations[tet.shape][tet.rotationIdx]
}

// Returns a copy of the mask that the tetromino is pointing to
// internally. This ensures that we never modify our rotations at any
// point and keep them safe.
func (tet *Tetromino) GetMask() []bool {
	mask := make([]bool, len(*tet.mask))
	copy(mask, *tet.mask)
	return mask
}

// Returns the mask that results from a left rotation
func (tet *Tetromino) GetLeftRotationMask() []bool {
	mask := make([]bool, len(*tet.mask))
	copy(mask, *rotations[tet.shape][tet.rotationIdx+1%4])
	return mask
}

// Returns the mask that results from a right rotation
func (tet *Tetromino) GetRightRotationMask() []bool {
	mask := make([]bool, len(*tet.mask))
	i := tet.rotationIdx - 1
	if i < 0 {
		i = 3
	}
	copy(mask, *rotations[tet.shape][i])
	return mask
}

// Creates a read only channel that sends random shapes. We
// paramaterize this with a seed
func ShapeGenerator(seed int64) <-chan Shape {
	r := rand.New(rand.NewSource(seed))

	shapes := make(chan Shape, 20)

	go func() {
		for {
			shapes <- Shape(r.Intn(int(maxShape)))
		}
	}()

	return shapes
}
