package main

import (
	"aletheiaware.com/perspectivefungo"
	"fmt"
	"log"
	"math"
	"net/url"
	"strconv"
	"syscall/js"
)

var (
	window       js.Value
	document     js.Value
	canvas       js.Value
	gl           js.Value
	d            *driver
	g            perspectivefungo.Game
	width        float64
	height       float64
	scale        float64
	threshold    float64
	dragging     bool
	rotated      bool
	lastX, lastY float64
)

func main() {
	js.Global().Set("renderFrame", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		renderFrame(args[0].Float())
		return nil
	}))
	js.Global().Set("startPuzzle", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if err := startPuzzle(args[0]); err != nil {
			js.Global().Call("alert", err.Error())
		}
		return nil
	}))
	window = js.Global().Get("window")
	document := js.Global().Get("document")
	canvas := document.Call("getElementById", "gocanvas")

	canvas.Call("addEventListener", "pointerdown", js.FuncOf(handleDown))
	canvas.Call("addEventListener", "pointercancel", js.FuncOf(handleCancel))
	canvas.Call("addEventListener", "pointermove", js.FuncOf(handleMove))
	canvas.Call("addEventListener", "pointerup", js.FuncOf(handleUp))

	gl = canvas.Call("getContext", "webgl", "{antialias: true}")
	if gl.IsUndefined() {
		js.Global().Call("alert", "Your browser doesn't appear to support WebGL")
		return
	}

	resizeCanvas := func() {
		width = canvas.Get("clientWidth").Float()
		height = canvas.Get("clientHeight").Float()
		canvas.Call("setAttribute", "width", width)
		canvas.Call("setAttribute", "height", height)

		if g != nil {
			g.Resize(float32(width), float32(height))
		}

		gl.Call("viewport", 0, 0, width, height)

		scale = math.Min(width, height)
		threshold = scale / 100
	}

	window.Call("addEventListener", "resize", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resizeCanvas()
		return nil
	}))

	resizeCanvas()

	d = NewDriver(gl)

	<-make(chan struct{})
}

func renderFrame(now float64) {
	if err := d.Loop(now); err != nil {
		log.Fatal(err)
	}

	loop()
}

func loop() {
	js.Global().Call("requestAnimationFrame", js.Global().Get("renderFrame"))
}

func startPuzzle(data js.Value) error {
	puzzle := &perspectivefungo.Puzzle{}

	puzzle.Size = uint(data.Get("size").Int())

	player := data.Get("player")
	if !player.IsUndefined() && !player.IsNull() {
		for i := 0; i < player.Get("length").Int(); i++ {
			puzzle.Player = append(puzzle.Player, player.Get(strconv.Itoa(i)).Int())
		}
	}

	goal := data.Get("goal")
	if !goal.IsUndefined() && !goal.IsNull() {
		for i := 0; i < goal.Get("length").Int(); i++ {
			puzzle.Goal = append(puzzle.Goal, goal.Get(strconv.Itoa(i)).Int())
		}
	}

	blocks := data.Get("blocks")
	if !blocks.IsUndefined() && !blocks.IsNull() {
		for i := 0; i < blocks.Get("length").Int(); i++ {
			puzzle.Blocks = append(puzzle.Blocks, blocks.Get(strconv.Itoa(i)).Int())
		}
	}

	portals := data.Get("portals")
	if !portals.IsUndefined() && !portals.IsNull() {
		for i := 0; i < portals.Get("length").Int(); i++ {
			puzzle.Portals = append(puzzle.Portals, portals.Get(strconv.Itoa(i)).Int())
		}
	}

	g = perspectivefungo.NewGame(puzzle)

	if err := d.Init(g); err != nil {
		log.Fatal(err)
	}

	g.Resize(float32(width), float32(height))

	gl.Call("viewport", 0, 0, width, height)

	loop()
	return nil
}

func handleDown(this js.Value, args []js.Value) interface{} {
	event := args[0]
	if !event.Get("isPrimary").Bool() {
		// Ignore
		return nil
	}
	if g != nil && g.Animating() {
		return nil
	}
	x := event.Get("clientX").Float()
	y := event.Get("clientY").Float()
	dragging = true
	rotated = false
	lastX, lastY = x, y
	return nil
}

func handleCancel(this js.Value, args []js.Value) interface{} {
	event := args[0]
	if !event.Get("isPrimary").Bool() {
		// Ignore
		return nil
	}
	if g != nil && g.Animating() {
		return nil
	}
	dragging = false
	rotated = false
	g.RotateToAxis()
	return nil
}

func handleMove(this js.Value, args []js.Value) interface{} {
	event := args[0]
	if !event.Get("isPrimary").Bool() {
		// Ignore
		return nil
	}
	if g != nil && g.Animating() {
		return nil
	}
	x := event.Get("clientX").Float()
	y := event.Get("clientY").Float()
	if dragging {
		deltaX := x - lastX
		deltaY := y - lastY
		if math.Abs(deltaX) > threshold || math.Abs(deltaY) > threshold {
			rotated = true
			radX := float32((deltaY / scale) * 2.0 * math.Pi)
			radY := float32((deltaX / scale) * 2.0 * math.Pi)
			g.Rotate(radX, radY)
		} else {
			// Don't update lastX or lastY
			return nil
		}
	}
	lastX, lastY = x, y
	return nil
}

func handleUp(this js.Value, args []js.Value) interface{} {
	event := args[0]
	if !event.Get("isPrimary").Bool() {
		// Ignore
		return nil
	}
	if g != nil && g.Animating() {
		return nil
	}
	dragging = false
	if rotated {
		g.RotateToAxis()
	} else if !g.HasGameStarted() {
		g.Start()
	} else if g.HasGameEnded() {
		s := g.Solution()
		if s == nil {
			g.Reset()
			g.Start()
		} else {
			x := event.Get("clientX").Float()
			if x < width/2 {
				g.Reset()
				g.Start()
			} else {
				params := url.Values{}
				params.Add("url", "https://perspective.fun/daily")
				params.Add("text", fmt.Sprintf("%s %.2fs\n\n", s.Start.UTC().Format("2006-01-02"), s.End.Sub(s.Start).Seconds()))
				params.Add("hashtags", "PerspectiveDailyPuzzle")
				window.Get("location").Set("href", "https://twitter.com/intent/tweet?"+params.Encode())
			}
		}
	} else {
		g.ReleaseBall()
	}
	rotated = false
	return nil
}
