package main

import (
	"aletheiaware.com/perspectivefungo"
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
	"strings"
)

type driver struct {
	game perspectivefungo.Game

	program uint32

	projectionUniform int32
	cameraUniform     int32
	modelUniform      int32
	lightUniform      int32
	colorUniform      int32
	vertexAttrib      int32
	normalAttrib      int32

	meshes map[string]*GLMesh

	previousTime float64
}

func NewDriver() *driver {
	return &driver{
		meshes: make(map[string]*GLMesh),
	}
}

func (d *driver) Init(game perspectivefungo.Game) error {
	fmt.Println("OpenGL vendor", gl.GoStr(gl.GetString(gl.VENDOR)))
	fmt.Println("OpenGL version", gl.GoStr(gl.GetString(gl.VERSION)))
	fmt.Println("OpenGL renderer", gl.GoStr(gl.GetString(gl.RENDERER)))
	fmt.Println("OpenGL shader language version", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))

	d.game = game

	p, err := createProgram()
	if err != nil {
		return err
	}
	d.program = p
	gl.UseProgram(d.program)

	logError()

	d.projectionUniform = gl.GetUniformLocation(d.program, gl.Str("u_Projection\x00"))
	d.cameraUniform = gl.GetUniformLocation(d.program, gl.Str("u_Camera\x00"))
	d.modelUniform = gl.GetUniformLocation(d.program, gl.Str("u_Model\x00"))
	d.lightUniform = gl.GetUniformLocation(d.program, gl.Str("u_LightPos\x00"))
	d.colorUniform = gl.GetUniformLocation(d.program, gl.Str("u_Color\x00"))
	d.vertexAttrib = gl.GetAttribLocation(d.program, gl.Str("a_Position\x00"))
	d.normalAttrib = gl.GetAttribLocation(d.program, gl.Str("a_Normal\x00"))

	logError()

	gl.ClearColor(perspectivefungo.BackgroundColor[0], perspectivefungo.BackgroundColor[1], perspectivefungo.BackgroundColor[2], perspectivefungo.BackgroundColor[3])
	gl.ClearDepth(1.0)

	logError()

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	logError()

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	logError()

	gl.Enable(gl.MULTISAMPLE)

	logError()

	if err := d.game.Init(d); err != nil {
		return err
	}

	logError()
	return nil
}

func (d *driver) Loop(now float64) error {
	d.previousTime = now

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(d.program)

	logError()

	if err := d.game.Loop(d); err != nil {
		return err
	}

	logError()
	return nil
}

func (d *driver) Now() float64 {
	return d.previousTime
}

func (d *driver) SetProjection(p *mgl32.Mat4) {
	gl.UniformMatrix4fv(d.projectionUniform, 1, false, &p[0])
}

func (d *driver) SetCamera(c *mgl32.Mat4) {
	gl.UniformMatrix4fv(d.cameraUniform, 1, false, &c[0])
}

func (d *driver) SetModel(m *mgl32.Mat4) {
	gl.UniformMatrix4fv(d.modelUniform, 1, false, &m[0])
}

func (d *driver) SetLight(l *mgl32.Vec3) {
	gl.Uniform3f(d.lightUniform, l[0], l[1], l[2])
}

func (d *driver) SetColor(c *mgl32.Vec4) {
	gl.Uniform4fv(d.colorUniform, 1, &c[0])
}

func createProgram() (uint32, error) {
	program := gl.CreateProgram()

	vertex, err := compileShader(perspectivefungo.VertexShader+"\x00", gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragment, err := compileShader(perspectivefungo.FragmentShader+"\x00", gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	gl.AttachShader(program, vertex)
	gl.AttachShader(program, fragment)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("Failed to link program: %v", log)
	}

	gl.DeleteShader(vertex)
	gl.DeleteShader(fragment)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	cSource, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, cSource, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("Failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func logError() {
	if err := gl.GetError(); err != 0 {
		fmt.Println("GL Error: ", err)
	}
}
