package space

import (
	"math"

	. "github.com/brendonh/glvec"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"

)

type RenderContext struct {
	QCamRotate Quat
	VCamTranslate Vec3

	MPerspective Mat4
	MView Mat4

	VLightDir Vec3
}

func NewRenderContext() *RenderContext {
	context := &RenderContext {
		VLightDir: Vec3 { 0.0, 0.0, -1.0 },
		VCamTranslate: Vec3 { 0.0, 0.0, 6.0 },
	}

	QIdent(&context.QCamRotate)

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
}

func (context *RenderContext) Resize(width, height int) {
	M4Perspective(&context.MPerspective, math.Pi / 4, 
		float32(width) / float32(height), 0.1, 100.0);

	gl.Viewport(0, 0, width, height)
}

func (context *RenderContext) StartFrame() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);
	context.TickCamera()
}

func (context *RenderContext) TickCamera() {

	// context.VCamTranslate[0] += 0.02
	// context.VCamTranslate[1] += 0.01

	var mCamTransform Mat4
	M4MakeTransform(&mCamTransform, &context.VCamTranslate)

	M4Inverse(&context.MView, &mCamTransform)
	// var q Quat
	// QRotAng(&q, 0.01, &Vec3 { 1.0, 0.0, 0.0 })

	// var camRotate = &context.QCamRotate

	// QMul(camRotate, camRotate, &q)
	// QRotAng(&q, 0.005, &Vec3 { 0.0, 1.0, 0.0 })
	// QMul(camRotate, camRotate, &q)

	// var mCamRot Mat4
	// QMat4(&mCamRot, camRotate)

	// M4MulM4(&context.MView, &mCamTransform, &mCamRot)
}
