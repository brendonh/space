package main

import (
	"fmt"
	"math"
	"os"

	"space"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"

	. "github.com/brendonh/glvec"
)


const (
	Title  = "SPACE"
	Width  = 800
	Height = 600
)


func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}


func main() {

	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("Can't init glfw!")
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(Width, Height, Title, nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	glfw.SwapInterval(1)

	gl.Init()

	if err := initScene(); err != nil {
		fmt.Fprintf(os.Stderr, "init: %s\n", err)
		return
	}
	defer destroyScene()

	for !window.ShouldClose() {
		drawScene()
		window.SwapBuffers()
		glfw.PollEvents()
	}

}


var camRotate Quat
var camDegree float32 = 0

var mPerspective Mat4
var vLightDir = Vec3 { 0.0, 0.0, -1.0 }

var cube *space.Cubes

func initScene() (err error) {
    gl.ClearColor(0.0, 0.0, 0.0, 1.0);
    gl.ClearDepth(1.0);
    gl.Enable(gl.DEPTH_TEST);
    gl.DepthFunc(gl.LEQUAL);

    gl.Enable(gl.CULL_FACE);
    gl.CullFace(gl.BACK);

	cube = space.NewCubes()

	QIdent(&camRotate)
	M4Perspective(&mPerspective, math.Pi / 4, 800.0 / 600.0, 0.1, 100.0);

	return nil
}


func drawScene() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);

	var mCamTransform Mat4
	M4MakeTransform(&mCamTransform, &Vec3{ 0.0, 0.0, -6.0 })
	
	var q Quat
	QRotAng(&q, 0.01, &Vec3 { 1.0, 0.0, 0.0 })
	QMul(&camRotate, &camRotate, &q)
	QRotAng(&q, 0.005, &Vec3 { 0.0, 1.0, 0.0 })
	QMul(&camRotate, &camRotate, &q)

	var mCamRot Mat4
	QMat4(&mCamRot, &camRotate)

	var mView Mat4
	M4MulM4(&mView, &mCamTransform, &mCamRot)

	cube.Render(&mPerspective, &mView, vLightDir)
}



func destroyScene() {

}
