// main package
package main

import (
	"embed"
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"math"
	"regexp"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed cube300.obj
//go:embed cube300r45.obj
//go:embed tri.obj
//go:embed cube.obj
var embedded embed.FS

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	whiteImage = ebiten.NewImage(3, 3)
)

func init() {
	whiteImage.Fill(color.White)
}

func newVertex(x, y, tx, ty float32) ebiten.Vertex {
	return ebiten.Vertex{
		DstX:   x,
		DstY:   y,
		SrcX:   tx,
		SrcY:   ty,
		ColorR: 1,
		ColorG: 0,
		ColorB: 0,
		ColorA: 0.2,
	}
}

// Color is a color
type Color struct {
	R, G, B, A float32
}

// Project 3d coords into 2d
func Project(points []*Vector3) []ebiten.Vertex {
	verts := make([]ebiten.Vertex, len(points))
	var focalLength float32 = 100
	for i := len(points) - 1; i > -1; i-- {
		p := points[i]
		x := p.X*(focalLength/p.Z) + screenWidth*0.5
		y := p.Y*(focalLength/p.Z) + screenHeight*0.5
		verts[i] = newVertex(x, y, 0, 0)
	}

	return verts
}

// CullBackfaces carries out backface culling
func CullBackfaces(indicies []uint16, points []*Vector3) []uint16 {
	renderable := make([]uint16, 0)
	for i := len(indicies) - 1; i > -1; i -= 3 {
		// backface culling
		p1 := points[indicies[i-0]]
		p2 := points[indicies[i-1]]
		p3 := points[indicies[i-2]]
		v1 := p2.Sub(p1)
		v2 := p3.Sub(p1)
		n := v1.CrossProduct(v2)
		if n.DotProduct(p1.Sub(&Vector3{0, 0, 0})) < 0 { // todo put camera position instead of zero vector
			renderable = append(renderable, indicies[i-0], indicies[i-1], indicies[i-2])
		}
	}
	return renderable
}

// Game struct
type Game struct {
	cube  *Cube
	model *Model
}

// Update runs every tick
func (g *Game) Update() error {
	return nil
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawTrianglesOptions{}
	op.Address = ebiten.AddressUnsafe
	// g.model.RotateY(0.01)
	// g.model.RotateX(0.01)
	screen.DrawTriangles(Project(g.model.Vertices), g.model.VertexIndicies, whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image), op)
	// screen.DrawTriangles(Project(g.model.Vertices), CullBackfaces(g.model.VertexIndicies, g.model.Vertices), whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image), op)

	// g.cube.RotateY(0.01)
	// g.cube.RotateX(0.005)
	// screen.DrawTriangles(Project(g.cube.Points), CullBackfaces(g.cube.Indices, g.cube.Points), whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image), op)
}

// Layout returns screen dimensions on resize
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// Model is a 3d model's data
type Model struct {
	Position *Vector3
	Rotation *Vector3

	Vertices []*Vector3
	Normals  []*Vector3
	UVs      []*Vector2

	VertexIndicies []uint16
	NormalIndicies []uint16
	UVIndicies     []uint16
}

func NewModel(position, rotation *Vector3) *Model {
	return &Model{
		Position:       position,
		Rotation:       rotation,
		Vertices:       make([]*Vector3, 0),
		Normals:        make([]*Vector3, 0),
		UVs:            make([]*Vector2, 0),
		VertexIndicies: make([]uint16, 0),
		NormalIndicies: make([]uint16, 0),
		UVIndicies:     make([]uint16, 0),
	}
}

// RotateY rotates the model around the Y axis
func (m *Model) RotateY(radian float64) {
	cosine := float32(math.Cos(radian))
	sine := float32(math.Sin(radian))
	for i := len(m.Vertices) - 1; i > -1; i-- {
		p := m.Vertices[i]
		x := (p.Z-m.Position.Z)*sine + (p.X-m.Position.X)*cosine
		z := (p.Z-m.Position.Z)*cosine - (p.X-m.Position.X)*sine
		p.X = x + m.Position.X
		p.Z = z + m.Position.Z
	}
}

// RotateX rotates the model around the Y axis
func (m *Model) RotateX(radian float64) {
	cosine := float32(math.Cos(radian))
	sine := float32(math.Sin(radian))
	for i := len(m.Vertices) - 1; i > -1; i-- {
		p := m.Vertices[i]
		y := (p.Y-m.Position.Y)*cosine - (p.Z-m.Position.Z)*sine
		z := (p.Y-m.Position.Y)*sine + (p.Z-m.Position.Z)*cosine
		p.Y = y + m.Position.Y
		p.Z = z + m.Position.Z
	}
}

func main() {
	log.SetFlags(log.Lshortfile)

	cube := NewCube(0, 0, 300, 150)
	model := NewModel(NewVector3(0, 0, 300), NewVector3(0, 0, 0))

	if b, err := embedded.ReadFile("tri.obj"); err == nil {
		s := string(b)
		r := regexp.MustCompile(("s.*\n|o.*\n|#.*\n")) // remove smoothing groups, object names and comments
		s = r.ReplaceAllString(s, "")
		f := strings.NewReader(s)
		for {
			var lineType string
			_, err := fmt.Fscanf(f, "%s", &lineType)
			if err != nil {
				if err == io.EOF {
					break
				}
			}

			switch lineType {
			case "v":
				var v1, v2, v3 float32
				_, err := fmt.Fscanf(f, "%f %f %f\n", &v1, &v2, &v3)
				if err != nil {
					log.Fatal(err)
				}
				model.Vertices = append(model.Vertices, NewVector3(v1*1000, v2*1000, v3*1000))
			case "vn":
				var v1, v2, v3 float32
				_, err := fmt.Fscanf(f, "%f %f %f\n", &v1, &v2, &v3)
				if err != nil {
					log.Fatal(err)
				}
				model.Normals = append(model.Normals, NewVector3(v1, v2, v3))
			case "vt":
				var v1, v2 float32
				_, err := fmt.Fscanf(f, "%f %f\n", &v1, &v2)
				if err != nil {
					log.Fatal(err)
				}
				model.UVs = append(model.UVs, NewVector2(v1, v2))
			case "f":
				// f v1/vt1/vn1 v2/vt2/vn2 v3/vt3/vn3
				var v1, v2, v3, vt1, vt2, vt3, vn1, vn2, vn3 uint16
				matches, err := fmt.Fscanf(f, "%d/%d/%d %d/%d/%d %d/%d/%d\n", &v1, &vt1, &vn1, &v2, &vt2, &vn2, &v3, &vt3, &vn3)
				if matches != 9 {
					log.Fatal("Incorrect file format, f must use format: f v1/vt1/vn1 v2/vt2/vn2 v3/vt3/vn3")
				}
				if err != nil {
					log.Fatal(err)
				}
				// model.VertexIndicies = append([]uint16{v1 - 1, v2 - 1, v3 - 1}, model.VertexIndicies...)
				model.VertexIndicies = append(model.VertexIndicies, v1-1, v2-1, v3-1)
				model.NormalIndicies = append(model.NormalIndicies, vn1-1, vn2-1, vn3-1)
				model.UVIndicies = append(model.UVIndicies, vt1-1, vt2-1, vt3-1)
			default:
			}
		}
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Polygons (Ebitengine Demo)")
	if err := ebiten.RunGame(&Game{
		model: model,
		cube:  cube,
	}); err != nil {
		log.Fatal(err)
	}
}
