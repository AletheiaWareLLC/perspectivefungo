package perspectivefungo

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
	"time"
)

const (
	ACCELERATION = 98.1
	INCREMENT    = 0.02
)

type Animation interface {
	Tick() bool // Return true if animation has completed
}

var (
	maxAngle = math.Pi / 20.
	xAxis    = mgl32.Vec4{1, 0, 0, 1}
	yAxis    = mgl32.Vec4{0, 1, 0, 1}
	zAxis    = mgl32.Vec4{0, 0, 1, 1}
	axes     = []mgl32.Vec4{xAxis, yAxis, zAxis}
)

type rotateToAxisAnimation struct {
	rotation            *mgl32.Mat4
	cameraEye           mgl32.Vec3
	cameraUp            mgl32.Vec3
	vectors             [3]mgl32.Vec3
	closestAxisIndexEye int
	closestAxisIndexUp  int
	closestAxisSignEye  float32
	closestAxisSignUp   float32
}

func NewRotateToAxisAnimation(rotation *mgl32.Mat4, cameraEye, cameraUp mgl32.Vec3) Animation {
	a := &rotateToAxisAnimation{
		rotation:  rotation,
		cameraEye: cameraEye.Normalize(),
		cameraUp:  cameraUp,
	}

	var (
		dotEye    [3]float32
		dotUp     [3]float32
		absDotEye [3]float32
		absDotUp  [3]float32
	)

	inverse := a.rotation.Inv()

	// For each axis,
	for i := 0; i < 3; i++ {
		// Multiply by the invert rotation matrix and normalize
		a.vectors[i] = inverse.Mul4x1(axes[i]).Vec3().Normalize()
		// And calculate how far it is from the camera vectors
		dotEye[i] = a.cameraEye.Dot(a.vectors[i])
		dotUp[i] = a.cameraUp.Dot(a.vectors[i])
		absDotEye[i] = abs(dotEye[i])
		absDotUp[i] = abs(dotUp[i])
	}

	// Determine which is the closest to camera eye
	if absDotEye[0] > absDotEye[1] && absDotEye[0] > absDotEye[2] {
		// fmt.Println("X is closest to camera eye")
		a.closestAxisIndexEye = 0
		absDotUp[0] = 0 // Make sure X cannot win up axis as well
	} else if absDotEye[1] > absDotEye[2] {
		// fmt.Println("Y is closest to camera eye")
		a.closestAxisIndexEye = 1
		absDotUp[1] = 0 // Make sure Y cannot win up axis as well
	} else {
		// fmt.Println("Z is closest to camera eye")
		a.closestAxisIndexEye = 2
		absDotUp[2] = 0 // Make sure Z cannot win up axis as well
	}

	// Determine which is the closest to camera up
	if absDotUp[0] > absDotUp[1] && absDotUp[0] > absDotUp[2] {
		// fmt.Println("X is closest to camera up")
		a.closestAxisIndexUp = 0
	} else if absDotUp[1] > absDotUp[2] {
		// fmt.Println("Y is closest to camera up")
		a.closestAxisIndexUp = 1
	} else {
		// fmt.Println("Z is closest to camera up")
		a.closestAxisIndexUp = 2
	}

	a.closestAxisSignEye = float32(math.Copysign(1, float64(dotEye[a.closestAxisIndexEye])))
	a.closestAxisSignUp = float32(math.Copysign(1, float64(dotUp[a.closestAxisIndexUp])))
	return a
}

func (a *rotateToAxisAnimation) Tick() bool {
	// Camera Eye
	a.vectors[a.closestAxisIndexEye] = a.rotation.Inv().Mul4x1(axes[a.closestAxisIndexEye]).Vec3().Mul(a.closestAxisSignEye).Normalize()

	angleEye := math.Acos(float64(a.vectors[a.closestAxisIndexEye].Dot(a.cameraEye)))

	angleEye = math.Min(maxAngle, angleEye)

	if angleEye != 0 {
		// fmt.Println("angleEye:", float32(angleEye))
		axisEye := a.cameraEye.Cross(a.vectors[a.closestAxisIndexEye]).Normalize()
		tempM := mgl32.HomogRotate3D(float32(angleEye), axisEye)
		*a.rotation = a.rotation.Mul4(tempM)
	}

	// Camera Up
	a.vectors[a.closestAxisIndexUp] = a.rotation.Inv().Mul4x1(axes[a.closestAxisIndexUp]).Vec3().Mul(a.closestAxisSignUp).Normalize()

	angleUp := math.Acos(float64(a.vectors[a.closestAxisIndexUp].Dot(a.cameraUp)))

	angleUp = math.Min(maxAngle, angleUp)

	if angleUp != 0 {
		// fmt.Println("angleUp:", float32(angleUp))
		axisUp := a.cameraUp.Cross(a.vectors[a.closestAxisIndexUp]).Normalize()
		tempM := mgl32.HomogRotate3D(float32(angleUp), axisUp)
		*a.rotation = a.rotation.Mul4(tempM)
	}

	if angleEye != 0 || angleUp != 0 {
		return false
	}

	for i := 0; i < 16; i++ {
		if a.rotation[i] > 0.5 {
			a.rotation[i] = 1.0
		} else if a.rotation[i] < -0.5 {
			a.rotation[i] = -1.0
		} else {
			a.rotation[i] = 0.0
		}
	}
	return true
}

type ReleaseBallAnimation interface {
	Animation
	Progress(float64) bool
}

type releaseBallAnimation struct {
	size        uint
	rotation    mgl32.Mat4
	player      *[3]float32
	initial     [3]float32
	goal        [3]float32
	blocks      [][3]float32
	portals     [][3]float32
	releaseAxis mgl32.Vec4
	start       time.Time
}

func NewReleaseBallAnimation(size uint, rotation mgl32.Mat4, player *[3]float32, goal [3]float32, blocks, portals [][3]float32) ReleaseBallAnimation {
	return &releaseBallAnimation{
		size:        size,
		rotation:    rotation,
		player:      player,
		initial:     [3]float32{player[0], player[1], player[2]},
		goal:        goal,
		blocks:      blocks,
		portals:     portals,
		releaseAxis: rotation.Inv().Mul4x1(mgl32.Vec4{0, -1, 0, 1}),
	}
}

func (a *releaseBallAnimation) Tick() bool {
	if a.start.IsZero() {
		a.start = time.Now()
	}
	return a.Progress(time.Now().Sub(a.start).Seconds())
}

func (a *releaseBallAnimation) Progress(time float64) bool {
	// S = ?
	// U = 0
	// V = ?
	// A = gravity
	// T = time
	// Solve for S (distance)
	// S = (U * T) + (0.5 * A * T * T)
	distance := (0 * time) + (0.5 * ACCELERATION * time * time)
	// fmt.Println("Distance:", distance)

	for j := 0; j < 3; j++ {
		a.player[j] = a.initial[j] + float32(distance)*a.releaseAxis[j]
	}
	// fmt.Println("Player", a.player)

	var (
		steps = 1 + uint(distance)
		step  uint
		limit = a.size * 10 // 10x size so player is offscreen, well out of bounds
		pos   [3]float32
		cell  [3]int
		next  [3]float32
	)

	for j := 0; j < 3; j++ {
		pos[j] = float32(math.Round(float64(a.initial[j])))
		cell[j] = int(pos[j])
		next[j] = pos[j] + a.releaseAxis[j]
	}
	for ; step < steps; step++ {
		// Check if player is out of bounds
		if step > limit {
			// fmt.Println("Player out of bounds")
			a.setPlayerCell(cell)
			return true
		}

		// Check if cell contains a goal
		if a.goal[0] == pos[0] && a.goal[1] == pos[1] && a.goal[2] == pos[2] {
			// fmt.Println("Player in Goal")
			a.setPlayerCell(cell)
			return true
		}

		if step > 0 {
			// Check if cell contains a portal
			for l, p := range a.portals {
				if p[0] == pos[0] && p[1] == pos[1] && p[2] == pos[2] {
					prev := [3]float32{
						p[0],
						p[1],
						p[2],
					}
					var pair [3]float32
					if l%2 == 0 {
						pair = a.portals[l+1]
					} else {
						pair = a.portals[l-1]
					}
					// fmt.Println("Player moved through Portal at", prev, "to", pair)
					// Move player to paired portal
					for j := 0; j < 3; j++ {
						a.player[j] += pair[j] - prev[j]
						pos[j] = float32(pair[j])
						cell[j] = int(pos[j])
						next[j] = pos[j] + a.releaseAxis[j]
					}
					break
				}
			}
		}

		// Check if next cell is blocked
		for _, b := range a.blocks {
			if b[0] == next[0] && b[1] == next[1] && b[2] == next[2] {
				// fmt.Println("Player stopped at Block")
				// TODO handle bounce
				a.setPlayerCell(cell)
				return true
			}
		}

		// Move to next cell
		for j := 0; j < 3; j++ {
			pos[j] = pos[j] + a.releaseAxis[j]
			cell[j] = int(pos[j])
			next[j] = pos[j] + a.releaseAxis[j]
		}
	}
	return false
}

type gameOverAnimation struct {
	ticks int
	model *mgl32.Mat4
}

func NewGameOverAnimation(model *mgl32.Mat4) Animation {
	return &gameOverAnimation{
		model: model,
	}
}

func (a *gameOverAnimation) Tick() bool {
	*a.model = a.model.Mul4(mgl32.Scale3D(0.75, 0.75, 0.75))
	a.ticks++
	return a.ticks >= 10
}

func (a *releaseBallAnimation) setPlayerCell(p [3]int) {
	for j := 0; j < 3; j++ {
		a.player[j] = float32(p[j])
	}
}

func abs(f float32) float32 {
	if f < 0 {
		return -f
	}
	return f
}
