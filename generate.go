package perspectivefungo

import (
	"math/rand"
	"time"
)

func Generate(size, blocks, portals uint) (*Puzzle, error) {
	rand.Seed(time.Now().UnixNano())

	occupied := make(map[string]bool, 2+blocks+portals)

	p := &Puzzle{
		Size: size,
	}
	p.Player = GenerateLocation(occupied, size)
	p.Goal = GenerateLocation(occupied, size)
	for i := uint(0); i < blocks; i++ {
		p.Blocks = append(p.Blocks, GenerateLocation(occupied, size)...)
	}
	for i := uint(0); i < portals/2; i++ {
		p.Portals = append(p.Portals, GenerateLocation(occupied, size)...)
		p.Portals = append(p.Portals, GenerateLocation(occupied, size)...)
	}
	return p, nil
}

func GenerateLocation(occupied map[string]bool, size uint) []int {
	var (
		x, y, z int
		key     string
	)
	for {
		x = RandomLocation(size)
		y = RandomLocation(size)
		z = RandomLocation(size)
		key = Key(x, y, z)
		if !occupied[key] {
			occupied[key] = true
			return []int{x, y, z}
		}
	}
}

func RandomLocation(size uint) int {
	s := int(size)
	return rand.Intn(s) - s/2
}
