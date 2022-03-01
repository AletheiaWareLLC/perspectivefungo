package perspectivefungo

import (
	"time"
)

type Solution struct {
	Start, End time.Time
	Progress   []*Snapshot
}

type Snapshot struct {
	Time     time.Time
	Position [3]float32
}
