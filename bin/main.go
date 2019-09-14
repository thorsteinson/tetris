package main

import (
	"flag"
	"image/color"
	"log"
	"time"

	"tetris/lib"
	"tetris/sdl"
)

func main() {
	debug := flag.Bool("debug", false, "Disable timer and allow free movement")
	level := flag.Int("level", 1, "Starting level (1-20)")
	flag.Parse()

	if *debug {
		log.Print("Debugging enabled")
	}

	evtMgr, disMgr := sdl.Init(*debug)

	game := lib.NewGame(time.Now().UnixNano(), *level)
	initState := game.Snap()

	palette := [7]color.RGBA{
		color.RGBA{224, 166, 20, 255},
		color.RGBA{52, 193, 21, 255},
		color.RGBA{139, 188, 176, 255},
		color.RGBA{39, 62, 165, 255},
		color.RGBA{0, 255, 255, 255},
		color.RGBA{185, 57, 214, 255},
		color.RGBA{214, 57, 60, 255},
	}

	boardComp := sdl.NewBoardComponent(initState.Board, palette, 800, 600)

	disMgr.Add(boardComp)
	disMgr.AddSurf(sdl.MakeGrid(800, 600))

	snaps := make(chan lib.GameSnapshot)

	go disMgr.Render(snaps)

	game.Play(evtMgr.C, snaps, *debug)
}
