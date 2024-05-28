package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// Cube struct
type Cube struct {
	Position *Vector3
	Points   []*Vector3
	Colors   []Color
	Indices  []uint16
	Vertices []ebiten.Vertex
}

// NewCube returns a &Cube
func NewCube(x, y, z, size float32) *Cube {
	cube := &Cube{
		Position: &Vector3{x, y, z},
		Points: []*Vector3{
			{x - size, y - size, z - size},
			{x + size, y - size, z - size},
			{x + size, y + size, z - size},
			{x - size, y + size, z - size},
			{x - size, y - size, z + size},
			{x + size, y - size, z + size},
			{x + size, y + size, z + size},
			{x - size, y + size, z + size},
		},
		Colors: []Color{
			{1, 0, 0, 1},
			{0, 1, 0, 1},
			{0, 0, 1, 1},
			{1, 0, 0, 1},
			{0, 1, 0, 1},
			{0, 0, 1, 1},
			{1, 0, 0, 1},
			{0, 1, 0, 1},
		},
	}
	cube.Indices = []uint16{
		0, 1, 2, 0, 2, 3,
		0, 4, 5, 0, 5, 1,
		1, 5, 6, 1, 6, 2,
		3, 2, 6, 3, 6, 7,
		0, 3, 7, 0, 7, 4,
		4, 7, 6, 4, 6, 5}
	return cube
}

// RotateY rotates the cube around the Y axis
func (c *Cube) RotateY(radian float64) {
	cosine := float32(math.Cos(radian))
	sine := float32(math.Sin(radian))
	for i := len(c.Points) - 1; i > -1; i-- {
		p := c.Points[i]
		x := (p.Z-c.Position.Z)*sine + (p.X-c.Position.X)*cosine
		z := (p.Z-c.Position.Z)*cosine - (p.X-c.Position.X)*sine
		p.X = x + c.Position.X
		p.Z = z + c.Position.Z
	}
}

// RotateX rotates the cube around the Y axis
func (c *Cube) RotateX(radian float64) {
	cosine := float32(math.Cos(radian))
	sine := float32(math.Sin(radian))
	for i := len(c.Points) - 1; i > -1; i-- {
		p := c.Points[i]
		y := (p.Y-c.Position.Y)*cosine - (p.Z-c.Position.Z)*sine
		z := (p.Y-c.Position.Y)*sine + (p.Z-c.Position.Z)*cosine
		p.Y = y + c.Position.Y
		p.Z = z + c.Position.Z
	}
}
