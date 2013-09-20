package render

import (
	"math"

	. "github.com/brendonh/glvec"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"

)

type Context struct {
	VCamPos Vec3

	MPerspective Mat4
	MView Mat4

	VLightDir Vec3
}


func NewContext() *Context {
	context := &Context {
		VLightDir: Vec3 { 0.0, 0.0, -1.0 },
	}

	return context
}

func (context *Context) Init() {
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


func (context *Context) Resize(width, height int) {
	M4Perspective(&context.MPerspective, math.Pi / 4, 
		float32(width) / float32(height), 0.1, 100.0);

	gl.Viewport(0, 0, width, height)
}

func (context *Context) StartFrame() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);
}