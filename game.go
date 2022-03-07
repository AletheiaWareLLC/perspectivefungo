package perspectivefungo

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"time"
)

type Game interface {
	Init(Driver) error
	Resize(float32, float32)
	Reset()
	Start()
	Loop(Driver) error
	Rotate(float32, float32)
	RotateToAxis()
	ReleaseBall()
	Animating() bool
	Solution() *Solution
	GameOver(bool)
	HasGameStarted() bool
	HasGameEnded() bool
}

type game struct {
	puzzle   *Puzzle
	solution *Solution

	scale    mgl32.Mat4
	rotation mgl32.Mat4

	cameraEye    mgl32.Vec3
	cameraLookAt mgl32.Vec3
	cameraUp     mgl32.Vec3

	projection mgl32.Mat4
	camera     mgl32.Mat4
	model      mgl32.Mat4
	light      mgl32.Vec3

	animation Animation

	player  [3]float32
	goal    [3]float32
	blocks  [][3]float32
	portals [][3]float32

	gameStarted, gameEnded bool
}

func NewGame(puzzle *Puzzle) Game {
	fmt.Println("Playing", puzzle)
	g := &game{
		puzzle: puzzle,
	}
	g.Reset()
	return g
}

func (g *game) Init(d Driver) error {
	if err := LoadAssets(d); err != nil {
		return err
	}
	g.Reset()
	return nil
}

func (g *game) Resize(width, height float32) {
	g.projection = NewProjection(width, height)
}

func (g *game) Reset() {
	scale := 2 / float32(g.puzzle.Size)
	g.scale = mgl32.Scale3D(scale, scale, scale)
	g.rotation = mgl32.Ident4()

	g.cameraEye = NewCameraEye()
	g.cameraLookAt = NewCameraLookAt()
	g.cameraUp = NewCameraUp()
	g.camera = mgl32.LookAtV(g.cameraEye, g.cameraLookAt, g.cameraUp)

	g.model = NewModel()

	g.light = NewLight()

	g.animation = nil

	for i := 0; i < 3; i++ {
		g.player[i] = float32(g.puzzle.Player[i])
		g.goal[i] = float32(g.puzzle.Goal[i])
	}

	g.blocks = nil
	for i := 0; i < len(g.puzzle.Blocks); i += 3 {
		g.blocks = append(g.blocks, [3]float32{
			float32(g.puzzle.Blocks[i]),
			float32(g.puzzle.Blocks[i+1]),
			float32(g.puzzle.Blocks[i+2]),
		})
	}

	g.portals = nil
	for i := 0; i < len(g.puzzle.Portals); i += 3 {
		g.portals = append(g.portals, [3]float32{
			float32(g.puzzle.Portals[i]),
			float32(g.puzzle.Portals[i+1]),
			float32(g.puzzle.Portals[i+2]),
		})
	}

	g.gameStarted = false
	g.gameEnded = false
}

func (g *game) Start() {
	g.solution = &Solution{
		Start: time.Now(),
	}
	g.gameStarted = true
}

func (g *game) Loop(d Driver) error {
	if a := g.animation; a != nil && a.Tick() {
		g.animation = nil
		if g.solution != nil {
			g.solution.Progress = append(g.solution.Progress, &Snapshot{
				Time:     time.Now(),
				Position: [3]float32{g.player[0], g.player[1], g.player[2]},
			})
		}
	}

	d.SetProjection(&g.projection)
	d.SetCamera(&g.camera)
	d.SetLight(&g.light)

	var err error
	if !g.gameStarted {
		err = g.showStart(d)
	} else if g.gameEnded {
		err = g.showEnd(d)
	} else {
		err = g.showGame(d)
	}
	return err
}

func (g *game) Rotate(radX, radY float32) {
	if !g.gameStarted || g.gameEnded {
		return
	}
	// fmt.Println("Rotate:", radX, radY)
	inverse := g.rotation.Inv()

	if radY != 0 {
		// Y
		temp := inverse.Mul4x1(mgl32.Vec4{0, 1, 0, 1}).Vec3().Normalize()
		g.rotation = g.rotation.Mul4(mgl32.HomogRotate3D(radY, temp))
	}

	if radX != 0 {
		// X
		temp := inverse.Mul4x1(mgl32.Vec4{1, 0, 0, 1}).Vec3().Normalize()
		g.rotation = g.rotation.Mul4(mgl32.HomogRotate3D(radX, temp))
	}
}

func (g *game) RotateToAxis() {
	if !g.gameStarted || g.gameEnded {
		return
	}
	// fmt.Println("RotateToAxis")
	g.animation = NewRotateToAxisAnimation(&g.rotation, g.cameraEye, g.cameraUp)
}

func (g *game) ReleaseBall() {
	if !g.gameStarted || g.gameEnded {
		return
	}
	// fmt.Println("ReleaseBall")
	g.animation = NewReleaseBallAnimation(g.puzzle.Size, g.rotation, &g.player, g.goal, g.blocks, g.portals)
}

func (g *game) Animating() bool {
	return g.animation != nil
}

func (g *game) Solution() *Solution {
	return g.solution
}

func (g *game) GameOver(won bool) {
	if won {
		g.solution.End = time.Now()
	} else {
		g.solution = nil
	}
	g.gameEnded = true
	g.animation = NewGameOverAnimation(&g.model)
}

func (g *game) HasGameStarted() bool {
	return g.gameStarted
}

func (g *game) HasGameEnded() bool {
	return g.gameEnded
}

func (g *game) showStart(d Driver) error {
	d.SetColor(&GameLostColor)
	temp := g.model.Mul4(mgl32.Scale3D(0.1, 0.1, 0.1)).Mul4(mgl32.Translate3D(0, -1, 0))
	d.SetModel(&temp)
	if err := d.DrawMesh("start"); err != nil {
		return err
	}
	return nil
}

func (g *game) showEnd(d Driver) error {
	var temp mgl32.Mat4
	if g.solution == nil {
		d.SetColor(&GameLostColor)
		temp = g.model.Mul4(mgl32.Translate3D(0, 0, 25)).Mul4(mgl32.Scale3D(1, 1, 0.1))
		d.SetModel(&temp)
		if err := d.DrawMesh("gameover"); err != nil {
			return err
		}

		// Draw retry button
		d.SetColor(&RetryColor)
		temp = g.model.Mul4(mgl32.Translate3D(0, -12, 25)).Mul4(mgl32.Scale3D(1, 1, 0.1))
		d.SetModel(&temp)
		if err := d.DrawMesh("retry"); err != nil {
			return err
		}
	} else {
		d.SetColor(&GoalColor)
		now := float32(d.Now())
		temp = g.model.Mul4(mgl32.Scale3D(40, 40, 40)).Mul4(mgl32.HomogRotate3D(now, mgl32.Vec3{1, 0, 0})).Mul4(mgl32.HomogRotate3D(now, mgl32.Vec3{0, 1, 0}))
		d.SetModel(&temp)
		if err := d.DrawMesh("goal"); err != nil {
			return err
		}

		// Show time taken
		d.SetColor(&GameWonColor)
		taken := fmt.Sprintf("%.2fs", g.solution.End.Sub(g.solution.Start).Seconds())
		offset := float32(len(taken)-1) / 2
		for i, c := range taken {
			temp = g.model.Mul4(mgl32.Translate3D((float32(i)-offset)*3.5, -2, 25))
			d.SetModel(&temp)
			if err := d.DrawMesh(string(c)); err != nil {
				return err
			}
		}

		// Draw retry button
		d.SetColor(&RetryColor)
		temp = g.model.Mul4(mgl32.Translate3D(-12, -12, 25)).Mul4(mgl32.Scale3D(1, 1, 0.1))
		d.SetModel(&temp)
		if err := d.DrawMesh("retry"); err != nil {
			return err
		}

		// Draw share button
		d.SetColor(&ShareColor)
		temp = g.model.Mul4(mgl32.Translate3D(12, -12, 25)).Mul4(mgl32.Scale3D(1, 1, 0.1))
		d.SetModel(&temp)
		if err := d.DrawMesh("share"); err != nil {
			return err
		}
	}
	return nil
}

func (g *game) showGame(d Driver) error {
	var temp mgl32.Mat4
	r := g.model.Mul4(g.scale).Mul4(g.rotation)

	for i, p := range g.portals {
		d.SetColor(&PortalColors[i/2])

		temp = r.Mul4(mgl32.Translate3D(p[0], p[1], p[2]))
		d.SetModel(&temp)

		if err := d.DrawMesh("portal"); err != nil {
			return err
		}
	}

	d.SetColor(&BlockColor)
	for _, b := range g.blocks {
		temp = r.Mul4(mgl32.Translate3D(b[0], b[1], b[2]))
		d.SetModel(&temp)

		if err := d.DrawMesh("block"); err != nil {
			return err
		}
	}

	d.SetColor(&GoalColor)
	temp = r.Mul4(mgl32.Translate3D(g.goal[0], g.goal[1], g.goal[2]))
	d.SetModel(&temp)

	if err := d.DrawMesh("goal"); err != nil {
		return err
	}

	d.SetColor(&PlayerColor)
	temp = r.Mul4(mgl32.Translate3D(g.player[0], g.player[1], g.player[2]))
	d.SetModel(&temp)

	if err := d.DrawMesh("player"); err != nil {
		return err
	}
	limit := float32(g.puzzle.Size)
	if abs(g.player[0]) > limit || abs(g.player[1]) > limit || abs(g.player[2]) > limit {
		g.GameOver(false)
	}

	if g.goal[0] == g.player[0] && g.goal[1] == g.player[1] && g.goal[2] == g.player[2] {
		g.GameOver(true)
	}

	return nil
}

func NewProjection(width, height float32) mgl32.Mat4 {
	// return mgl32.Perspective(mgl32.DegToRad(45.0), width/height, 0.1, 10.0)

	var (
		left   float32
		right  float32
		bottom float32
		top    float32
		near   float32
		far    float32
	)
	if width > height {
		ratio := width / height
		// The height will stay the same while the width will vary as per aspect ratio.
		left = -ratio
		right = ratio
		bottom = -1
		top = 1
	} else {
		ratio := height / width
		// The width will stay the same while the height will vary as per aspect ratio.
		left = -1
		right = 1
		bottom = -ratio
		top = ratio
	}
	near = 1
	far = 5

	return mgl32.Frustum(left, right, bottom, top, near, far)
}

func NewCameraEye() mgl32.Vec3 {
	return mgl32.Vec3{0, 0, 3}
}

func NewCameraLookAt() mgl32.Vec3 {
	return mgl32.Vec3{0, 0, 0}
}

func NewCameraUp() mgl32.Vec3 {
	return mgl32.Vec3{0, 1, 0}
}

func NewModel() mgl32.Mat4 {
	return mgl32.Ident4()
}

func NewLight() mgl32.Vec3 {
	return mgl32.Vec3{0, 1, 1}
}
