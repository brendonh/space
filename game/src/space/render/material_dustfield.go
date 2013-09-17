package render

import (
	"fmt"

	"github.com/go-gl/gl"
	. "github.com/brendonh/glvec"
)

type DustfieldMaterial struct {
	starBuffer gl.Buffer
}

func NewDustfieldMaterial() *DustfieldMaterial {
	return &DustfieldMaterial {
		starBuffer: gl.GenBuffer(),
	}
}

// XXX TODO: Parameterize
func (sm *DustfieldMaterial) Render(mP, mV *Mat4, camPos Vec3) {

	program, err := ShaderCache.GetShader("starfield", "starfield.vert", "starfield.frag")
	if err != nil {
		panic(fmt.Sprintf("Couldn't get starfield material: %s", err))
	}

	program.Use()

    sm.starBuffer.Bind(gl.ARRAY_BUFFER)

	var firstPoint = []float32{ 0.0, 0.0, 0.0, 1.0 }

    gl.BufferData(gl.ARRAY_BUFFER, 3 * 4, firstPoint, gl.STATIC_DRAW);

    aStarPosition := program.GetAttribLocation("aStarPosition")
    aStarPosition.AttribPointer(3, gl.FLOAT, false, 0, uintptr(0))
	aStarPosition.EnableArray()
	defer aStarPosition.DisableArray()

    uCenterPosition := program.GetUniformLocation("uCenterPosition")
	uCenterPosition.Uniform3f(camPos[0], camPos[1], 0.0)

    uPerspective := program.GetUniformLocation("uPMatrix")
    uPerspective.UniformMatrix4fv(false, *mP)
	
    uView := program.GetUniformLocation("uVMatrix")
	uView.UniformMatrix4fv(false, *mV)

	startX := floorMod(camPos[0], 5.0)
	startY := floorMod(camPos[1], 5.0)

    uBasePosition := program.GetUniformLocation("uBasePosition")

	for x := startX - 5; x <= startX + 5; x += 5 {
		for y := startY - 5; y <= startY + 5; y += 5 {
			uBasePosition.Uniform3f(float32(x), float32(y), -2.5)
			gl.DrawArraysInstanced(gl.POINTS, 0, 1, 20)
		}
	}
}


func floorMod (val, quot float32) float64 {
	q := val / quot
	if val < 0 {
		q -= 1
	}
	return float64(int(q)) * float64(quot)
}