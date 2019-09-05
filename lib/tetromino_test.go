package lib

import (
	"reflect"
	"testing"
	"time"
)

func TestTetRotation(t *testing.T) {
	// One property of rotation a tetromino is that after rotating it
	// 4 times, in either direction, it should be the same as when it
	// started.

	for _, s := range shapes {
		tet := NewTet(s)

		initMask := tet.GetMask()

		for i := 0; i < 4; i++ {
			tet.RotLeft()
		}

		if !reflect.DeepEqual(initMask, tet.GetMask()) {
			t.Errorf("Shape %v is not same after consecutive left rotations", s)
		}

		for i := 0; i < 4; i++ {
			tet.RotRight()
		}

		if !reflect.DeepEqual(initMask, tet.GetMask()) {
			t.Errorf("Shape %v is not same after consecutive right rotations", s)
		}
	}
}

func TestShapeGeneration(t *testing.T) {
	// Receives a bunch of shapes, and looks at the distribution of
	// the shapes. While random, it should be relatively equal. If
	// it's too wonky, throw an error

	shapeGen := ShapeGenerator(time.Now().UnixNano())
	counts := make(map[Shape]int)
	const SAMPLES = 1000

	var s Shape
	for i := 0; i < SAMPLES; i++ {
		s = <-shapeGen
		counts[s]++
	}

	const eqDist = 1 / float64(maxShape+1)
	// Allow for 20 percent variation
	const epsilon = eqDist * 0.35
	for _, c := range counts {
		distribution := float64(c) / float64(SAMPLES)
		if distribution > eqDist+epsilon || distribution < eqDist-epsilon {
			t.Error("Distribution of values is too far off")
			t.Errorf("Equal distribution: %v", eqDist)
			t.Errorf("Experiment distribution: %v", distribution)
		}
	}
}
