package render

import (
	"fmt"
	"math"
	"sort"

	. "github.com/brendonh/glvec"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"

)

type WindowConfig struct {
	Title string
	Width int
	Height int
	Fullscreen bool
}

var defaultWindowConfig = WindowConfig{
	Title: "SPACE",
	Width: 800,
	Height: 600,
	Fullscreen: false,
}

type Context struct {
	Config *WindowConfig
	Window *glfw.Window
	RenderQueue RenderQueue

	MPerspective Mat4
	MView Mat4

	MCamRotate Mat3
	MPerspectiveInverse Mat4

	VCamPos Vec3
	VLightDir Vec3
}


func NewContext() *Context {
	var light = Vec3 { 0.0, -1.0, -2.0 }
	V3Normalize(&light, light)

	context := &Context {
		Config: &defaultWindowConfig,
		VLightDir: light,
	}
	return context
}

func (context *Context) Init() {
	glfw.SetErrorCallback(func (err glfw.ErrorCode, desc string) {
		fmt.Printf("%v: %v\n", err, desc)
	})


	if !glfw.Init() {
		panic("Can't init glfw!")
	}

	context.initWindow()
	gl.Init()
	context.initGL()
}

func (context *Context) ToggleFullscreen() {
	_, err := context.Window.GetMonitor()
	isFullscreen := err == nil

	if isFullscreen {
		context.Config = &defaultWindowConfig
	} else {
		context.Config = &WindowConfig{ Fullscreen: true }
	}
	context.initWindow()
	context.initGL()
}

func (context *Context) initWindow() {
	glfw.WindowHint(glfw.Samples, 1)

	var monitor *glfw.Monitor
	var width, height int
	var config = context.Config

	if config.Fullscreen {
		var err error
		monitor, err = glfw.GetPrimaryMonitor()
		if err != nil {
			panic(err)
		}
		
		mode, err := monitor.GetVideoMode()
		width = mode.Width
		height = mode.Height
	} else {
		width = config.Width
		height = config.Height
	}

	window, err := glfw.CreateWindow(width, height, config.Title, monitor, context.Window)
	if err != nil {
		panic(err)
	}

	if context.Window != nil {
		context.Window.Destroy()
	}

	context.Window = window

	window.MakeContextCurrent()
}

func (context *Context) initGL() {
	glfw.SwapInterval(0)

    gl.ClearColor(0.0, 0.0, 0.0, 1.0)
    gl.ClearDepth(1.0)
    gl.Enable(gl.DEPTH_TEST)
    gl.DepthFunc(gl.LEQUAL)

    gl.Enable(gl.CULL_FACE)
    gl.CullFace(gl.BACK)

	gl.Enable(gl.PROGRAM_POINT_SIZE)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.Enable(gl.MULTISAMPLE)

	width, height := context.Window.GetSize()
	context.Resize(width, height)
}


func (context *Context) Resize(width, height int) {
	M4Perspective(&context.MPerspective, math.Pi / 4, 
		float32(width) / float32(height), 1.0, 100.0);
	
	M4Inverse(&context.MPerspectiveInverse, &context.MPerspective)

	gl.Viewport(0, 0, width, height)
}

func (context *Context) StartFrame() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);
}

func (context *Context) Enqueue(materialID MaterialID, args interface{}) {
	context.RenderQueue = append(context.RenderQueue, RenderQueueEntry {
		MaterialID: materialID,
		Args: args,
	})
}

func (context *Context) FlushQueue() {
	var matID MaterialID = -1
	var mat Material = nil

	sort.Sort(context.RenderQueue)

	for _, entry := range context.RenderQueue {

		if entry.MaterialID != matID {
			if matID != -1 {
				mat.Cleanup()
			}

			mat = materials[entry.MaterialID]
			mat.Prepare(context)

			matID = entry.MaterialID
		}

		mat.Render(entry.Args)
	}

	if matID != -1 {
		mat.Cleanup()
	}

	context.RenderQueue = context.RenderQueue[:0]

	context.Window.SwapBuffers()
}


// --------------------------------------

type RenderQueueEntry struct {
	MaterialID MaterialID
	Args interface{}
}

type RenderQueue []RenderQueueEntry

func (q RenderQueue) Len() int {
	return len(q)
}

func (q RenderQueue) Less(i, j int) bool {
	return q[i].MaterialID < q[j].MaterialID
}

func (q RenderQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

