package main

import (
	"fmt"
	"github.com/justinclift/webgl"
	"syscall/js"
)

type GLMesh struct {
	Count        int32
	VertexBuffer *js.Value
	NormalBuffer *js.Value
}

func (d *driver) AddMesh(id string, count int32, vertices, normals []float32) error {
	mesh := &GLMesh{
		Count: count,
	}

	vb := d.gl.Call("createBuffer")
	mesh.VertexBuffer = &vb
	d.gl.Call("bindBuffer", webgl.ARRAY_BUFFER, mesh.VertexBuffer)
	d.gl.Call("bufferData", webgl.ARRAY_BUFFER, webgl.SliceToTypedArray(vertices), webgl.STATIC_DRAW)

	nb := d.gl.Call("createBuffer")
	mesh.NormalBuffer = &nb
	d.gl.Call("bindBuffer", webgl.ARRAY_BUFFER, mesh.NormalBuffer)
	d.gl.Call("bufferData", webgl.ARRAY_BUFFER, webgl.SliceToTypedArray(normals), webgl.STATIC_DRAW)

	logError(d.gl)

	d.meshes[id] = mesh
	return nil
}

func (d *driver) DrawMesh(id string) error {
	m, ok := d.meshes[id]
	if !ok {
		return fmt.Errorf("Unrecognized mesh: %s", id)
	}

	d.gl.Call("bindBuffer", webgl.ARRAY_BUFFER, m.VertexBuffer)
	d.gl.Call("enableVertexAttribArray", d.vertexAttrib)
	d.gl.Call("vertexAttribPointer", d.vertexAttrib, 3, webgl.FLOAT, false, 0, 0)

	d.gl.Call("bindBuffer", webgl.ARRAY_BUFFER, m.NormalBuffer)
	d.gl.Call("enableVertexAttribArray", d.normalAttrib)
	d.gl.Call("vertexAttribPointer", d.normalAttrib, 3, webgl.FLOAT, false, 0, 0)

	d.gl.Call("drawArrays", webgl.TRIANGLES, 0, m.Count)

	d.gl.Call("disableVertexAttribArray", d.vertexAttrib)
	d.gl.Call("disableVertexAttribArray", d.normalAttrib)

	logError(d.gl)
	return nil
}
