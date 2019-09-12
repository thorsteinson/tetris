package sdl

import (
	gosdl "github.com/veandco/go-sdl2/sdl"
)

type DisplayMgr struct{
	window *gosdl.Window
	name string
}

func NewDisplayMgr(name string, xres, yres int) *DisplayMgr {
	mgr := &DisplayMgr{name: name}
	mgr.CreateWindow(xres, yres)

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
