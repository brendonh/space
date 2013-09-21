package render

import (
	"fmt"

	"github.com/go-gl/gl"
)

type MaterialID int64

var materials []Material
var materialsByName map[string]Material

func init() {
	materialsByName = make(map[string]Material)
}

func registerMaterial(m Material, name string) MaterialID {
	materials = append(materials, m)

	_, ok := materialsByName[name]
	if ok {
		panic("Duplicate material name: " + name)
	}

	materialsByName[name] = m
	return MaterialID(len(materials) - 1)
}

func GetMaterialID(name string) (MaterialID, bool) {
	m, ok := materialsByName[name]
	if !ok {
		return 0, false
	}
	return m.GetID(), true
}

// ----------------------------------------------

type Material interface {
	GetID() MaterialID
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
		Program: program,
	}
}

func (m *BaseMaterial) GetID() MaterialID {
	return m.ID
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
