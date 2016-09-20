package main

import (
	"strings"
	"testing"
)

func TestNormalizeImageName(t *testing.T) {
	expected := "alpine:3.1"
	got := NormalizeImageName(expected)
	if !strings.EqualFold(expected, got) {
		t.Errorf("NormalizeImageName(%v) should have equaled %v", expected, got)
	}
}
