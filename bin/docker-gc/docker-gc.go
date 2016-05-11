package main

import (
	"fmt"

	recorder "github.com/docwhat/docker-gc/mem_recorder"
)

func main() {
	fmt.Println("Press Control-C to exit...")
	recorder.Start()
}
