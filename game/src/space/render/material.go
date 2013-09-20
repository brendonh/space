package render

import (
	"fmt"
	"math"

	"github.com/go-gl/gl"
)

type MaterialID uint64

var lastMaterialID MaterialID

func getMaterialID() MaterialID {
	if (lastMaterialID == math.MaxUint64) {
		panic("Ran out of IDs!")
	}
	lastMaterialID++
	return lastMaterialID
}

// ----------------------------------------------

type Material interface {
	Prepare(*Context)
	Cleanup()

	Render(args interface{})
}

type BaseMaterial struct {
	ID MaterialID
	Program *gl.Program
	AttribLocations []gl.AttribLocation
	UniformLocations []gl.UniformLocation
}

func NewBaseMaterial(shaderTag string, shaderSpecs... ShaderSpec) *BaseMaterial {
	program, err := ShaderCache.GetShader(shaderTag, shaderSpecs...)
	if err != nil {
		panic(fmt.Sprintf("Couldn't get cube material: %s", err))
	}

	return &BaseMaterial {
		ID: getMaterialID(),
		Program: program,
	}
}


func (m *BaseMaterial) GetAttribLocation(name string) gl.AttribLocation {
	loc := m.Program.GetAttribLocation(name)
	if loc == -1 {
		panic("Couldn't find attrib location: " + name)
	}
	return loc
}

func (m *BaseMaterial) GetUniformLocation(name string) gl.UniformLocation {
	loc := m.Program.GetUniformLocation(name)
	if loc == -1 {
		panic("Couldn't find uniform location: " + name)
	}
	return loc
}

func (m *BaseMaterial) EnableAttribs() {
	for _, attr := range m.AttribLocations {
		attr.EnableArray()
	}
}

func (m *BaseMaterial) DisableAttribs() {
	for _, attr := range m.AttribLocations {
		attr.DisableArray()
	}
}
