package main

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
)

type GLMesh struct {
	Count        int32
	VertexBuffer uint32
	NormalBuffer uint32
}

func (d *driver) AddMesh(id string, count int32, vertices, normals []float32) error {
	mesh := &GLMesh{
		Count: count,
	}

	size := int(mesh.Count * 3 * 4) // *3:xyz, *4:sizeof(float32)

	gl.GenBuffers(1, &mesh.VertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.VertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, size, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(uint32(d.vertexAttrib))
	gl.VertexAttribPointerWithOffset(uint32(d.vertexAttrib), 3, gl.FLOAT, false, 0, 0)

	gl.GenBuffers(1, &mesh.NormalBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.NormalBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, size, gl.Ptr(normals), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(uint32(d.normalAttrib))
	gl.VertexAttribPointerWithOffset(uint32(d.normalAttrib), 3, gl.FLOAT, false, 0, 0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	logError()

	d.meshes[id] = mesh
	return nil
}

func (d *driver) DrawMesh(id string) error {
	m, ok := d.meshes[id]
	if !ok {
		return fmt.Errorf("Unrecognized mesh: %s", id)
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, m.VertexBuffer)
	gl.EnableVertexAttribArray(uint32(d.vertexAttrib))
	gl.VertexAttribPointerWithOffset(uint32(d.vertexAttrib), 3, gl.FLOAT, false, 0, 0)

	gl.BindBuffer(gl.ARRAY_BUFFER, m.NormalBuffer)
	gl.EnableVertexAttribArray(uint32(d.normalAttrib))
	gl.VertexAttribPointerWithOffset(uint32(d.normalAttrib), 3, gl.FLOAT, false, 0, 0)

	gl.DrawArrays(gl.TRIANGLES, 0, m.Count)

	gl.DisableVertexAttribArray(uint32(d.vertexAttrib))
	gl.DisableVertexAttribArray(uint32(d.normalAttrib))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return nil
}
