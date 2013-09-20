package render

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/go-gl/gl"
)

var BASE_PATH = "data/shaders/"

type shaderCache struct {
	shaders map[string]gl.Program
}

var ShaderCache = &shaderCache {
	shaders: make(map[string]gl.Program),
}

type ShaderSpec struct {
	Type gl.GLenum
	Name string
}

func (c *shaderCache) GetShader(tag string, shaderSpecs... ShaderSpec) (*gl.Program, error) {
	program, ok := c.shaders[tag]
	if ok {
		return &program, nil
	}

    program = gl.CreateProgram();

	for _, spec := range shaderSpecs {
		source, err := loadShaderSource(spec.Name)
		if err != nil {
			return nil, err
		}

		shader, err := compileShader(spec.Type, source)
		if err != nil {
			return nil, err
		}

		program.AttachShader(*shader)
	}

    program.Link()

	if program.Get(gl.LINK_STATUS) != gl.TRUE {
		panic("linker error: " + program.GetInfoLog())
	}

	c.shaders[tag] = program

	return &program, nil
}


func compileShader(shaderType gl.GLenum, source string) (*gl.Shader, error) {
    var shader = gl.CreateShader(shaderType)
    shader.Source(source)
    shader.Compile()

    if shader.Get(gl.COMPILE_STATUS) != gl.TRUE {
        return nil, errors.New(fmt.Sprintf("Error compiling shader: %v\n%s", 
			shader.GetInfoLog(), source))
    }

    return &shader, nil
}


func loadShaderSource(filename string) (string, error) {
	content, err := ioutil.ReadFile(BASE_PATH + filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

