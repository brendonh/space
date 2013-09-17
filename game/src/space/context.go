package space

import (
	"math"

	. "github.com/brendonh/glvec"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"

)

type RenderContext struct {
	VCamTranslate Vec3
	VCamPos Vec3

	MPerspective Mat4
	MView Mat4

	VLightDir Vec3

	FollowPhysics *SpacePhysics
	FrameDelta float64
}

func NewRenderContext() *RenderContext {
	context := &RenderContext {
		VLightDir: Vec3 { 0.0, 0.0, -1.0 },
		VCamTranslate: Vec3 { 0.0, -2.0, 6.0 },
	}

	return context
}

func (context *RenderContext) Init() {
	glfw.SwapInterval(1)

	gl.Init()

    gl.ClearColor(0.0, 0.0, 0.0, 1.0)
    gl.ClearDepth(1.0)
    gl.Enable(gl.DEPTH_TEST)
    gl.DepthFunc(gl.LEQUAL)

    gl.Enable(gl.CULL_FACE)
    gl.CullFace(gl.BACK)

	gl.Enable(gl.PROGRAM_POINT_SIZE)

	gl.Enable( gl.BLEND )
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.Enable(gl.MULTISAMPLE)
}

func (context *RenderContext) FollowEntity(e *Entity) {
	context.FollowPhysics = e.GetComponent("struct_spacephysics").(*SpacePhysics)
}

func (context *RenderContext) Resize(width, height int) {
	M4Perspective(&context.MPerspective, math.Pi / 4, 
		float32(width) / float32(height), 0.1, 100.0);

	gl.Viewport(0, 0, width, height)
}

func (context *RenderContext) StartFrame(delta float64) {
	context.FrameDelta = delta
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);
	context.SetCamera()
}

func (context *RenderContext) SetCamera() {
	phys := context.FollowPhysics

	var center, up Vec3

	var pos = phys.InterpolatePosition(context.FrameDelta)
	
	center = Vec3 { float32(pos.PosX), float32(pos.PosY), 0.0 }
	V3Add(&context.VCamPos, center, context.VCamTranslate)
	up = Vec3 { 0.0, 1.0, 1.0 }

	M4LookAt(&context.MView, context.VCamPos, center, up)
}
