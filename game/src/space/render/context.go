package render

import (
	"math"
	"sort"

	. "github.com/brendonh/glvec"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"

)

type Context struct {
	VCamPos Vec3

	MPerspective Mat4
	MView Mat4

	VLightDir Vec3

	RenderQueue RenderQueue
}


func NewContext() *Context {
	context := &Context {
		VLightDir: Vec3 { 0.0, -1.0, -2.0 },
	}

	return context
}

func (context *Context) Init() {
	glfw.SwapInterval(0)

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

