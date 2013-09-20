package render

import (
	"github.com/go-gl/gl"
)

type Material interface {
	Init()
	Prepare(*Context)
	Render(args interface{})
	Cleanup()
}

type BaseMaterial struct {
	Program *gl.Program
	AttribLocations []gl.AttribLocation
	UniformLocations []gl.UniformLocation
}