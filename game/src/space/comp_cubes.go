package space

import (
	"fmt"
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
	BaseComponent

	Physics *SpacePhysics
	CubeMaterialID render.MaterialID
	GridMaterialID render.MaterialID

	verts []float32
	edges []float32
	cubeCount int

	glVerts gl.Buffer
	glEdges gl.Buffer

	triCount int
	edgeCount int

	mModel Mat4
}

func NewCubesComponent() *CubesComponent {
	comp := &CubesComponent {
		BaseComponent: NewBaseComponent(),
		CubeMaterialID: render.GetCubeMaterialID(),
		GridMaterialID: render.GetGridMaterialID(),
		glVerts: gl.GenBuffer(),
		glEdges: gl.GenBuffer(),
		triCount: 0,
	}
	M4MakeScale(&comp.mModel, 0.2)
	return comp
}

func (c *CubesComponent) Init() {
	c.Physics = c.Entity.GetComponent("struct_spacephysics").(*SpacePhysics)
	globalDispatch.Listen(c, "gl_init", c.OnGLInit)
}

func (c *CubesComponent) Tag() string {
	return "cubes"
}

func (c *CubesComponent) Render(context *render.Context, alpha float64) {
	var mModelView Mat4
	var mPhysics = c.Physics.GetModelMatrix(alpha)
	M4MulM4(&mModelView, &mPhysics, &c.mModel)
	M4MulM4(&mModelView, &context.MView, &mModelView)

	context.Enqueue(c.CubeMaterialID, render.CubeRenderArguments{
		MModelView: mModelView,
		Verts: c.glVerts,
		TriCount: c.triCount,
	})


	context.Enqueue(c.GridMaterialID, render.GridRenderArguments{
		MModelView: mModelView,
		Edges: c.glEdges,
		EdgeCount: c.edgeCount,
	})

}

var CUBE_VERTS = (3+3+3)*(3+3)*6

func (c *CubesComponent) SetCubes(cubes []Cube) {
	c.verts = make([]float32, 0, CUBE_VERTS * len(cubes))
	c.edges = make([]float32, 0, (2 * 4) * len(cubes))

	for _, cube := range cubes {
		c.addCube(cube)
	}

	fmt.Printf("Edges:", len(c.edges))

	c.cubeCount = len(cubes)
	c.RefreshGLBuffer()
}

func(c *CubesComponent) RefreshGLBuffer() {
    c.glVerts.Bind(gl.ARRAY_BUFFER)
    gl.BufferData(gl.ARRAY_BUFFER, len(c.verts) * 4, c.verts, gl.STATIC_DRAW)

	c.glEdges.Bind(gl.ARRAY_BUFFER)
    gl.BufferData(gl.ARRAY_BUFFER, len(c.edges) * 4, c.edges, gl.STATIC_DRAW)

	c.triCount = 12 * c.cubeCount
	c.edgeCount = 4 * c.cubeCount
}

func (c *CubesComponent) OnGLInit(args interface{}) {
	c.RefreshGLBuffer()
}



func (c *CubesComponent) addCube(cube Cube) {

    // Front
	var q Quat
	QIdent(&q)

	var yAxis = Vec3{ 0.0, 1.0, 0.0 }
	var xAxis = Vec3{ 1.0, 0.0, 0.0 }

	var pos = Vec3{ float32(cube.X), float32(cube.Y), 0 }
	var color = Vec3 { cube.Color.R, cube.Color.G, cube.Color.B }

    c.addFace(q, pos, Vec3 {0.0, 0.0, -1.0}, color)

    // Right
	QRotAng(&q, math.Pi / 2, yAxis)
    c.addFace(q, pos, Vec3 {1.0, 0.0, 0.0}, color)

    // Back
	QRotAng(&q, math.Pi, yAxis)
    c.addFace(q, pos, Vec3 {0.0, 0.0, 1.0}, color)

    // Left
	QRotAng(&q, -math.Pi/2, yAxis)
    c.addFace(q, pos, Vec3 {-1.0, 0.0, 0.0}, color)

    // Top
	QRotAng(&q, math.Pi / 2, xAxis)
    c.addFace(q, pos, Vec3 {0.0, -1.0, 0.0}, color)

    // Bottom
	QRotAng(&q, -math.Pi / 2, xAxis)
    c.addFace(q, pos, Vec3 {0.0, 1.0, 0.0}, color)

	// Edges
	c.addEdges(pos)
}



var vertVecs = []Vec3 {
	Vec3 {  1.0,  1.0,  1.0 },
	Vec3 { -1.0,  1.0,  1.0 },
	Vec3 {  1.0, -1.0,  1.0 },
	Vec3 {  1.0, -1.0,  1.0 },
	Vec3 { -1.0,  1.0,  1.0 },
	Vec3 { -1.0, -1.0,  1.0 },
}

func (c *CubesComponent) addFace(rot Quat, pos Vec3, normal Vec3, color Vec3) {
	for _, v := range vertVecs {
		var temp Mat3
		QMat3(&temp, rot)
		M3MulV3(&v, &temp, v)

		c.verts = append(c.verts, 
			round(v[0]) + pos[0], 
			round(v[1]) + pos[1], 
			round(v[2]) + pos[2])

		c.verts = append(c.verts, normal[0], normal[1], normal[2])
		c.verts = append(c.verts, color[0], color[1], color[2])
	}
}

var edgeVecs = []Vec3 {
	Vec3 {  1.0,  1.0, 1.001 },
	Vec3 {  1.0, -1.0, 1.001 },
	Vec3 {  1.0, -1.0, 1.001 },
	Vec3 { -1.0, -1.0, 1.001 },
	Vec3 { -1.0, -1.0, 1.001 },
	Vec3 { -1.0,  1.0, 1.001 },
	Vec3 { -1.0,  1.0, 1.001 },
	Vec3 {  1.0,  1.0, 1.001 },
}

func (c *CubesComponent) addEdges(pos Vec3) {
	for _, v := range edgeVecs {
		V3Add(&v, pos, v)
		c.edges = append(c.edges, v[0], v[1], v[2])
	}
}

func round(f float32) float32 {
	if f < 0 {
		return float32(math.Ceil(float64(f) - 0.5))
	}
	return float32(math.Floor(float64(f) + 0.5))
}