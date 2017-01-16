package main

import (
	"sync"
	"time"
)

// MemRecorder stores all seen image tags in memory.
type MemRecorder struct {
	imageTags map[string]time.Time
	mutex     sync.Mutex
}

// NewMemRecorder initializes a new MemRecorder for use.
func NewMemRecorder() *MemRecorder {
	var r MemRecorder
	r.imageTags = make(map[string]time.Time)
	return &r
}

// SawImageTagAt records when a tag was last seen.
//
// Duplicate or older times will be ignored.
func (r *MemRecorder) SawImageTagAt(tag string, when time.Time) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if old, ok := r.imageTags[tag]; ok {
		if when.Before(old) {
			return // We don't need to adjust the value.
		}
	}
	r.imageTags[tag] = when
}

// SawImageTag records a tag being seen now.
func (r *MemRecorder) SawImageTag(tag string) {
	r.SawImageTagAt(tag, time.Now())
}

// Forget that you saw a tag.
func (r *MemRecorder) Forget(tag string) {
	r.mutex.Lock()
	delete(r.imageTags, tag)
	r.mutex.Unlock()
}

// Sweep runs a function on all tag and timestamp pairs.
func (r *MemRecorder) Sweep(sweeper SweepHandler) {
	copiedImageTags := make(map[string]time.Time)
	r.mutex.Lock()
	for tag, when := range r.imageTags {
		copiedImageTags[tag] = when
	}
	r.mutex.Unlock()

	for tag, when := range copiedImageTags {
		if sweeper(tag, when) {
			r.Forget(tag)
		}
	}
}
