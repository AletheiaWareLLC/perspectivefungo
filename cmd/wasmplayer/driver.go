package main

import (
	"aletheiaware.com/perspectivefungo"
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/justinclift/webgl"
	"syscall/js"
	"unsafe"
)

type driver struct {
	game perspectivefungo.Game

	gl      js.Value
	program js.Value

	projectionUniform js.Value
	cameraUniform     js.Value
	modelUniform      js.Value
	lightUniform      js.Value
	colorUniform      js.Value
	vertexAttrib      js.Value
	normalAttrib      js.Value

	meshes map[string]*GLMesh

	previousTime float64
}

func NewDriver(gl js.Value) *driver {
	return &driver{
		gl:     gl,
		meshes: make(map[string]*GLMesh),
	}
}

func (d *driver) Init(game perspectivefungo.Game) error {
	d.game = game

	if err := loadAssets(); err != nil {
		return err
	}

	p, err := createProgram(d.gl)
	if err != nil {
		return err
	}
	d.program = p
	d.gl.Call("useProgram", d.program)

	logError(d.gl)

	d.projectionUniform = d.gl.Call("getUniformLocation", d.program, "u_Projection")
	d.cameraUniform = d.gl.Call("getUniformLocation", d.program, "u_Camera")
	d.modelUniform = d.gl.Call("getUniformLocation", d.program, "u_Model")
	d.lightUniform = d.gl.Call("getUniformLocation", d.program, "u_LightPos")
	d.colorUniform = d.gl.Call("getUniformLocation", d.program, "u_Color")
	d.vertexAttrib = d.gl.Call("getAttribLocation", d.program, "a_Position")
	d.gl.Call("enableVertexAttribArray", d.vertexAttrib)
	d.gl.Call("vertexAttribPointer", d.vertexAttrib, 3, webgl.FLOAT, false, 0, 0)
	d.normalAttrib = d.gl.Call("getAttribLocation", d.program, "a_Normal")
	d.gl.Call("enableVertexAttribArray", d.normalAttrib)
	d.gl.Call("vertexAttribPointer", d.normalAttrib, 3, webgl.FLOAT, false, 0, 0)

	logError(d.gl)

	d.gl.Call("clearColor", perspectivefungo.BackgroundColor[0], perspectivefungo.BackgroundColor[1], perspectivefungo.BackgroundColor[2], perspectivefungo.BackgroundColor[3])
	d.gl.Call("clearDepth", 1.0)
	d.gl.Call("enable", webgl.DEPTH_TEST)
	d.gl.Call("depthFunc", webgl.LEQUAL)
	d.gl.Call("enable", webgl.BLEND)
	d.gl.Call("blendFunc", webgl.SRC_ALPHA, webgl.ONE_MINUS_SRC_ALPHA)

	logError(d.gl)

	if err := d.game.Init(d); err != nil {
		return err
	}

	logError(d.gl)
	return nil
}

func (d *driver) Loop(now float64) error {
	now /= 1000
	d.previousTime = now

	d.gl.Call("clear", webgl.COLOR_BUFFER_BIT)
	d.gl.Call("clear", webgl.DEPTH_BUFFER_BIT)

	d.gl.Call("useProgram", d.program)

	logError(d.gl)

	if err := d.game.Loop(d); err != nil {
		return err
	}

	logError(d.gl)
	return nil
}

func (d *driver) Now() float64 {
	return d.previousTime
}

func (d *driver) SetProjection(p *mgl32.Mat4) {
	var projectionBuffer *[16]float32
	projectionBuffer = (*[16]float32)(unsafe.Pointer(p))
	typedProjectionBuffer := webgl.SliceToTypedArray([]float32((*projectionBuffer)[:]))
	d.gl.Call("uniformMatrix4fv", d.projectionUniform, false, typedProjectionBuffer)
}

func (d *driver) SetCamera(c *mgl32.Mat4) {
	var cameraBuffer *[16]float32
	cameraBuffer = (*[16]float32)(unsafe.Pointer(c))
	typedCameraBuffer := webgl.SliceToTypedArray([]float32((*cameraBuffer)[:]))
	d.gl.Call("uniformMatrix4fv", d.cameraUniform, false, typedCameraBuffer)
}

func (d *driver) SetModel(m *mgl32.Mat4) {
	var modelBuffer *[16]float32
	modelBuffer = (*[16]float32)(unsafe.Pointer(m))
	typedModelBuffer := webgl.SliceToTypedArray([]float32((*modelBuffer)[:]))
	d.gl.Call("uniformMatrix4fv", d.modelUniform, false, typedModelBuffer)
}

func (d *driver) SetLight(l *mgl32.Vec3) {
	d.gl.Call("uniform3f", d.lightUniform, l[0], l[1], l[2])
}

func (d *driver) SetColor(c *mgl32.Vec4) {
	d.gl.Call("uniform4f", d.colorUniform, c[0], c[1], c[2], c[3])
}

func createProgram(gl js.Value) (js.Value, error) {
	program := gl.Call("createProgram")

	vertex, err := compileShader(gl, perspectivefungo.VertexShader, webgl.VERTEX_SHADER)
	if err != nil {
		return js.Undefined(), err
	}

	fragment, err := compileShader(gl, perspectivefungo.FragmentShader, webgl.FRAGMENT_SHADER)
	if err != nil {
		return js.Undefined(), err
	}

	gl.Call("attachShader", program, vertex)
	gl.Call("attachShader", program, fragment)
	gl.Call("linkProgram", program)

	if !gl.Call("getProgramParameter", program, webgl.LINK_STATUS).Bool() {
		log := gl.Call("getProgramInfoLog", program)
		fmt.Println("Log:", log)
		return js.Undefined(), fmt.Errorf("Failed to link program: %v", log)
	}

	gl.Call("deleteShader", vertex)
	gl.Call("deleteShader", fragment)

	return program, nil
}

func compileShader(gl js.Value, source string, shaderType int) (js.Value, error) {
	shader := gl.Call("createShader", shaderType)

	gl.Call("shaderSource", shader, source)
	gl.Call("compileShader", shader)

	if !gl.Call("getShaderParameter", shader, webgl.COMPILE_STATUS).Bool() {
		log := gl.Call("getShaderInfoLog", shader)
		fmt.Println("Log:", log)
		return js.Undefined(), fmt.Errorf("Failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func logError(gl js.Value) {
	if err := gl.Call("getError").Int(); err != 0 {
		fmt.Println("GL Error: ", err)
	}
}
