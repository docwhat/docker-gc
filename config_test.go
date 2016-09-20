package main

import (
	"testing"
	"time"
)

func TestMaxAgeOfImagesFlag(t *testing.T) {
	var testTable = []struct {
		in  []string
		out time.Duration
	}{
		{[]string{}, 168 * time.Hour},
		{[]string{"--max-image-age=20m"}, 20 * time.Minute},
		{[]string{"--max-image-age", "7h"}, 7 * time.Hour},
		{[]string{"-m", "5h"}, 5 * time.Hour},
	}

	for _, tt := range testTable {
		config := newAppConfig(tt.in)
		got := config.MaxAgeOfImages
		if got != tt.out {
			t.Errorf("args %q config.MaxAgeOfImages => %q, want %q", tt.in, got, tt.out)
		}
	}
}

func TestSweeperTimeFlag(t *testing.T) {
	var testTable = []struct {
		in  []string
		out time.Duration
	}{
		{[]string{}, 15 * time.Minute},
		{[]string{"--sweeper-time=11m"}, 11 * time.Minute},
		{[]string{"--sweeper-time", "2h"}, 2 * time.Hour},
		{[]string{"-s", "30s"}, 30 * time.Second},
	}

	for _, tt := range testTable {
		config := newAppConfig(tt.in)
		got := config.SweeperTime
		if got != tt.out {
			t.Errorf("args %q config.SweeperTime => %q, want %q", tt.in, got, tt.out)
		}
	}
}

func TestDangleSafeDurationFlag(t *testing.T) {
	var testTable = []struct {
		in  []string
		out time.Duration
	}{
		{[]string{}, 30 * time.Minute},
		{[]string{"--dangle-safe-duration=11m"}, 11 * time.Minute},
		{[]string{"--dangle-safe-duration", "2h"}, 2 * time.Hour},
		{[]string{"-d", "30s"}, 30 * time.Second},
	}

	for _, tt := range testTable {
		config := newAppConfig(tt.in)
		got := config.DangleSafeDuration
		if got != tt.out {
			t.Errorf("args %q config.DangleSafeDuration => %q, want %q", tt.in, got, tt.out)
		}
	}
}
