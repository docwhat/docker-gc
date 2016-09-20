package memrecorder

import (
	"testing"
	"time"

	"docwhat.org/docker-gc/types"
)

func TestMatchesInterface(t *testing.T) {
	var _ types.Recorder = NewMemRecorder()
}

func TestSawImageTagAt(t *testing.T) {
	r := NewMemRecorder()

	tag := "flibbit"

	now := time.Now()
	then := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)

	r.SawImageTagAt(tag, then)
	r.Sweep(func(_ string, w time.Time) bool {
		if w != then {
			t.Error("The timestamp wasn't stored for the tag")
		}
		return false
	})

	r.SawImageTagAt(tag, now)
	r.Sweep(func(_ string, w time.Time) bool {
		if w != now {
			t.Errorf("The timestamp should be now, not %v", w)
		}
		return false
	})

	r.SawImageTagAt(tag, then)
	r.Sweep(func(_ string, w time.Time) bool {
		if w != now {
			t.Errorf("The timestamp should remain now, not %v", w)
		}
		return false
	})
}

func TestSawImageTag(t *testing.T) {
	r := NewMemRecorder()

	tag := "frobnitz"

	r.SawImageTag(tag)
	r.Sweep(func(found string, _ time.Time) bool {
		if found != tag {
			t.Error("I expected to find the tag")
		}
		return false
	})
}

func TestForget(t *testing.T) {
	r := NewMemRecorder()

	expectedTag := "frobnitz"
	unexpectedTag := "foobar"

	r.SawImageTag(expectedTag)
	r.SawImageTag(unexpectedTag)
	r.Forget(unexpectedTag)
	r.Sweep(func(found string, _ time.Time) bool {
		if found != expectedTag {
			t.Errorf("I expected to find the tag %v; found %v", expectedTag, found)
		}
		return false
	})
}

func TestSweep(t *testing.T) {
	r := NewMemRecorder()

	tag := "frobnitz"

	r.SawImageTag(tag)

	r.Sweep(func(_ string, _ time.Time) bool {
		return false
	})

	r.Sweep(func(found string, _ time.Time) bool {
		if found != tag {
			t.Error("I expected to find the tag")
		}
		return true
	})

	r.Sweep(func(_ string, _ time.Time) bool {
		t.Error("I shouldn't find any tags")
		return false
	})
}
