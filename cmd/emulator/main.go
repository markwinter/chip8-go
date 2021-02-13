package main

import (
	"flag"
	"github.com/markwinter/chip8-go/pkg/chipeight"
	"log"
)

var (
	file = flag.String("file", "", "Path to the Chip8 ROM to load")
)

func main() {
	flag.Parse()

	if *file == "" {
		log.Fatalf("Must provide path to Chip8 ROM")
	}

	c8 := chipeight.NewChipeight()

	loadErr := c8.Load(*file)
	if loadErr != nil {
		log.Fatalf("failed to load file: %v", loadErr)
	}

	c8.Run()
}
