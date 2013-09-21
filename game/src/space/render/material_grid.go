package render

import (
	"github.com/go-gl/gl"
	. "github.com/brendonh/glvec"
)

const (
	GridAttr_VertexPosition = iota
	GridUnif_mPerspective
	GridUnif_mModelView
)

type GridMaterial struct {
	*BaseMaterial
}


func GetGridMaterialID() MaterialID {

	id, ok := GetMaterialID("grid")
	if ok {
		return id
	}

	m := &GridMaterial{
		NewBaseMaterial("grid",
			ShaderSpec{ gl.VERTEX_SHADER, "grid.vert" }, 
			ShaderSpec{ gl.FRAGMENT_SHADER, "grid.frag" },
		),
	}

	m.AttribLocations = []gl.AttribLocation {
		GridAttr_VertexPosition: m.GetAttribLocation("aVertexPosition"),
	}

	m.UniformLocations = []gl.UniformLocation {
		GridUnif_mPerspective: m.GetUniformLocation("uPerspective"),
		GridUnif_mModelView: m.GetUniformLocation("uModelView"),
	}

	m.ID = registerMaterial(m, "grid")

	return m.ID
}


func (m *GridMaterial) Prepare(context *Context) {
	m.Program.Use()
	m.EnableAttribs()

	uPerspective := m.UniformLocations[GridUnif_mPerspective]
	uPerspective.UniformMatrix4fv(false, context.MPerspective)

	gl.Disable(gl.LINE_SMOOTH)
	gl.LineWidth(2)
}


func (m *GridMaterial) Cleanup() {
	m.DisableAttribs()
}


type GridRenderArguments struct {
	MModelView Mat4
	Edges gl.Buffer
	EdgeCount int
}

func (m *GridMaterial) Render(args interface{}) {
	gridArgs := args.(GridRenderArguments)
	
	gridArgs.Edges.Bind(gl.ARRAY_BUFFER)
	
	aVertexPosition := m.AttribLocations[GridAttr_VertexPosition]
	
	aVertexPosition.AttribPointer(3, gl.FLOAT, false, 0, uintptr(0))
	
	uModelView := m.UniformLocations[GridUnif_mModelView]
	uModelView.UniformMatrix4fv(false, gridArgs.MModelView)

	gl.DrawArrays(gl.LINES, 0, gridArgs.EdgeCount * 2)
}


