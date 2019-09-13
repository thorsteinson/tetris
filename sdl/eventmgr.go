package sdl

import (
	gosdl "github.com/veandco/go-sdl2/sdl"
	"tetris/lib"
)

type EventMgr struct {
	C       chan lib.Movement
}

var defaultInputMap map[gosdl.Keycode]lib.Movement
var debugInputMap map[gosdl.Keycode]lib.Movement

func init() {
	defaultInputMap = map[gosdl.Keycode]lib.Movement{
		gosdl.K_DOWN:  lib.MOVE_DOWN,
		gosdl.K_UP:    lib.MOVE_ROTATE_LEFT,
		gosdl.K_LEFT:  lib.MOVE_LEFT,
		gosdl.K_RIGHT: lib.MOVE_RIGHT,
		gosdl.K_a:     lib.MOVE_ROTATE_LEFT,
		gosdl.K_d:     lib.MOVE_ROTATE_RIGHT,
		gosdl.K_SPACE: lib.MOVE_SLAM,
	}

	debugInputMap = map[gosdl.Keycode]lib.Movement{
		gosdl.K_UP:    lib.MOVE_UP,
		gosdl.K_DOWN:  lib.MOVE_DOWN,
		gosdl.K_LEFT:  lib.MOVE_LEFT,
		gosdl.K_RIGHT: lib.MOVE_RIGHT,
		gosdl.K_a:     lib.MOVE_ROTATE_LEFT,
		gosdl.K_d:     lib.MOVE_ROTATE_RIGHT,
		gosdl.K_SPACE: lib.MOVE_SLAM,
		gosdl.K_s:     lib.MOVE_FORCE_DOWN,
	}
}

func NewEventMgr(inC chan gosdl.Event, debug bool) *EventMgr {
	mapping := defaultInputMap
	if debug {
		mapping = debugInputMap
	}

	outC := make(chan lib.Movement)

	// Process all events
	go func() {
		var code gosdl.Keycode
		var kevt *gosdl.KeyboardEvent
		for evt := range inC {
			if evt.GetType() == gosdl.KEYDOWN {
				// Must be a keyboard event, cast it, and then get
				// it's keycode
				kevt = evt.(*gosdl.KeyboardEvent)
				code = kevt.Keysym.Sym
				outMove, ok := mapping[code]
				if ok {
					outC <- outMove
				}

			}
			// Ignore any events that don't fit the mapping
		}
	}()

	return &EventMgr{outC}
}
