package perspectivefungo

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Driver interface {
	Init(Game) error
	Loop(float64) error
	AddMesh(string, int32, []float32, []float32) error
	DrawMesh(string) error
	Now() float64
	SetProjection(*mgl32.Mat4)
	SetCamera(*mgl32.Mat4)
	SetModel(*mgl32.Mat4)
	SetLight(*mgl32.Vec3)
	SetColor(*mgl32.Vec4)
}
