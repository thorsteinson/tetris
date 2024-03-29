package sdl

import (
	gosdl "github.com/veandco/go-sdl2/sdl"

	"image/color"

	"tetris/lib"
)

// A palette must have exactly 7 colors
type Palette [7]color.RGBA

func LookupColor(tc lib.TileColor, p Palette) color.RGBA {
	if tc == lib.EMPTY {
		// Pure black
		return color.RGBA{0, 0, 0, 255}
	}

	// Since the empty color has no actual color, shift every number
	// down by one
	return p[int(tc)-1]
}

// Helper function for creating surfaces
func NewSurface(w, h int) *gosdl.Surface {
	s, err := gosdl.CreateRGBSurfaceWithFormat(
		0,                  // Flags, ignore
		int32(w), int32(h), // Width and height
		32,                               // Depth
		uint32(gosdl.PIXELFORMAT_RGBA32), // RGBA Format
	)
	if err != nil {
		panic(err)
	}

	return s
}

// Helper function for quickly making rects
func Rect(x, y, w, h int) gosdl.Rect {
	return gosdl.Rect{
		int32(x),
		int32(y),
		int32(w),
		int32(h),
	}
}

func ColorMap(color color.RGBA) uint32 {
	if pxFormat == nil {
		panic("pxFormat has not been initialized. Must initialize before calling this function")
	}
	r, g, b, a := color.RGBA()
	return gosdl.MapRGBA(pxFormat, uint8(r), uint8(g), uint8(b), uint8(a))
}

// Helper function for coloring rectangles in surfaces. Allows us to
// use ordinary colors from image/color without hasslej
func FillRect(surf *gosdl.Surface, rect gosdl.Rect, color color.RGBA) {
	err := surf.FillRect(&rect, ColorMap(color))
	if err != nil {
		panic(err)
	}
}

// Clears the specified surface, by setting everything to the
// transparent color
func ClearSurface(surf *gosdl.Surface) {
	fillColor := color.RGBA{255, 255, 255, 255}
	r := Rect(0, 0, int(surf.W), int(surf.H))
	FillRect(surf, r, fillColor)
}

type Component interface {
	GetSurface() *gosdl.Surface
	Draw()
	Update(lib.GameSnapshot)
}

type BoardComponent struct {
	board   lib.Board
	palette Palette
	surf    *gosdl.Surface
	w       int
	h       int
}

const W_MIN = 50
const H_MIN = 100

func NewBoardComponent(initBoard lib.Board, p Palette, w int, h int) *BoardComponent {
	if w < W_MIN || h < H_MIN {
		panic("Cannot create board component that small. Minimum supported size is 50x100")
	}

	return &BoardComponent{
		board:   initBoard,
		palette: p,
		surf:    NewSurface(w, h),
		w:       w,
		h:       h,
	}
}

func (bc *BoardComponent) GetSurface() *gosdl.Surface {
	return bc.surf
}

func (bc *BoardComponent) Draw() {
	// Update the surface with the contents of the board

	// Figure out what the size of each rect should be
	var rectSize int
	if bc.h/bc.w >= 2 {
		rectSize = bc.w / 10
	} else {
		rectSize = bc.h / 20
	}

	realW := rectSize * 10
	realH := rectSize * 20

	xOff := (bc.w - realW) / 2
	yOff := (bc.h - realH) / 2

	// Draw tiles (0,0) -> (0,19), by drawing rectangles
	var rect gosdl.Rect
	var tc lib.TileColor
	for y := 0; y < 20; y++ {
		for x := 0; x < 10; x++ {
			tc = bc.board.GetTile(x, y)
			rect = Rect(xOff+x*rectSize, yOff+(20-y-1)*rectSize, rectSize, rectSize)
			FillRect(bc.surf, rect, LookupColor(tc, bc.palette))
		}
	}
}

func (bc *BoardComponent) Update(snap lib.GameSnapshot) {
	b := snap.Board
	if b != bc.board {
		ClearSurface(bc.surf)
		bc.board = b
		bc.Draw()
	}
}

// Creates a grid that is meant to be directly overlayed on top of a
// board, so it's more apparent how the tetrominos are layed out. This
// is a static component, so the surface is returned directly
func MakeGrid(w, h int) *gosdl.Surface {
	if w < W_MIN || h < H_MIN {
		panic("Cannot create a grid of provided size. Minumum supported size is 50x100")
	}

	surf := NewSurface(w, h)

	var lineSize int
	if h/w >= 2 {
		lineSize = w / 10
	} else {
		lineSize = h / 20
	}

	realW := lineSize * 10
	realH := lineSize * 20

	xOff := (w - realW) / 2
	yOff := (h - realH) / 2

	LINE_COLOR := color.RGBA{200, 200, 200, 200}

	var line gosdl.Rect
	// Draw horizonal lines
	for y := 0; y < 20; y++ {
		line = Rect(xOff, yOff+y*lineSize, realW, 1)
		FillRect(surf, line, LINE_COLOR)
	}
	line = Rect(xOff, yOff+20*lineSize-1, realW, 1)
	FillRect(surf, line, LINE_COLOR)
	// Draw vertical lines
	for x := 0; x < 10; x++ {
		line = Rect(xOff+x*lineSize, yOff, 1, realH)
		FillRect(surf, line, LINE_COLOR)
	}
	line = Rect(xOff+10*lineSize-1, yOff, 1, realH)
	FillRect(surf, line, LINE_COLOR)

	return surf
}
