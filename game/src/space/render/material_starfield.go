package render

import (
	"fmt"

	"github.com/go-gl/gl"
	. "github.com/brendonh/glvec"
)

type StarfieldMaterial struct {
	quad gl.Buffer
}

func NewStarfieldMaterial() *StarfieldMaterial {
	var quadVerts = []float32 {
 		 1.0,  1.0, 0.0,
		-1.0,  1.0, 0.0,
		 1.0, -1.0, 0.0,
		 1.0, -1.0, 0.0,
		-1.0,  1.0, 0.0,
		-1.0, -1.0, 0.0,
	}

    glBuf := gl.GenBuffer()
    glBuf.Bind(gl.ARRAY_BUFFER)
    gl.BufferData(gl.ARRAY_BUFFER, len(quadVerts) * 4, quadVerts, gl.STATIC_DRAW);

	return &StarfieldMaterial {
		quad: glBuf,
	}
}


func (sm *StarfieldMaterial) Render(mP, mV *Mat4, stars []Vec3) {

	program, err := ShaderCache.GetShader("starfield", "starfield.vert", "starfield.frag")
	if err != nil {
		panic(fmt.Sprintf("Couldn't get starfield material: %s", err))
	}

	program.Use()

    sm.quad.Bind(gl.ARRAY_BUFFER)

    aVertexPosition := program.GetAttribLocation("aVertexPosition")
    aVertexPosition.AttribPointer(3, gl.FLOAT, false, 0, uintptr(0))
	aVertexPosition.EnableArray()
	defer aVertexPosition.DisableArray()

    pUniform := program.GetUniformLocation("uPMatrix")
    pUniform.UniformMatrix4fv(false, *mP)

    mvUniform := program.GetUniformLocation("uMVMatrix")

	var mM, mTranslate, mMV Mat4	

	for _, starPos := range stars {
		M4MakeScale(&mM, 0.01)
		M4MakeTransform(&mTranslate, &starPos)

		M4MulM4(&mM, &mTranslate, &mM)
		M4MulM4(&mMV, mV, &mM)

		mvUniform.UniformMatrix4fv(false, mMV)
		gl.DrawArrays(gl.TRIANGLES, 0, 6)
	}
}


