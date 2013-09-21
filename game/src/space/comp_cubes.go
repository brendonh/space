package space

import (
	"math"

	"space/render"

	"github.com/go-gl/gl"
	. "github.com/brendonh/glvec"
)

type CubeColor struct {
	R, G, B float32
}

type Cube struct {
	X, Y int
	Color CubeColor
}

type CubesComponent struct {
	Entity *Entity
	Physics *SpacePhysics
	MaterialID render.MaterialID

	verts gl.Buffer
	count int

	mModel Mat4
}

func NewCubesComponent() *CubesComponent {
	comp := &CubesComponent {
		MaterialID: render.GetCubeMaterialID(),
		verts: gl.GenBuffer(),
		count: 0,
	}
	M4MakeScale(&comp.mModel, 0.2)
	return comp
}

func (c *CubesComponent) Init() {
	c.Physics = c.Entity.GetComponent("struct_spacephysics").(*SpacePhysics)
}

func (c *CubesComponent) Tag() string {
	return "cubes"
}

func (c *CubesComponent) SetEntity(e *Entity) {
	c.Entity = e
}

func (c *CubesComponent) Render(context *render.Context, alpha float64) {
	var mModelView Mat4
	var mPhysics = c.Physics.GetModelMatrix(alpha)
	M4MulM4(&mModelView, &mPhysics, &c.mModel)
	M4MulM4(&mModelView, &context.MView, &mModelView)

	context.Enqueue(c.MaterialID, render.CubeRenderArguments{
		MModelView: mModelView,
		Verts: c.verts,
		TriCount: c.count,
	})
}

var CUBE_VERTS = (3+3+3)*(3+3)*6

func (c *CubesComponent) SetCubes(cubes []Cube) {

	var buf = make([]float32, CUBE_VERTS * len(cubes))
	var start = 0

	for _, cube := range cubes {
		addCube(buf[start:start + CUBE_VERTS], cube)
		start += CUBE_VERTS
	}

    c.verts.Bind(gl.ARRAY_BUFFER)
    gl.BufferData(gl.ARRAY_BUFFER, len(buf) * 4, buf, gl.STATIC_DRAW)
	
	c.count = 12 * len(cubes)

}

func addCube(buf []float32, cube Cube) {
    var bufOffset = 0
	var nextSlot = func() []float32 {
		start := bufOffset
		bufOffset += (3+3+3) * (3+3)
		return buf[start:bufOffset]
	}

    // Front
	var q Quat
	QIdent(&q)

	var yAxis = Vec3{ 0.0, 1.0, 0.0 }
	var xAxis = Vec3{ 1.0, 0.0, 0.0 }

	var pos = Vec3{ float32(cube.X), float32(cube.Y), 0 }
	var color = Vec3 { cube.Color.R, cube.Color.G, cube.Color.B }

    addFace(nextSlot(), q, pos, Vec3 {0.0, 0.0, -1.0}, color)

    // Right
	QRotAng(&q, math.Pi / 2, yAxis)
    addFace(nextSlot(), q, pos, Vec3 {1.0, 0.0, 0.0}, color)

    // Back
	QRotAng(&q, math.Pi, yAxis)
    addFace(nextSlot(), q, pos, Vec3 {0.0, 0.0, 1.0}, color)

    // Left
	QRotAng(&q, -math.Pi/2, yAxis)
    addFace(nextSlot(), q, pos, Vec3 {-1.0, 0.0, 0.0}, color)

    // Top
	QRotAng(&q, math.Pi / 2, xAxis)
    addFace(nextSlot(), q, pos, Vec3 {0.0, -1.0, 0.0}, color)

    // Bottom
	QRotAng(&q, -math.Pi / 2, xAxis)
    addFace(nextSlot(), q, pos, Vec3 {0.0, 1.0, 0.0}, color)
}


var vertVecs = []Vec3 {
	Vec3 {  1.0,  1.0,  1.0 },
	Vec3 { -1.0,  1.0,  1.0 },
	Vec3 {  1.0, -1.0,  1.0 },
	Vec3 {  1.0, -1.0,  1.0 },
	Vec3 { -1.0,  1.0,  1.0 },
	Vec3 { -1.0, -1.0,  1.0 },
}

func addFace(buf []float32, rot Quat, pos Vec3, normal Vec3, color Vec3) {
	var bufOffset = 0
	for _, v := range vertVecs {
		var temp Mat3
		QMat3(&temp, rot)
		M3MulV3(&v, &temp, v)
		
		// Vert
		buf[bufOffset + 0] = round(v[0]) + pos[0]
		buf[bufOffset + 1] = round(v[1]) + pos[1]
		buf[bufOffset + 2] = round(v[2]) + pos[2]
		
		// Normal
		buf[bufOffset + 3] = normal[0]
		buf[bufOffset + 4] = normal[1]
		buf[bufOffset + 5] = normal[2]
		
		// Color
		buf[bufOffset + 6] = color[0]
		buf[bufOffset + 7] = color[1]
		buf[bufOffset + 8] = color[2]

		bufOffset += 9
	}
}

func round(f float32) float32 {
	if f < 0 {
		return float32(math.Ceil(float64(f) - 0.5))
	}
	return float32(math.Floor(float64(f) + 0.5))
}