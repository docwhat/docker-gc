package main

import (
	"fmt"
	"time"

	recorder "github.com/docwhat/docker-gc/memrecorder"
)

type imageTagRecorder interface {
	SawImageTagAt(tag string, when time.Time)
}

func main() {
	fmt.Println("Press Control-C to exit...")
	tagger := recorder.NewMemRecorder()
	schedule(tagger)
	fmt.Printf("%v\n", tagger)
}

func schedule(tagger imageTagRecorder) {
	fmt.Printf("%v\n", tagger)
	tagger.SawImageTagAt("busybox", time.Now())
}
