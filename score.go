package perspectivefungo

const (
	BAD  = -1
	GOOD = 1
)

var (
	left       = []int{-1, 0, 0}
	right      = []int{1, 0, 0}
	down       = []int{0, -1, 0}
	up         = []int{0, 1, 0}
	backward   = []int{0, 0, -1}
	foreward   = []int{0, 0, 1}
	directions = [][]int{
		left,
		right,
		down,
		up,
		backward,
		foreward,
	}
)

func Score(puzzle *Puzzle) (uint, uint) {
	blocks := make(map[string]bool)
	for i := 0; i < len(puzzle.Blocks); i += 3 {
		blocks[Key(puzzle.Blocks[i], puzzle.Blocks[i+1], puzzle.Blocks[i+2])] = true
	}
	portals := make(map[string][]int)
	var pair *[]int
	for i := 0; i < len(puzzle.Portals); i += 3 {
		x := puzzle.Portals[i]
		y := puzzle.Portals[i+1]
		z := puzzle.Portals[i+2]
		if pair == nil {
			pair = &[]int{x, y, z}
		} else {
			p := *pair
			portals[Key(x, y, z)] = p
			portals[Key(p[0], p[1], p[2])] = []int{x, y, z}
			pair = nil
		}
	}
	tested := make(map[string]int)
	visited := make(map[string]bool)
	rotations, direction := ScoreDirections(puzzle.Size, puzzle.Player, puzzle.Goal, blocks, portals, tested, visited, false)
	penalty := uint(0)
	// Check all blocks were visited
	for i := 0; i < len(puzzle.Blocks); i += 3 {
		if !visited[Key(puzzle.Blocks[i], puzzle.Blocks[i+1], puzzle.Blocks[i+2])] {
			// log.Println("Unvisited Blocks: " + b.String())
			penalty++
		}
	}
	// Check all portals were visited
	for i := 0; i < len(puzzle.Portals); i += 3 {
		if !visited[Key(puzzle.Portals[i], puzzle.Portals[i+1], puzzle.Portals[i+2])] {
			// log.Println("Unvisited Portal: " + p.String())
			penalty++
			penalty++// Double penalty to encourage all portals to be visited
		}
	}
	if rotations < 0 {
		return 0, penalty
	}
	if direction[0] != down[0] || direction[1] != down[1] || direction[2] != down[2] {
		// Add initial rotation
		rotations += 1
	}
	return uint(rotations), penalty
}

func ScoreDirections(size uint, player, goal []int, blocks map[string]bool, portals map[string][]int, tested map[string]int, visited map[string]bool, portaled bool) (int, []int) {
	min := BAD
	dir := down
	posKey := Key(player[0], player[1], player[2])
	for _, d := range directions {
		key := posKey + Key(d[0], d[1], d[2])
		rotations, ok := tested[key]
		if !ok {
			tested[key] = BAD // Set now, update later to avoid loops
			rotations = ScoreDirection(size, []int{player[0], player[1], player[2]}, goal, blocks, portals, d, tested, visited, portaled)
			tested[key] = rotations
		}
		if rotations >= 0 && (rotations < min || min == BAD) {
			min = rotations
			dir = d
		}
	}
	return min, dir
}

func ScoreDirection(size uint, player, goal []int, blocks map[string]bool, portals map[string][]int, direction []int, tested map[string]int, visited map[string]bool, portaled bool) int {
	// log.Println("Scoring Direction:", direction)
	rotations := 0
	// Tracks portal usage to prevent infinite portal loops
	usage := make(map[string]int)
	for step := 0; ; step++ {
		// log.Println("Player:", player)
		if Abs(player[0]) > size || Abs(player[1]) > size || Abs(player[2]) > size {
			// log.Println("Out of Bounds")
			return BAD
		}
		if goal[0] == player[0] && goal[1] == player[1] && goal[2] == player[2] {
			// log.Println("Goal")
			return rotations
		}
		if step > 0 && !portaled {
			key := Key(player[0], player[1], player[2])
			link, ok := portals[key]
			if ok {
				uses, ok := usage[key]
				if !ok {
					uses = 0
				}
				if uses < 100 {
					// log.Println("Portal")
					player[0] = link[0]
					player[1] = link[1]
					player[2] = link[2]
					portaled = true
					usage[key] = uses + 1
					visited[key] = true
					visited[Key(link[0], link[1], link[2])] = true
					continue
				} else {
					// log.Println("Infinite Portal Loop")
					return BAD
				}
			}
		}
		next := []int{
			player[0] + direction[0],
			player[1] + direction[1],
			player[2] + direction[2],
		}
		if key := Key(next[0], next[1], next[2]); blocks[key] {
			// log.Println("Blocks")
			visited[key] = true
			r, _ := ScoreDirections(size, player, goal, blocks, portals, tested, visited, portaled)
			if r >= 0 {
				return r + rotations + GOOD
			}
			return r
		}
		player[0] = next[0]
		player[1] = next[1]
		player[2] = next[2]
		portaled = false
	}
}

func Abs(a int) uint {
	if a < 0 {
		return uint(-a)
	}
	return uint(a)
}
