package types

import "time"

// A SweepHandler is a method to run on each image tag.
//
// Return true if you deleted the tag.
type SweepHandler func(tag string, when time.Time) bool

// A Recorder keeps track of when images are seen.
type Recorder interface {
	SawImageTag(string)
	SawImageTagAt(string, time.Time)
	Sweep(SweepHandler)
	Forget(string)
}
