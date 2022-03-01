package main

import (
	"aletheiaware.com/perspectivefungo"
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const DAILY_URL = "http://localhost/daily.json"

var (
	width  = 800
	height = 600
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	/*
		width = 200
		height = 200
		game := &Recorder{
			Output: os.Args[1],
			Mesh:   "goal",
			Color:  perspectivefungo.GoalColor,
			Count:  64 - 1,
		}
	*/

	if err := glfw.Init(); err != nil {
		log.Fatal("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(width, height, "Perspective", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}

	d := NewDriver()

	puzzle, err := loadPuzzle()
	if err != nil {
		log.Fatal(err)
	}
	game := perspectivefungo.NewGame(puzzle)

	if err := d.Init(game); err != nil {
		log.Fatal(err)
	}

	game.Resize(float32(width), float32(height))

	window.SetCloseCallback(func(w *glfw.Window) {
		// Do Nothing
	})

	window.SetPosCallback(func(w *glfw.Window, x, y int) {
		// Do Nothing
	})

	scale := math.Min(float64(width), float64(height))
	threshold := scale / 100

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
		if game != nil {
			game.Resize(float32(width), float32(height))
		}
		scale = math.Min(float64(width), float64(height))
		threshold = scale / 100
	})

	window.SetRefreshCallback(func(w *glfw.Window) {
		// Do Nothing
	})

	var (
		dragging     bool
		rotated      bool
		lastX, lastY float64
	)
	window.SetCursorPosCallback(func(w *glfw.Window, x, y float64) {
		if game == nil || game.Animating() {
			return
		}
		if dragging {
			deltaX := x - lastX
			deltaY := y - lastY
			if math.Abs(deltaX) > threshold || math.Abs(deltaY) > threshold {
				rotated = true
				radX := float32((deltaY / scale) * 2.0 * math.Pi)
				radY := float32((deltaX / scale) * 2.0 * math.Pi)
				game.Rotate(radX, radY)
			} else {
				// Don't update lastX or lastY
				return
			}
		}
		lastX, lastY = x, y
	})

	window.SetMouseButtonCallback(func(w *glfw.Window, btn glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if btn != glfw.MouseButton1 {
			// Do Nothing
			return
		}
		if game == nil || game.Animating() {
			return
		}
		switch action {
		case glfw.Press:
			dragging = true
		case glfw.Release:
			dragging = false
			if rotated {
				game.RotateToAxis()
			} else if !game.HasGameStarted() {
				game.Start()
			} else if game.HasGameEnded() {
				s := game.Solution()
				if s == nil {
					game.Reset()
					game.Start()
				} else {
					if lastX < float64(width)/2 {
						game.Reset()
						game.Start()
					} else {
						log.Println("Not Supported")
					}
				}
			} else {
				game.ReleaseBall()
			}
		}
		rotated = false
	})

	for !window.ShouldClose() {
		if err := d.Loop(glfw.GetTime()); err != nil {
			log.Fatal(err)
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func loadPuzzle() (*perspectivefungo.Puzzle, error) {
	var reader io.Reader
	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			return nil, err
		}
		defer file.Close()
		reader = file
	} else {
		response, err := http.Get(DAILY_URL)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()
		reader = response.Body
	}
	var puzzle perspectivefungo.Puzzle
	if err := json.NewDecoder(reader).Decode(&puzzle); err != nil {
		return nil, err
	}
	return &puzzle, nil
}