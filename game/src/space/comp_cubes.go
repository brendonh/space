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
	Top, Left, Bottom, Right bool
	EdgeTop, EdgeRight bool
}

type CubesComponent struct {
	BaseComponent

	Physics *SpacePhysics
	Control *ShipControl

	CubeMaterialID render.MaterialID
	GridMaterialID render.MaterialID

	MModel Mat4

	ShowEdges bool

	verts []float32
	edges []float32
	cubeCount int

	glVerts gl.Buffer
	glEdges gl.Buffer

	triCount int
	edgeCount int


}

func NewCubesComponent() *CubesComponent {
	comp := &CubesComponent {
		BaseComponent: NewBaseComponent(),
		CubeMaterialID: render.GetCubeMaterialID(),
		GridMaterialID: render.GetGridMaterialID(),
		ShowEdges: true,
		glVerts: gl.GenBuffer(),
		glEdges: gl.GenBuffer(),
		triCount: 0,
	}
	M4Ident(&comp.MModel)
	return comp
}

func (c *CubesComponent) Init() {
	c.Physics = c.Entity.GetComponent("struct_spacephysics").(*SpacePhysics)
	var control = c.Entity.FindComponent("struct_shipcontrol")
	if control != nil {
		c.Control = control.(*ShipControl)
	}
	globalDispatch.Listen(c, "gl_init", c.OnGLInit)
}

func (c *CubesComponent) Tag() string {
	return "cubes"
}

func (c *CubesComponent) Render(context *render.Context, alpha float64) {
	var mModelView Mat4
	var mPhysics = c.Physics.GetModelMatrix(alpha)
	M4MulM4(&mModelView, &mPhysics, &c.MModel)
	M4MulM4(&mModelView, &context.MView, &mModelView)

	if c.triCount > 0 {
		context.Enqueue(c.CubeMaterialID, render.CubeRenderArguments{
			MModelView: mModelView,
			Verts: c.glVerts,
			TriCount: c.triCount,
		})
	}

	if c.ShowEdges && c.edgeCount > 0 {
		var active []int
		if (c.Control != nil) {
			active = []int { 
				c.Control.ActiveTile[0] * 2,
				c.Control.ActiveTile[1] * 2,
			}
		}

		context.Enqueue(c.GridMaterialID, render.GridRenderArguments{
			MModelView: mModelView,
			Edges: c.glEdges,
			EdgeCount: c.edgeCount,
			Active: active,
		})
	}

}

var CUBE_VERTS = (3+3+3)*(3+3)*6

func (c *CubesComponent) SetCubes(cubes []Cube) {
	c.checkSides(cubes)

	c.verts = c.edges[:0]
	c.edges = c.edges[:0]

	for _, cube := range cubes {
		c.addCube(cube)
		c.addEdges(cube)
	}

	c.cubeCount = len(cubes)
	c.RefreshGLBuffer()
}

func(c *CubesComponent) RefreshGLBuffer() {
	c.triCount = len(c.verts) / ((3+3+3) + (3+3))
	c.edgeCount = len(c.edges) / 3

	if c.triCount > 0 {
		c.glVerts.Bind(gl.ARRAY_BUFFER)
		gl.BufferData(gl.ARRAY_BUFFER, len(c.verts) * 4, c.verts, gl.STATIC_DRAW)
	}

	if c.edgeCount > 0 {
		c.glEdges.Bind(gl.ARRAY_BUFFER)
		gl.BufferData(gl.ARRAY_BUFFER, len(c.edges) * 4, c.edges, gl.STATIC_DRAW)
	}
}

func (c *CubesComponent) OnGLInit(args interface{}) {
	c.RefreshGLBuffer()
}


func (c *CubesComponent) checkSides(cubes []Cube) {
	for i := range cubes {
		cube := &cubes[i]
		x, y := cube.X, cube.Y
		cube.Top = true
		cube.Left = true
		cube.Bottom = true
		cube.Right = true
		for _, prev := range cubes {
			px, py := prev.X, prev.Y
			if px == x && py == y + 1 {
				cube.Top = false
			} else if px == x - 1 && py == y {
				cube.Left = false
			} else if px == x && py == y - 1 {
				cube.Bottom = false
			} else if px == x + 1 && py == y {
				cube.Right = false
			}
		}
	}

}


func (c *CubesComponent) addCube(cube Cube) {
	var yAxis = Vec3{ 0.0, 1.0, 0.0 }
	var xAxis = Vec3{ 1.0, 0.0, 0.0 }

	var pos = Vec3{ float32(cube.X * 2), float32(cube.Y * 2), 0 }
	var color = Vec3{ cube.Color.R, cube.Color.G, cube.Color.B }

	var q Quat

	// Front
	QIdent(&q)
    c.addFace(q, pos, Vec3 {0.0, 0.0, -1.0}, color)

    // Back
	QRotAng(&q, math.Pi, yAxis)
    c.addFace(q, pos, Vec3 {0.0, 0.0, 1.0}, color)

	if cube.Top {
		QRotAng(&q, math.Pi / 2, xAxis)
		c.addFace(q, pos, Vec3 {0.0, -1.0, 0.0}, color)
	}

	if cube.Left {
		QRotAng(&q, math.Pi / 2, yAxis)
		c.addFace(q, pos, Vec3 {1.0, 0.0, 0.0}, color)
	}

	if cube.Bottom {
		QRotAng(&q, -math.Pi / 2, xAxis)
		c.addFace(q, pos, Vec3 {0.0, 1.0, 0.0}, color)
	}

	if cube.Right {
		QRotAng(&q, -math.Pi/2, yAxis)
		c.addFace(q, pos, Vec3 {-1.0, 0.0, 0.0}, color)
	}
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
}

func (c *CubesComponent) addEdges(cube Cube) {
	var pos = Vec3{ float32(cube.X * 2), float32(cube.Y * 2), 0 }

	var addEdge = func(verts... Vec3) {
		for _, v := range verts {
			var nv Vec3
			V3Add(&nv, pos, v)
			c.edges = append(c.edges, nv[0], nv[1], nv[2])
		}
	}

	// Top
	if cube.Top {
		addEdge(Vec3 { -1.0,  1.0, 1.001 }, Vec3 {  1.0,  1.0, 1.001 })
	}
	
	// Left
	addEdge(Vec3 {  -1.0,  1.0, 1.001 }, Vec3 {  -1.0, -1.0, 1.001 })

	// Bottom
	addEdge(Vec3 { -1.0, -1.0, 1.001 }, Vec3 { 1.0, -1.0, 1.001 })

	// Right
	if cube.Right {
		addEdge(Vec3 { 1.0,  1.0, 1.001 }, Vec3 { 1.0, -1.0, 1.001 })
	}

}

func round(f float32) float32 {
	if f < 0 {
		return float32(math.Ceil(float64(f) - 0.5))
	}
	return float32(math.Floor(float64(f) + 0.5))
}