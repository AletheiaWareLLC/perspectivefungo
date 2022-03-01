package main

import (
	"aletheiaware.com/perspectivefungo"
	"encoding/json"
	"flag"
	"log"
	"os"
)

func main() {
	flag.Parse()

	var p perspectivefungo.Puzzle

	args := flag.Args()
	reader := os.Stdin
	if len(args) > 0 {
		log.Println("Reading:", args[0])
		file, err := os.Open(args[0])
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		reader = file
	}
	if err := json.NewDecoder(reader).Decode(&p); err != nil {
		log.Fatal(err)
	}

	rotations, penalties := perspectivefungo.Score(&p)
	log.Println("Rotations:", rotations)
	log.Println("Penalties:", penalties)
}
