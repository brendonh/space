package render

import (
	"github.com/go-gl/gl"
	. "github.com/brendonh/glvec"
)

const (
	CubeAttr_VertexPosition = iota
	CubeAttr_VertexNormal
	CubeAttr_VertexColor

	CubeUnif_vLightDirection
	CubeUnif_mPerspective

	CubeUnif_mModelView
	CubeUnif_mNormal
)

type CubeMaterial struct {
	*BaseMaterial
}


func GetCubeMaterialID() MaterialID {

	id, ok := GetMaterialID("cubes")
	if ok {
		return id
	}

	m := &CubeMaterial{
		NewBaseMaterial("cube",
			ShaderSpec{ gl.VERTEX_SHADER, "cube.vert" }, 
			ShaderSpec{ gl.FRAGMENT_SHADER, "cube.frag" },
		),
	}

	m.AttribLocations = []gl.AttribLocation {
		CubeAttr_VertexPosition: m.GetAttribLocation("aVertexPosition"),
		CubeAttr_VertexNormal: m.GetAttribLocation("aVertexNormal"),
		CubeAttr_VertexColor: m.GetAttribLocation("aVertexColor"),
	}

	m.UniformLocations = []gl.UniformLocation {
		CubeUnif_vLightDirection: m.GetUniformLocation("uLightDirection"),
		CubeUnif_mPerspective: m.GetUniformLocation("uPerspective"),
		CubeUnif_mModelView: m.GetUniformLocation("uModelView"),
		CubeUnif_mNormal: m.GetUniformLocation("uNormalMatrix"),
	}

	m.ID = registerMaterial(m, "cubes")

	return m.ID
}


func (cm *CubeMaterial) Prepare(context *Context) {
	cm.Program.Use()
	cm.EnableAttribs()

	uPerspective := cm.UniformLocations[CubeUnif_mPerspective]
	uPerspective.UniformMatrix4fv(false, context.MPerspective)

	uLightDirection := cm.UniformLocations[CubeUnif_vLightDirection]
	vLight := &context.VLightDir
	uLightDirection.Uniform3f(vLight[0],vLight[1],vLight[2])
}


func (cm *CubeMaterial) Cleanup() {
	cm.DisableAttribs()
}


type CubeRenderArguments struct {
	MModelView Mat4
	Verts gl.Buffer
	TriCount int
}

var CM_OFFSET_VERTICES = uintptr(0)
var CM_OFFSET_NORMALS  = uintptr(3 * 4)
var CM_OFFSET_COLORS   = uintptr(6 * 4)


func (cm *CubeMaterial) Render(args interface{}) {
	cubeArgs := args.(CubeRenderArguments)

	cubeArgs.Verts.Bind(gl.ARRAY_BUFFER)
	
	aVertexPosition := cm.AttribLocations[CubeAttr_VertexPosition]
	aVertexNormal := cm.AttribLocations[CubeAttr_VertexNormal]
	aVertexColor := cm.AttribLocations[CubeAttr_VertexColor]
	
	aVertexPosition.AttribPointer(3, gl.FLOAT, false, 9 * 4, CM_OFFSET_VERTICES)
	aVertexNormal.AttribPointer(3, gl.FLOAT, false, 9 * 4, CM_OFFSET_NORMALS)
	aVertexColor.AttribPointer(3, gl.FLOAT, false, 9 * 4, CM_OFFSET_COLORS)
	
	uModelView := cm.UniformLocations[CubeUnif_mModelView]
	uModelView.UniformMatrix4fv(false, cubeArgs.MModelView)
	
	var mMVN Mat3
	M4RotationMatrix(&mMVN, &cubeArgs.MModelView)
	
	uNormal := cm.UniformLocations[CubeUnif_mNormal]
	uNormal.UniformMatrix3fv(false, mMVN)
	
	gl.DrawArrays(gl.TRIANGLES, 0, cubeArgs.TriCount * 3)
}


