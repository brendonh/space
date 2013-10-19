package space

import (
	"space/render"

	"github.com/go-gl/gl"
	. "github.com/brendonh/glvec"

)

type AvatarRenderer struct {
	BaseComponent
	Position *AvatarPosition

	MModel Mat4
	MModelFrame Mat4

	verts []float32
	colors []float32
	triCount int
	glVerts gl.Buffer
	glColors gl.Buffer
	CubeMaterialID render.MaterialID
}

func NewAvatarRenderer() *AvatarRenderer {
	comp := &AvatarRenderer{
		BaseComponent: NewBaseComponent(),
		glVerts: gl.GenBuffer(),
		glColors: gl.GenBuffer(),
		CubeMaterialID: render.GetCubeMaterialID(),
	}

	M4MakeScale(&comp.MModel, 0.5)

	M4SetTransform(&comp.MModel, Vec3{ 0, 0, 1 })
	

	comp.makeCube()
	return comp
}

func (c *AvatarRenderer) Init() {
	c.Position = c.Entity.GetComponent("struct_avatarposition").(*AvatarPosition)
}

func (c *AvatarRenderer) Tag() string {
	return "avatar"
}

func (c *AvatarRenderer) Render(context *render.Context, alpha float32) {
	var mPosition = c.Position.GetModelMatrix(alpha)

	M4MulM4(&c.MModelFrame, &mPosition, &c.MModel)
	M4MulM4(&c.MModelFrame, &context.MView, &c.MModelFrame)

	context.Enqueue(c.CubeMaterialID, render.CubeRenderArguments{
		MModelView: c.MModelFrame,
		Verts: c.glVerts,
		Colors: c.glColors,
		TriCount: c.triCount,
	})
	
}

func (c *AvatarRenderer) HandleMouse(Ray) bool {
	return false
}


func (c *AvatarRenderer) makeCube() {
	var cube = Cube{ 0, 0, CubeColor{ 1.0, 0.3, 0.3 }, CubeFacesAll() }

	c.verts = addCubeFaces(c.verts, cube, Vec3{ 0, 0, 0 })
	c.colors = addCubeColors(c.colors, cube)

	c.triCount = len(c.verts) / (3 * (3+3))

	c.glVerts.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(c.verts) * 4, c.verts, gl.STATIC_DRAW)

	c.glColors.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(c.colors) * 4, c.colors, gl.STATIC_DRAW)
}