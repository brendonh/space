package space

import (
	"math"

	. "github.com/brendonh/glvec"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"

)

type RenderContext struct {
	QCamRotate Quat

	MPerspective Mat4
	MView Mat4
	VLightDir Vec3
}

func NewRenderContext() *RenderContext {
	context := &RenderContext {
		VLightDir: Vec3 { 0.0, 0.0, -1.0 },
	}

	QIdent(&context.QCamRotate)
	M4Perspective(&context.MPerspective, math.Pi / 4, 800.0 / 600.0, 0.1, 100.0);

	return context
}

func (context *RenderContext) Init() {
	glfw.SwapInterval(1)

	gl.Init()

    gl.ClearColor(0.0, 0.0, 0.0, 1.0);
    gl.ClearDepth(1.0);
    gl.Enable(gl.DEPTH_TEST);
    gl.DepthFunc(gl.LEQUAL);

    gl.Enable(gl.CULL_FACE);
    gl.CullFace(gl.BACK);
}

func (context *RenderContext) StartFrame() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);
	context.TickCamera()
}

func (context *RenderContext) TickCamera() {
	var mCamTransform Mat4
	M4MakeTransform(&mCamTransform, &Vec3{ 0.0, 0.0, -6.0 })
	
	var q Quat
	QRotAng(&q, 0.01, &Vec3 { 1.0, 0.0, 0.0 })

	var camRotate = &context.QCamRotate

	QMul(camRotate, camRotate, &q)
	QRotAng(&q, 0.005, &Vec3 { 0.0, 1.0, 0.0 })
	QMul(camRotate, camRotate, &q)

	var mCamRot Mat4
	QMat4(&mCamRot, camRotate)

	M4MulM4(&context.MView, &mCamTransform, &mCamRot)
}
