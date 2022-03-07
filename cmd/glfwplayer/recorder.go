package main

import (
	"aletheiaware.com/perspectivefungo"
	"errors"
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"os"
)

type Recorder struct {
	Output        string
	Mesh          string
	Color         mgl32.Vec4
	Count         int
	width, height int
	cameraEye     mgl32.Vec3
	cameraLookAt  mgl32.Vec3
	cameraUp      mgl32.Vec3
	projection    mgl32.Mat4
	camera        mgl32.Mat4
	model         mgl32.Mat4
	light         mgl32.Vec3
	frame         int
	frames        [][]uint8
}

func (r *Recorder) Init(d perspectivefungo.Driver) error {
	fmt.Println("Init")
	if err := perspectivefungo.LoadAssets(d); err != nil {
		return err
	}
	r.Reset()
	return nil
}

func (r *Recorder) Resize(width, height float32) {
	fmt.Println(width, height)
	r.width = int(width)
	r.height = int(height)
	r.projection = perspectivefungo.NewProjection(width, height)
}

func (r *Recorder) Reset() {
	r.cameraEye = perspectivefungo.NewCameraEye()
	r.cameraLookAt = perspectivefungo.NewCameraLookAt()
	r.cameraUp = perspectivefungo.NewCameraUp()
	r.camera = mgl32.LookAtV(r.cameraEye, r.cameraLookAt, r.cameraUp)
	r.model = perspectivefungo.NewModel()
	r.light = perspectivefungo.NewLight()
}

func (r *Recorder) Start() {
	//
}

func (r *Recorder) Loop(d perspectivefungo.Driver) error {
	length := r.width * r.height * 4

	if r.frame > 0 {
		frame := make([]uint8, length)
		r.frames = append(r.frames, frame)
		gl.ReadBuffer(gl.FRONT)
		gl.ReadPixels(int32(r.width)/2, int32(r.height)/2, int32(r.width), int32(r.height), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(frame))
	}

	if r.frame >= r.Count {
		switch r.Count {
		case 1:
			frame := r.frames[0]
			img := image.NewRGBA(image.Rect(0, 0, r.width, r.height))
			for x := 0; x < r.width; x++ {
				for y := 0; y < r.height; y++ {
					i := ((r.height-y-1)*r.width + x) * 4
					img.Set(x, y, color.RGBA{frame[i], frame[i+1], frame[i+2], frame[i+3]})
				}
			}
			// Write Output
			f, err := os.OpenFile(r.Output, os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				return err
			}
			if err := png.Encode(f, img); err != nil {
				return err
			}
			if err := f.Close(); err != nil {
				return err
			}
		default:
			// Create Gif Palette
			downsize := make(map[string]color.Color)
			colors := make(map[string]color.Color)
			for _, f := range r.frames {
				for i := 0; i < length; i += 4 {
					r := f[i]
					g := f[i+1]
					b := f[i+2]
					a := f[i+3]
					original := color.RGBA{R: r, G: g, B: b, A: a}
					if r == 255 && g == 255 && b == 255 && a == 255 {
						downsize[fmt.Sprintf("%d,%d,%d,%d", original.R, original.G, original.B, original.A)] = original
						colors[fmt.Sprintf("%d,%d,%d,%d", original.R, original.G, original.B, original.A)] = original
					} else {
						limited := color.RGBA{
							R: (r >> 3) << 3,
							G: (g >> 3) << 3,
							B: (b >> 3) << 3,
							A: (a >> 3) << 3,
						}
						downsize[fmt.Sprintf("%d,%d,%d,%d", original.R, original.G, original.B, original.A)] = limited
						colors[fmt.Sprintf("%d,%d,%d,%d", limited.R, limited.G, limited.B, limited.A)] = limited
					}
				}
			}
			fmt.Println("Colors:", len(colors))

			var palette []color.Color
			for _, c := range colors {
				palette = append(palette, c)
			}

			// Create Gif Frames
			var images []*image.Paletted
			var delays []int
			for _, f := range r.frames {
				img := image.NewPaletted(image.Rect(0, 0, r.width, r.height), palette)
				images = append(images, img)
				delays = append(delays, 1)
				for x := 0; x < r.width; x++ {
					for y := 0; y < r.height; y++ {
						i := ((r.height-y-1)*r.width + x) * 4
						img.Set(x, y, downsize[fmt.Sprintf("%d,%d,%d,%d", f[i], f[i+1], f[i+2], f[i+3])])
					}
				}
			}

			// Write Output
			f, err := os.OpenFile(r.Output, os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				return err
			}
			if err := gif.EncodeAll(f, &gif.GIF{
				Image: images,
				Delay: delays,
			}); err != nil {
				return err
			}
			if err := f.Close(); err != nil {
				return err
			}
		}
		return errors.New("Done")
	}

	d.SetProjection(&r.projection)
	d.SetCamera(&r.camera)
	d.SetLight(&r.light)

	d.SetColor(&r.Color)
	angle := float32(r.frame) / 10
	temp := r.model.Mul4(mgl32.Scale3D(2, 2, 2)).Mul4(mgl32.HomogRotate3D(angle, mgl32.Vec3{1, 0, 0})).Mul4(mgl32.HomogRotate3D(angle, mgl32.Vec3{0, 1, 0}))
	d.SetModel(&temp)
	if err := d.DrawMesh(r.Mesh); err != nil {
		return err
	}

	r.frame++
	return nil
}

func (r *Recorder) Rotate(float32, float32) {
	// Do Nothing
}

func (r *Recorder) RotateToAxis() {
	// Do Nothing
}

func (r *Recorder) ReleaseBall() {
	// Do Nothing
}

func (r *Recorder) Animating() bool {
	return false
}

func (r *Recorder) GameOver(bool) {
	// Do Nothing
}

func (r *Recorder) Solution() *perspectivefungo.Solution {
	return nil
}

func (r *Recorder) HasGameStarted() bool {
	return false
}

func (r *Recorder) HasGameEnded() bool {
	return false
}
