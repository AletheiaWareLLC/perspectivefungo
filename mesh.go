package perspectivefungo

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func LoadOFFMesh(d Driver, id string, data []byte, smooth bool) error {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	if !scanner.Scan() {
		return fmt.Errorf("Failed to read .OFF header")
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	header := strings.Fields(scanner.Text())
	if header[0] != "OFF" {
		return fmt.Errorf("Invalid .OFF header")
	}
	vertexCount, err := strconv.Atoi(header[1])
	if err != nil {
		return err
	}
	faceCount, err := strconv.Atoi(header[2])
	if err != nil {
		return err
	}

	var (
		count    int32
		vertices []float32
		normals  []float32

		tempVertices [][]float32
		tempNormals  [][]float32
		tempFaces    [][3]int
	)

	addOffFace := func(ss ...string) error {
		var face [3]int
		var vs [3][]float32
		for i := 0; i < 3; i++ {
			f, err := strconv.Atoi(ss[i])
			if err != nil {
				return err
			}
			face[i] = f
			vs[i] = tempVertices[face[i]]
			vertices = append(vertices, vs[i]...)
			count++
		}
		tempFaces = append(tempFaces, face)
		// calculate face edges
		e0 := [3]float32{
			vs[1][0] - vs[0][0],
			vs[1][1] - vs[0][1],
			vs[1][2] - vs[0][2],
		}
		e1 := [3]float32{
			vs[2][0] - vs[0][0],
			vs[2][1] - vs[0][1],
			vs[2][2] - vs[0][2],
		}
		// calculate face normal
		n0 := e0[1]*e1[2] - e0[2]*e1[1]
		n1 := e0[2]*e1[0] - e0[0]*e1[2]
		n2 := e0[0]*e1[1] - e0[1]*e1[0]
		// normalize face normal
		length := float32(math.Sqrt(float64((n0 * n0) + (n1 * n1) + (n2 * n2))))
		if length > 0 {
			n0 = n0 / length
			n1 = n1 / length
			n2 = n2 / length
		}
		n := []float32{
			n0,
			n1,
			n2,
		}
		tempNormals = append(tempNormals, n)
		return nil
	}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if vertexCount > 0 {
			var v []float32
			for _, s := range parts[0:3] {
				f, err := strconv.ParseFloat(s, 32)
				if err != nil {
					return err
				}
				v = append(v, float32(f))
			}
			tempVertices = append(tempVertices, v)
			vertexCount--
		} else if faceCount > 0 {
			count, err := strconv.Atoi(parts[0])
			if err != nil {
				return err
			}
			if count > 3 {
				// Convert Polygon to Triangles
				for i := 1; i < count-1; i++ {
					if err := addOffFace(parts[1], parts[i+1], parts[i+2]); err != nil {
						return err
					}
				}
			} else {
				// Handle triangle
				if err := addOffFace(parts[1:4]...); err != nil {
					return err
				}
			}
			faceCount--
		} else {
			fmt.Println("Ignoring:", line)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if smooth {
		vns := make(map[int][]float32, count)
		// Loop once to add face normals to vertex normals
		for f, face := range tempFaces {
			fn := tempNormals[f]
			for i := 0; i < 3; i++ {
				vn := vns[face[i]]
				if vn == nil {
					vn = []float32{0, 0, 0}
				}
				vn[0] += fn[0]
				vn[1] += fn[1]
				vn[2] += fn[2]
				vns[face[i]] = vn
			}
		}
		// Loop again to add tempNormals to normals
		for _, face := range tempFaces {
			for i := 0; i < 3; i++ {
				vn := vns[face[i]]
				// Normalize normal
				length := float32(math.Sqrt(float64((vn[0] * vn[0]) + (vn[1] * vn[1]) + (vn[2] * vn[2]))))
				if length > 0 {
					vn[0] /= length
					vn[1] /= length
					vn[2] /= length
				}
				normals = append(normals, vn...)
			}
		}
	} else {
		for _, temp := range tempNormals {
			// Add once for each vertex
			for i := 0; i < 3; i++ {
				normals = append(normals, temp...)
			}
		}
	}
	if err := d.AddMesh(id, count, vertices, normals); err != nil {
		return err
	}
	return nil
}
