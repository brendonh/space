package render

import (
	"fmt"

	"github.com/go-gl/gl"
	. "github.com/brendonh/glvec"
)

func RenderCubeMaterial(mP, mMV *Mat4, vLight Vec3, verts gl.Buffer, count int) {
	program, err := ShaderCache.GetShader("cube", "cube.vert", "cube.frag")
	if err != nil {
		panic(fmt.Sprintf("Couldn't get cube material: %s", err))
	}

	program.Use()

    verts.Bind(gl.ARRAY_BUFFER)

    aVertexPosition := program.GetAttribLocation("aVertexPosition")
    aVertexPosition.AttribPointer(3, gl.FLOAT, false, 9 * 4, uintptr(0))
	aVertexPosition.EnableArray()
	defer aVertexPosition.DisableArray()

    aVertexNormal := program.GetAttribLocation("aNormal")
    aVertexNormal.AttribPointer(3, gl.FLOAT, false, 9 * 4, uintptr(3 * 4))
	aVertexNormal.EnableArray()
	defer aVertexNormal.DisableArray()

    aVertexColor := program.GetAttribLocation("aVertexColor")
    aVertexColor.AttribPointer(3, gl.FLOAT, false, 9 * 4, uintptr(6 * 4))
	aVertexColor.EnableArray()
	defer aVertexNormal.DisableArray()

    pUniform := program.GetUniformLocation("uPMatrix")
    pUniform.UniformMatrix4fv(false, *mP)

    mvUniform := program.GetUniformLocation("uMVMatrix")
    mvUniform.UniformMatrix4fv(false, *mMV)

	var mMVN Mat3
	M4RotationMatrix(&mMVN, mMV)
    mvNormal := program.GetUniformLocation("uNormalMatrix")
    mvNormal.UniformMatrix3fv(false, mMVN)

    uLightDirection := program.GetUniformLocation("uLightDirection")
    uLightDirection.Uniform3f(vLight[0],vLight[1],vLight[2])

    gl.DrawArrays(gl.TRIANGLES, 0, count*3)
}


