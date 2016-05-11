package memrecorder

import (
	"testing"
	"time"
)

func TestSawImageTagAt(t *testing.T) {
	r := NewMemRecorder()

	tag := "flibbit"

	now := time.Now()
	then := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)

	r.SawImageTagAt(tag, then)
	r.Sweep(func(_ string, w time.Time) {
		if w != then {
			t.Error("The timestamp wasn't stored for the tag")
		}
	})

	r.SawImageTagAt(tag, now)
	r.Sweep(func(_ string, w time.Time) {
		if w != now {
			t.Errorf("The timestamp should be now, not %v", w)
		}
	})

	r.SawImageTagAt(tag, then)
	r.Sweep(func(_ string, w time.Time) {
		if w != now {
			t.Errorf("The timestamp should remain now, not %v", w)
		}
	})
}
