package main

import (
	"flag"
	"log"

	"tetris/sdl"
)

func main() {
	debug := flag.Bool("debug", false, "Disable timer and allow free movement")
	flag.Parse()

	if *debug {
		log.Print("Debugging enabled")
	}

	evtMgr, _ := sdl.Init(*debug)

	for e := range evtMgr.C {
		log.Print("New Event:", e)
	}
}
