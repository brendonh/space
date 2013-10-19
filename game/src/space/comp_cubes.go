package space

import (
	"math"

	"space/render"

	"github.com/go-gl/gl"
	. "github.com/brendonh/glvec"
)

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
	colors []float32
	edges []float32

	glVerts gl.Buffer
	glColors gl.Buffer
	glEdges gl.Buffer

	triCount int
	edgeCount int
}

func NewCubesComponent() *CubesComponent {
	comp := &CubesComponent{
		BaseComponent: NewBaseComponent(),
		CubeMaterialID: render.GetCubeMaterialID(),
		GridMaterialID: render.GetGridMaterialID(),
		ShowEdges: true,
		glVerts: gl.GenBuffer(),
		glColors: gl.GenBuffer(),
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
	case "update_colors":
		c.SetColorOverrides(args.([]CubeColorOverride))
	}
}

func (c *CubesComponent) Render(context *render.Context, alpha float32) {
	var mPhysics = c.Physics.GetModelMatrix(alpha)

	M4MulM4(&c.MModelFrame, &mPhysics, &c.MModel)

	var mModelView Mat4
	M4MulM4(&mModelView, &context.MView, &c.MModelFrame)

	if c.triCount > 0 {
		context.Enqueue(c.CubeMaterialID, render.CubeRenderArguments{
			MModelView: mModelView,
			Verts: c.glVerts,
			Colors: c.glColors,
			TriCount: c.triCount,
		})
	}

	if c.ShowEdges && c.edgeCount > 0 {
		var active []float32
		var tile = c.Rooms.SelectedTile
		if tile != nil {
			shipPos := tile.GetShipPos()
			active = []float32 { 
				(float32(shipPos.X) + c.cubes.Center[0]) * 2,
				(float32(shipPos.Y) + c.cubes.Center[1]) * 2, 
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

func (c *CubesComponent) HandleMouse(ray Ray, action MouseAction) bool {
	worldPos, ok := ray.PlaneIntersect(Plane{ 
		Point: Vec3{ 0.0, 0.0, 1.0 },
		Normal: Vec3{ 0.0, 0.0, 1.0 },
	})

	if ok {
		var modelInv Mat4
		M4Inverse(&modelInv, &c.MModelFrame)
		M4MulV3(&worldPos, &modelInv, worldPos)
		
		V3Add(&worldPos, worldPos, Vec3{ 1, 1, 1 })

		shipPos := c.ModelToTile(worldPos)
		if !c.Rooms.SetSelectedTile(shipPos) {
			return false
		}

		if action == MOUSE_PRESS {
			c.Rooms.TriggerSelectedTile()
		}

		return true		
	}

	c.Rooms.ClearSelectedTile()
	return false
}


func (c *CubesComponent) ModelToTile(worldPos Vec3) Vec2i {
	V3ScalarMul(&worldPos, worldPos, 0.5)
	
	V3Sub(&worldPos, worldPos, c.cubes.Center)
	
	return Vec2i{
		int(math.Floor(float64(worldPos[0]))),
		int(math.Floor(float64(worldPos[1]))),
	}
}


func(c *CubesComponent) RefreshGLBuffer() {
	c.triCount = len(c.verts) / (3 * (3 + 3))
	c.edgeCount = len(c.edges) / 3

	if c.triCount > 0 {
		c.glVerts.Bind(gl.ARRAY_BUFFER)
		gl.BufferData(gl.ARRAY_BUFFER, len(c.verts) * 4, c.verts, gl.STATIC_DRAW)
	}

	if c.edgeCount > 0 {
		c.glEdges.Bind(gl.ARRAY_BUFFER)
		gl.BufferData(gl.ARRAY_BUFFER, len(c.edges) * 4, c.edges, gl.STATIC_DRAW)
	}

	c.RefreshColorBuffer()
}

func (c *CubesComponent) RefreshColorBuffer() {
	if c.triCount > 0 {
		c.glColors.Bind(gl.ARRAY_BUFFER)
		gl.BufferData(gl.ARRAY_BUFFER, len(c.colors) * 4, c.colors, gl.DYNAMIC_DRAW)
	}
}

func (c *CubesComponent) OnGLInit(args interface{}) {
	c.RefreshGLBuffer()
}


type CubeColorOverride struct {
	X, Y int
	Color CubeColor
}

func (c *CubesComponent) SetColorOverrides(overrides []CubeColorOverride) {
	c.colors = c.colors[:0]
	for _, cube := range c.cubes.Cubes {
		var addColor *CubeColor
		for _, override := range overrides {
			if cube.Pos.X == override.X && cube.Pos.Y == override.Y {
				addColor = &override.Color
				break
			}
		}
		c.colors = addCubeColors(c.colors, cube, addColor)		
	}
	c.RefreshColorBuffer()
}


func (c *CubesComponent) setCubes(cubeSet *CubeSet) {
	c.cubes = cubeSet
	c.checkSides(cubeSet.Cubes)

	c.verts = c.verts[:0]
	c.colors = c.colors[:0]
	c.edges = c.edges[:0]

	for _, cube := range cubeSet.Cubes {
		c.verts = addCubeFaces(c.verts, cube, cubeSet.Center)
		c.colors = addCubeColors(c.colors, cube, nil)
		c.edges = addCubeEdges(c.edges, cube, cubeSet.Center)
	}

	c.RefreshGLBuffer()
}

func (c *CubesComponent) checkSides(cubes []Cube) {
	for i := range cubes {
		cube := &cubes[i]
		x, y := cube.Pos.X, cube.Pos.Y
		cube.Faces = CubeFacesAll()
		for _, prev := range cubes {
			px, py := prev.Pos.X, prev.Pos.Y
			if px == x && py == y + 1 {
				cube.Faces.Unset(CUBE_TOP)
			} else if px == x - 1 && py == y {
				cube.Faces.Unset(CUBE_LEFT)
			} else if px == x && py == y - 1 {
				cube.Faces.Unset(CUBE_BOTTOM)
			} else if px == x + 1 && py == y {
				cube.Faces.Unset(CUBE_RIGHT)			
			}
		}
	}
}



