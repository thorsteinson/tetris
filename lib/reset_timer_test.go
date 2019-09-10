package lib

import (
	"testing"
	"time"
)

func TestResetTimerNoReset(t *testing.T) {
	const DURATION = time.Microsecond * 500
	// Create a timer that resets

	rTimer := NewResetTimer(DURATION)

	t0 := time.Now()
	// Receive 3 values
	<-rTimer.out
	<-rTimer.out
	<-rTimer.out

	t1 := time.Now()

	// Calculate the duration. Note that this will fail at smaller
	// times because there just isn't enough accuracy.
	dur := t1.Sub(t0)

	minDur := DURATION * 3
	maxDur := DURATION * 4

	rTimer.Stop()

	if dur < minDur || dur > maxDur {
		t.Errorf("Timing error: Expected duration between %v and %v , Found %v", minDur, maxDur, dur)
	}
}

// Test that the reset time actually will reset when we send a signal
func TestResetTimerReset(t *testing.T) {
	const DURATION = time.Microsecond * 500

	rTimer := NewResetTimer(DURATION)

	// Receive 3 values

	// Create a ticker that sends values at a regular interval, that's
	// less than our duration
	ticker := time.NewTicker(DURATION / 2)

	t0 := time.Now()

	<-ticker.C
	rTimer.Reset()
	// 250
	<-ticker.C
	rTimer.Reset()
	// 250
	<-rTimer.out
	// 500

	// We would expect the timer to be between 1000 and 1500 microseconds

	t1 := time.Now()

	minDur := DURATION * 2
	maxDur := DURATION * 3

	dur := t1.Sub(t0)

	rTimer.Stop()

	if dur < minDur || dur > maxDur {
		t.Errorf("Duration outside minimum of %v and max of %v; found %v", minDur, maxDur, dur)
	}
}

func TestResetTimerChangeDuration(t *testing.T) {
	const SMALL_DUR = time.Microsecond * 10
	const LARGE_DUR = time.Microsecond * 500

	// Create a timer with a ridiculously small
	rTimer := NewResetTimer(SMALL_DUR)

	t0 := time.Now()

	rTimer.duration = LARGE_DUR
	<-rTimer.out
	<-rTimer.out

	t1 := time.Now()

	if dur := t1.Sub(t0); dur < LARGE_DUR {
		t.Errorf("Changing duration failed expected duration over %v, found %v", LARGE_DUR, dur)
	}
}
