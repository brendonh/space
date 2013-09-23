package render

import (
	"github.com/go-gl/gl"
	. "github.com/brendonh/glvec"
)

const (
	DustUnif_mPerspective = iota
	DustUnif_mView

	DustUnif_vBasePosition
	DustUnif_vCenterPosition
	DustUnif_fDustBoxSize
)

type DustMaterial struct {
	*BaseMaterial
	starBuffer gl.Buffer
}

func GetDustMaterialID() MaterialID {

	id, ok := GetMaterialID("dust")
	if ok {
		return id
	}

	m := &DustMaterial {
		NewBaseMaterial("starfield", 
			ShaderSpec{ gl.VERTEX_SHADER, "dust.vert" }, 
			ShaderSpec{ gl.FRAGMENT_SHADER, "dust.frag" },
		),
		gl.GenBuffer(),
	}

	m.UniformLocations = []gl.UniformLocation {
		DustUnif_mPerspective: m.GetUniformLocation("uPerspective"),
		DustUnif_mView: m.GetUniformLocation("uView"),

		DustUnif_vBasePosition: m.GetUniformLocation("uBasePosition"),
		DustUnif_vCenterPosition: m.GetUniformLocation("uCenterPosition"),
		DustUnif_fDustBoxSize: m.GetUniformLocation("uDustBoxSize"),
	}

	m.ID = registerMaterial(m, "dust")

	return m.ID
}


func (dm *DustMaterial) Prepare(context *Context) {
	dm.Program.Use()
	dm.EnableAttribs()

    uPerspective := dm.UniformLocations[DustUnif_mPerspective]
	uPerspective.UniformMatrix4fv(false, context.MPerspective)

    uView := dm.UniformLocations[DustUnif_mView]
	uView.UniformMatrix4fv(false, context.MView)

    uCenterPosition := dm.UniformLocations[DustUnif_vCenterPosition]
	uCenterPosition.Uniform3f(context.VCamPos[0], context.VCamPos[1], 0.0)
}


func (dm *DustMaterial) Cleanup() {
	dm.DisableAttribs()
}

type DustArguments struct {
	BoxSize float32
	Corner Vec3
	Count int
}

func (dm *DustMaterial) Render(args interface{}) {
	var dustArgs = args.(DustArguments)
	var corner = dustArgs.Corner

	uBasePosition := dm.UniformLocations[DustUnif_vBasePosition]
	uBasePosition.Uniform3f(corner[0], corner[1], corner[2])

	uDustBoxSize := dm.UniformLocations[DustUnif_fDustBoxSize]
	uDustBoxSize.Uniform1f(dustArgs.BoxSize)

	gl.DrawArraysInstanced(gl.POINTS, 0, 1, dustArgs.Count)
}

