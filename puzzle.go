package perspectivefungo

import (
	"fmt"
)

type Puzzle struct {
	Size    uint  `json:"size"`
	Player  []int `json:"player"`
	Goal    []int `json:"goal"`
	Blocks  []int `json:"blocks"`
	Portals []int `json:"portals"`
}

func Key(x, y, z int) string {
	return fmt.Sprintf("%d,%d,%d", x, y, z)
}
