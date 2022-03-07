package perspectivefungo

import (
	"github.com/go-gl/mathgl/mgl32"
)

var (
	BackgroundColor = mgl32.Vec4{1, 1, 1, 1}
	GameLostColor   = mgl32.Vec4{0.9, 0, 0, 1}
	GameWonColor    = mgl32.Vec4{0.9, 0.9, 0, 1}
	RetryColor      = mgl32.Vec4{0.9, 0, 0, 1}
	ShareColor      = mgl32.Vec4{0, 0.8, 0, 1}
	BlockColor      = mgl32.Vec4{0.5, 0.5, 0.5, 1}
	GoalColor       = mgl32.Vec4{0, 0.8, 0, 0.9}
	PlayerColor     = mgl32.Vec4{1, 1, 0, 1}
	PortalColors    = []mgl32.Vec4{
		mgl32.Vec4{0, 0, 0.8, 0.9},
		mgl32.Vec4{0.27, 0, 0.8, 0.9},
		mgl32.Vec4{0.54, 0, 0.8, 0.9},
	}
)
