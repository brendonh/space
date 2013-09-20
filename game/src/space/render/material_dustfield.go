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
)

type DustMaterial struct {
	*BaseMaterial
	starBuffer gl.Buffer
}

func NewDustMaterial() *DustMaterial {
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
	}

	return m
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


// XXX TODO: Parameterize
func (dm *DustMaterial) Render(args interface{}) {
	var camPos = args.(Vec3)

	startX := floorMod(camPos[0], 5.0)
	startY := floorMod(camPos[1], 5.0)

	uBasePosition := dm.UniformLocations[DustUnif_vBasePosition]

	for x := startX - 5; x <= startX + 5; x += 5 {
		for y := startY - 5; y <= startY + 5; y += 5 {
			uBasePosition.Uniform3f(float32(x), float32(y), -2.5)
			gl.DrawArraysInstanced(gl.POINTS, 0, 1, 20)
		}
	}
}


func floorMod (val, quot float32) int {
	q := val / quot
	if val <= 0 {
		q -= 1
	}
	return int(q) * int(quot)
}