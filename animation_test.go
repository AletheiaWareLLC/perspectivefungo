package perspectivefungo_test

import (
	"aletheiaware.com/perspectivefungo"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReleaseBallAnimation(t *testing.T) {
	size := uint(5)
	rotation := mgl32.Ident4()
	t.Run("Goal", func(t *testing.T) {
		player := [3]float32{0, 1, 0}
		goal := [3]float32{0, -1, 0}
		blocks := [][3]float32{}
		portals := [][3]float32{}
		a := perspectivefungo.NewReleaseBallAnimation(size, rotation, &player, goal, blocks, portals)
		// After 1 second, player should be in goal
		assert.True(t, a.Progress(1))
		assert.Equal(t, float32(0), player[0])
		assert.Equal(t, float32(-1), player[1])
		assert.Equal(t, float32(0), player[2])
	})
	t.Run("Block", func(t *testing.T) {
		player := [3]float32{0, 1, 0}
		goal := [3]float32{1, 0, 0}
		blocks := [][3]float32{
			{0, -1, 0},
		}
		portals := [][3]float32{}
		a := perspectivefungo.NewReleaseBallAnimation(size, rotation, &player, goal, blocks, portals)
		// After 1 second, player should be stopped at block
		assert.True(t, a.Progress(1))
		assert.Equal(t, float32(0), player[0])
		assert.Equal(t, float32(0), player[1])
		assert.Equal(t, float32(0), player[2])
	})
	t.Run("Portal", func(t *testing.T) {
		player := [3]float32{0, 1, 0}
		goal := [3]float32{1, -1, 0}
		blocks := [][3]float32{}
		portals := [][3]float32{
			{0, -1, 0},
			{1, 1, 0},
		}
		a := perspectivefungo.NewReleaseBallAnimation(size, rotation, &player, goal, blocks, portals)
		// After 1 second, player should be through portal and in goal
		assert.True(t, a.Progress(1))
		assert.Equal(t, float32(1), player[0])
		assert.Equal(t, float32(-1), player[1])
		assert.Equal(t, float32(0), player[2])
	})
}
