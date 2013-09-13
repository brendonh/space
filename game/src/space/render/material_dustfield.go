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


func (sm *DustfieldMaterial) Render(mP, mV *Mat4, stars []float32) {

	program, err := ShaderCache.GetShader("starfield", "starfield.vert", "starfield.frag")
	if err != nil {
		panic(fmt.Sprintf("Couldn't get starfield material: %s", err))
	}

	program.Use()

    sm.starBuffer.Bind(gl.ARRAY_BUFFER)

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


