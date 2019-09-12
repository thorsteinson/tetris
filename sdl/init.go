package sdl

import (
	gosdl "github.com/veandco/go-sdl2/sdl"
)

// Initializes SDL and starts everything related to it. This must be
// called before other managers are initialized, since they rely on
// the functionality here. Also listens for the quit event and exits
// if we attempt to close the window
func Init(debug bool) (*EventMgr, *DisplayMgr) {
	const DEFAULT_XRES = 800
	const DEFAULT_YRES = 600
	if err := gosdl.Init(gosdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	eventChan := make(chan gosdl.Event)

	go func() {
		for event := gosdl.WaitEvent(); true; event = gosdl.WaitEvent() {
			switch event.(type) {
			case *gosdl.QuitEvent:
				gosdl.Quit()
			}

			eventChan <- event
		}
	}()

	return NewEventMgr(eventChan, debug), NewDisplayMgr("Tetris", DEFAULT_XRES, DEFAULT_YRES)
}
