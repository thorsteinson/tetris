package sdl

import (
	gosdl "github.com/veandco/go-sdl2/sdl"

	"log"
	"os"
)

// This is needed for properly converting colors. For some reason we
// need to allocate this, so might as well do it at initialization
var pxFormat *gosdl.PixelFormat

// Initializes SDL and starts everything related to it. This must be
// called before other managers are initialized, since they rely on
// the functionality here. Also listens for the quit event and exits
// if we attempt to close the window
func Init(xres, yres int, debug bool) (*EventMgr, *DisplayMgr) {
	if err := gosdl.Init(gosdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	eventChan := make(chan gosdl.Event)

	go func() {
		for event := gosdl.WaitEvent(); true; event = gosdl.WaitEvent() {
			switch event.(type) {
			case *gosdl.QuitEvent:
				log.Print("Quit Event Received, exiting...")
				gosdl.Quit()
				os.Exit(0)

			default:
				eventChan <- event
			}
		}
	}()

	format, err := gosdl.AllocFormat(uint(gosdl.PIXELFORMAT_RGBA32))
	if err != nil {
		panic(err)
	}
	pxFormat = format

	if !debug {
		// Load the theme song
		audioMgr := &AudioMgr{}
		audioMgr.Init()
		err := audioMgr.Loop(SONG_PATH)
		if err != nil {
			panic(err)
		}
	}

	return NewEventMgr(eventChan, debug), NewDisplayMgr("Tetris", xres, yres)
}
