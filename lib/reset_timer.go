package lib

import (
	"time"
)

// A timer that sends a signal after a specified duration, but has a
// reset channel that will reset the timer. As long as the signal
// get's applied, it will never fire
type ResetTimer struct {
	reset         chan struct{}
	internalTimer *time.Timer
	out           chan struct{}
}

func NewResetTimer(duration time.Duration) *ResetTimer {
	timer := time.NewTimer(duration)
	resetChan := make(chan struct{})
	outChan := make(chan struct{})

	// Manages logic internally for handling signals to the channels
	go func() {
		for {
			select {
			case <-resetChan:
				timer.Stop()
				timer.Reset(duration)
			case <-timer.C:
				timer.Reset(duration)
				var outSignal struct{}
				outChan <- outSignal
			}
		}
	}()

	return &ResetTimer{
		reset:         resetChan,
		internalTimer: timer,
		out:           outChan,
	}

}

func (t *ResetTimer) Reset() {
	var resetSig struct{}
	t.reset <- resetSig
}

func (t *ResetTimer) Stop() {
	t.internalTimer.Stop()
}
