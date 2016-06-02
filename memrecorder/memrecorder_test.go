package memrecorder

import (
	"testing"
	"time"

	"github.com/docwhat/docker-gc/types"
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
