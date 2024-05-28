package main

// Vector2 struct
type Vector2 struct {
	X, Y float32
}

// NewVector2 returns a Vector3
func NewVector2(x, y float32) *Vector2 {
	return &Vector2{x, y}
}

// Vector3 struct
type Vector3 struct {
	X, Y, Z float32
}

// NewVector3 returns a Vector3
func NewVector3(x, y, z float32) *Vector3 {
	return &Vector3{x, y, z}
}

// DotProduct returns the dot product of 2 Vector3s
func (v *Vector3) DotProduct(o *Vector3) float32 {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}

// CrossProduct returns the cross product of 2 Vector3s
func (v *Vector3) CrossProduct(o *Vector3) *Vector3 {
	return &Vector3{
		X: v.Y*o.Z - v.Z*o.Y,
		Y: v.Z*o.X - v.X*o.Z,
		Z: v.X*o.Y - v.Y*o.X,
	}
}

// Sub returns v-o
func (v *Vector3) Sub(o *Vector3) *Vector3 {
	return &Vector3{X: v.X - o.X, Y: v.Y - o.Y, Z: v.Z - o.Z}
}
