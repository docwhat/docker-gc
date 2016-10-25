package main

import (
	"testing"
	"time"
)

type configFlagTestTableRow struct {
	attrName string
	in       []string
	isGood   func(appConfig) bool
}

var configFlagTestTable = []configFlagTestTableRow{
	// MaxAgeOfImages
	{"MaxAgeOfImages", []string{}, func(c appConfig) bool { return c.MaxAgeOfImages == 168*time.Hour }},
	{"MaxAgeOfImages", []string{"--max-image-age=20m"}, func(c appConfig) bool { return c.MaxAgeOfImages == 20*time.Minute }},
	{"MaxAgeOfImages", []string{"--max-image-age", "7h"}, func(c appConfig) bool { return c.MaxAgeOfImages == 7*time.Hour }},
	{"MaxAgeOfImages", []string{"-m", "5h"}, func(c appConfig) bool { return c.MaxAgeOfImages == 5*time.Hour }},

	// SweeperTime
	{"SweeperTime", []string{}, func(c appConfig) bool { return c.SweeperTime == 15*time.Minute }},
	{"SweeperTime", []string{"--sweeper-time=11m"}, func(c appConfig) bool { return c.SweeperTime == 11*time.Minute }},
	{"SweeperTime", []string{"--sweeper-time", "2h"}, func(c appConfig) bool { return c.SweeperTime == 2*time.Hour }},
	{"SweeperTime", []string{"-s", "30s"}, func(c appConfig) bool { return c.SweeperTime == 30*time.Second }},

	// DangleSafeDuration
	{"DangleSafeDuration", []string{}, func(c appConfig) bool { return c.DangleSafeDuration == 30*time.Minute }},
	{"DangleSafeDuration", []string{"--dangle-safe-duration=11m"}, func(c appConfig) bool { return c.DangleSafeDuration == 11*time.Minute }},
	{"DangleSafeDuration", []string{"--dangle-safe-duration", "2h"}, func(c appConfig) bool { return c.DangleSafeDuration == 2*time.Hour }},
	{"DangleSafeDuration", []string{"-d", "30s"}, func(c appConfig) bool { return c.DangleSafeDuration == 30*time.Second }},

	// Quiet
	{"Quiet", []string{}, func(c appConfig) bool { return !c.Quiet }},
	{"Quiet", []string{"-q"}, func(c appConfig) bool { return c.Quiet }},
	{"Quiet", []string{"--quiet"}, func(c appConfig) bool { return c.Quiet }},
}

func TestConfigFlagTestTable(t *testing.T) {
	for _, tt := range configFlagTestTable {
		config, err := newAppConfig(tt.in)
		if err != nil {
			t.Errorf("args %q; got error %q", tt.in, err)
		}
		if !tt.isGood(config) {
			t.Errorf("args %v config => %+v", tt.in, config)
		}
	}
}

// See issue #4
func TestBogusFlag(t *testing.T) {
	args := []string{"--bogus-flag"}
	_, err := newAppConfig(args)
	if err == nil {
		t.Errorf("args %q; expected error, but got nothing", args)
	}
}
