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

type CubeSet struct {
	Cubes []Cube
	Center Vec3
}

type CubesComponent struct {
	BaseComponent

	Physics *SpacePhysics
	Rooms *RoomsComponent

	CubeMaterialID render.MaterialID
	GridMaterialID render.MaterialID

	MModel Mat4
	MModelFrame Mat4

	ShowEdges bool

	cubes *CubeSet

	verts []float32
	edges []float32

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
	c.Rooms = c.Entity.GetComponent("rooms").(*RoomsComponent)
	globalDispatch.Listen(c, "gl_init", c.OnGLInit)
}

func (c *CubesComponent) Tag() string {
	return "cubes"
}

func (c *CubesComponent) Event(tag string, args interface{}) {
	switch(tag) {
	case "update_cubes":
		c.setCubes(args.(*CubeSet))
	}
}

func (c *CubesComponent) Render(context *render.Context, alpha float64) {
	var mPhysics = c.Physics.GetModelMatrix(alpha)

	M4MulM4(&c.MModelFrame, &mPhysics, &c.MModel)

	var mModelView Mat4
	M4MulM4(&mModelView, &context.MView, &c.MModelFrame)

	if c.triCount > 0 {
		context.Enqueue(c.CubeMaterialID, render.CubeRenderArguments{
			MModelView: mModelView,
			Verts: c.glVerts,
			TriCount: c.triCount,
		})
	}

	if c.ShowEdges && c.edgeCount > 0 {
		var active []int
		var tile = c.Rooms.SelectedTile
		if tile != nil {
			active = []int { 
				int(float32(tile.X) - c.cubes.Center[0]) * 2, 
				int(float32(tile.Y) - c.cubes.Center[1]) * 2, 
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

func (c *CubesComponent) HandleMouse(ray Ray) bool {
	worldPos, ok := ray.PlaneIntersect(Plane{ 
		Point: Vec3{ 0.0, 0.0, 1.0 },
		Normal: Vec3{ 0.0, 0.0, 1.0 },
	})

	if ok {
		var modelInv Mat4
		M4Inverse(&modelInv, &c.MModelFrame)
		M4MulV3(&worldPos, &modelInv, worldPos)
		
		// Adjust for cube vert offset and scale
		V3Add(&worldPos, worldPos, Vec3{ 1, 1, 1 })
		V3ScalarMul(&worldPos, worldPos, 0.5)
		
		V3Add(&worldPos, worldPos, c.cubes.Center)
		
		// TODO: Something other than this	
		x := int(math.Floor(float64(worldPos[0]))) 
		y := int(math.Floor(float64(worldPos[1])))
		for _, cube := range c.cubes.Cubes {
			if cube.X == x && cube.Y == y {
				c.Rooms.SetSelectedTile(x, y)
				return true
			}
		}
	}

	c.Rooms.ClearSelectedTile()
	return false
}

var CUBE_VERTS = (3+3+3)*(3+3)*6


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


func (c *CubesComponent) setCubes(cubeSet *CubeSet) {
	c.cubes = cubeSet
	c.checkSides(cubeSet.Cubes)

	c.verts = c.edges[:0]
	c.edges = c.edges[:0]

	for _, cube := range cubeSet.Cubes {
		c.addCube(cube, cubeSet.Center)
		c.addEdges(cube, cubeSet.Center)
	}

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


func (c *CubesComponent) addCube(cube Cube, center Vec3) {
	var yAxis = Vec3{ 0.0, 1.0, 0.0 }
	var xAxis = Vec3{ 1.0, 0.0, 0.0 }

	var pos = Vec3{ 
		(float32(cube.X) - center[0]) * 2.0,
		(float32(cube.Y) - center[1]) * 2.0, 
		0,
	}

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

func (c *CubesComponent) addEdges(cube Cube, center Vec3) {
	var pos = Vec3{ 
		(float32(cube.X) - center[0]) * 2.0,
		(float32(cube.Y) - center[1]) * 2.0,
		0,
	}

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