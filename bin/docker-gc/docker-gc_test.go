package main

import (
	"strings"
	"testing"
)

func TestNormalizeRepoTag(t *testing.T) {
	expected := "alpine:3.1"
	got := normalizeRepoTag(expected)
	if !strings.EqualFold(expected, got) {
		t.Errorf("normalizeRepoTag(%v) should have equaled %v", expected, got)
	}
}
