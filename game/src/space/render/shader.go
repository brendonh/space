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

func (c *shaderCache) GetShader(tag string, vert string, frag string) (*gl.Program, error) {
	program, ok := c.shaders[tag]
	if ok {
		return &program, nil
	}

	vertSource, err := loadShaderSource(vert)
	if err != nil {
		return nil, err
	}

	fragSource, err := loadShaderSource(frag)
	if err != nil {
		return nil, err
	}

	vertShader, err := compileShader(gl.VERTEX_SHADER, vertSource)
	if err != nil {
		return nil, err
	}

	fragShader, err := compileShader(gl.FRAGMENT_SHADER, fragSource)
	if err != nil {
		return nil, err
	}

    program = gl.CreateProgram();
	program.AttachShader(*vertShader)
	program.AttachShader(*fragShader)
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

