package render

import (
	"fmt"
	"math/rand"

	"github.com/go-gl/gl"
	. "github.com/brendonh/glvec"
)

type StarfieldMaterial struct {
	starBuffer gl.Buffer
	printed bool
}

func NewStarfieldMaterial() *StarfieldMaterial {
	return &StarfieldMaterial {
		starBuffer: gl.GenBuffer(),
		printed: false,
	}
}


func (sm *StarfieldMaterial) Render(mP, mV *Mat4, stars []float32) {

	program, err := ShaderCache.GetShader("starfield", "starfield.vert", "starfield.frag")
	if err != nil {
		panic(fmt.Sprintf("Couldn't get starfield material: %s", err))
	}

	program.Use()

    sm.starBuffer.Bind(gl.ARRAY_BUFFER)

	var proj, invProj Mat4
	M4MulM4(&proj, mP, mV)
	M4Inverse(&invProj, &proj)

	for i := 0; i < len(stars); i += 4 {
		var starPos = Vec4 { stars[i], stars[i+1], stars[i+2], 1.0 }
		var clipPos Vec4
		
		M4MulV4(&clipPos, &proj, starPos)
		var w = clipPos[3]
		V4ScalarDiv(&clipPos, &clipPos, w)
		
		if clipPos[0] < -1.1 {
			clipPos[0] = 1.0 + (rand.Float32() * 0.1)
			clipPos[1] = (rand.Float32() * 2.2) - 1.1
		} else if clipPos[0] > 1.1 {
			clipPos[0] = -(1.0 + (rand.Float32() * 0.1))
			clipPos[1] = (rand.Float32() * 2.2) - 1.1
		} else if clipPos[1] < -1.1 {
			clipPos[1] = 1.0 + (rand.Float32() * 0.1)
			clipPos[0] = (rand.Float32() * 2.2) - 1.1
		} else if clipPos[1] > 1.1 {
			clipPos[1] = -(1.0 + (rand.Float32() * 0.1))
			clipPos[0] = (rand.Float32() * 2.2) - 1.1
		} else {
			continue
		}

		M4MulV4(&starPos, &invProj, clipPos)
		V4ScalarMul(&starPos, &starPos, w)
		stars[i] = starPos[0]
		stars[i+1] = starPos[1]
	}

    gl.BufferData(gl.ARRAY_BUFFER, len(stars) * 4, stars, gl.DYNAMIC_DRAW);

    aVertexPosition := program.GetAttribLocation("aStarPosition")
    aVertexPosition.AttribPointer(4, gl.FLOAT, false, 0, uintptr(0))
	aVertexPosition.EnableArray()
	defer aVertexPosition.DisableArray()

    pUniform := program.GetUniformLocation("uPMatrix")
    pUniform.UniformMatrix4fv(false, *mP)

    uView := program.GetUniformLocation("uVMatrix")
	uView.UniformMatrix4fv(false, *mV)

	gl.DrawArrays(gl.POINTS, 0, len(stars) / 4)
}


