package render

import (
	"github.com/go-gl/gl"
	. "github.com/brendonh/glvec"
)

const (
	GridAttr_VertexPosition = iota
	GridUnif_mPerspective
	GridUnif_mModelView
	GridUnif_bActive
	GridUnif_vActiveCoords
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
		GridUnif_bActive: m.GetUniformLocation("uActive"),
		GridUnif_vActiveCoords: m.GetUniformLocation("uActiveCoords"),
	}

	m.ID = registerMaterial(m, "grid")

	return m.ID
}


func (m *GridMaterial) Prepare(context *Context) {
	m.Program.Use()
	m.EnableAttribs()

	uPerspective := m.UniformLocations[GridUnif_mPerspective]
	uPerspective.UniformMatrix4fv(false, context.MPerspective)

	gl.Enable(gl.LINE_SMOOTH)
	gl.LineWidth(3)
}


func (m *GridMaterial) Cleanup() {
	m.DisableAttribs()
}


type GridRenderArguments struct {
	MModelView Mat4
	Edges gl.Buffer
	EdgeCount int
	Active []int
}

func (m *GridMaterial) Render(args interface{}) {
	gridArgs := args.(GridRenderArguments)
	
	gridArgs.Edges.Bind(gl.ARRAY_BUFFER)
	
	aVertexPosition := m.AttribLocations[GridAttr_VertexPosition]
	
	aVertexPosition.AttribPointer(3, gl.FLOAT, false, 0, uintptr(0))
	
	uModelView := m.UniformLocations[GridUnif_mModelView]
	uModelView.UniformMatrix4fv(false, gridArgs.MModelView)

	uActive := m.UniformLocations[GridUnif_bActive]
	uActiveCoords := m.UniformLocations[GridUnif_vActiveCoords]

	var active = gridArgs.Active
	if active == nil {
		uActive.Uniform1i(0)
	} else {
		uActive.Uniform1i(1)
		uActiveCoords.Uniform2i(active[0], active[1])
	}

	gl.DrawArrays(gl.LINES, 0, gridArgs.EdgeCount * 6)
}


