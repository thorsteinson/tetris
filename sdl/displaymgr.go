package sdl

import (
	gosdl "github.com/veandco/go-sdl2/sdl"

	"tetris/lib"
)

type DisplayMgr struct {
	window     *gosdl.Window
	winSurf    *gosdl.Surface
	name       string
	components []Component
	surfaces   []*gosdl.Surface
}

func NewDisplayMgr(name string, xres, yres int) *DisplayMgr {
	mgr := &DisplayMgr{name: name,}
	mgr.CreateWindow(xres, yres)
	s, err := mgr.window.GetSurface()
	if err != nil {
		panic(err)
	}
	mgr.winSurf = s

	return mgr
}

func (mgr *DisplayMgr) CreateWindow(xres, yres int) {
	window, err := gosdl.CreateWindow(
		mgr.name,
		gosdl.WINDOWPOS_UNDEFINED,
		gosdl.WINDOWPOS_UNDEFINED,
		int32(xres), int32(yres), gosdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	mgr.window = window
}

func (mgr *DisplayMgr) Add(comp Component) {
	mgr.surfaces = append(mgr.surfaces, comp.GetSurface())
	mgr.components = append(mgr.components, comp)
}

// Renders visuals to the screen
func (mgr *DisplayMgr) Render(snapshots chan lib.GameSnapshot) {
	for snap := range snapshots {
		for _, comp := range mgr.components {
			comp.Update(snap)
		}

		for  _,s := range mgr.surfaces {
			// blit each surface in order onto the screen, and then
			// update it
			// r := Rect(0, 0, int(mgr.winSurf.W), int(mgr.winSurf.H))
			s.Blit(nil, mgr.winSurf, nil)
		}

		mgr.window.UpdateSurface()
	}
}
