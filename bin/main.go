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
	flag.Parse()

	if *debug {
		log.Print("Debugging enabled")
	}

	evtMgr, disMgr := sdl.Init(*debug)

	game := lib.NewGame(time.Now().UnixNano())
	initState := game.Snap()

	palette := [7]color.RGBA{
		color.RGBA{255, 0, 0, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{0, 0, 255, 255},
		color.RGBA{255, 0, 255, 255},
		color.RGBA{0, 255, 255, 255},
		color.RGBA{255, 255, 0, 255},
		color.RGBA{255, 255, 125, 255},
	}

	boardComp := sdl.NewBoardComponent(initState.Board, palette, 500, 500)

	disMgr.Add(boardComp)

	snaps := make(chan lib.GameSnapshot)

	go disMgr.Render(snaps)

	game.Play(evtMgr.C, snaps, *debug)
}
