package main

import (
	"flag"
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

	evtMgr, _ := sdl.Init(*debug)

	game := lib.NewGame(time.Now().UnixNano())

	snaps := make(chan lib.GameSnapshot)
	go func() {
		for snap := range snaps {
			log.Print(snap.Position)
		}
	}()

	game.Listen(evtMgr.C, snaps, *debug)
}
