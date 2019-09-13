package sdl

import (
	"github.com/veandco/go-sdl2/mix"
)

type AudioMgr struct {}

const SONG_PATH = "../assets/theme.ogg"

const CHUNK_SIZE = 4096

// Starts the mixer for SDL. We only need to deal with ogg files
func (mgr *AudioMgr) Init() {
	var err error
	err = mix.Init(mix.INIT_OGG)
	err = mix.OpenAudio(mix.DEFAULT_FREQUENCY, mix.DEFAULT_FORMAT, mix.DEFAULT_CHANNELS, CHUNK_SIZE)
	if err != nil {
		panic(err)
	}

	// Set the volume way down, so our ears aren't blasted
	mix.Volume(-1, 5)
}

func (mgr *AudioMgr) Loop(filepath string) error {
	music, err := mix.LoadMUS(filepath)
	if err != nil {
		return err
	}

	// Loops the song forever
	music.Play(-1)

	return nil
}
