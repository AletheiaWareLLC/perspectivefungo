package main

import (
	"aletheiaware.com/perspectivefungo"
	"encoding/json"
	"flag"
	"log"
	"math"
	"os"
)

const MAX_TRIES = 1000000000

var (
	size      = flag.Uint("size", 5, "Puzzle size")
	rotations = flag.Uint("rotations", 1, "Minimum puzzle rotations")
	penalties = flag.Uint("penalties", math.MaxUint, "Maximum puzzle penalties")
	blocks    = flag.Uint("blocks", 4, "Number of blocks")
	portals   = flag.Uint("portals", 2, "Number of portals")
)

func main() {
	flag.Parse()

	var puzzle *perspectivefungo.Puzzle

	for i, best := 0, *rotations/2; i < MAX_TRIES && best < *rotations; i++ {
		p, err := perspectivefungo.Generate(*size, *blocks, *portals)
		if err != nil {
			log.Fatal(err)
		}
		rs, ps := perspectivefungo.Score(p)
		if best < rs && ps <= *penalties {
			log.Println("Iteration:", i)
			log.Println("Size:", *size)
			log.Println("Rotations:", rs, "/", *rotations)
			log.Println("Penalties:", ps)
			puzzle = p
			best = rs
		}
	}

	args := flag.Args()
	writer := os.Stdout
	if len(args) > 0 {
		log.Println("Writing:", args[0])
		file, err := os.Create(args[0])
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		writer = file
	}
	if err := json.NewEncoder(writer).Encode(puzzle); err != nil {
		log.Fatal(err)
	}
}
