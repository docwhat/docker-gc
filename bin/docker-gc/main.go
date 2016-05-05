package main

import (
	"fmt"

	recorder "github.com/docwhat/docker-gc/recorder"
)

func main() {
	fmt.Println("Press Control-C to exit...")
	recorder.Start()
}
