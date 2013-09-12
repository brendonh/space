package components

import (
	"math"

	. "space"
	"space/render"

	"github.com/go-gl/gl"
	. "github.com/brendonh/glvec"
)

type CubesComponent struct {
	verts gl.Buffer
	count int
}

func NewCubesComponent() *CubesComponent {
	return &CubesComponent {
		verts: makeCubeBuffer(),
		count: 36,
	}
}

func (c *CubesComponent) Render(context *RenderContext) {
	// XXX TODO: Mix M into mV
	render.RenderCubeMaterial(
		&context.MPerspective, &context.MView, context.VLightDir,
		c.verts, c.count)
}



func makeCubeBuffer() gl.Buffer {
    
	var buf = make([]float32, (3+3+3)*6*6)
    var bufOffset = 0

    var vertVecs = []Vec3 {
		Vec3 {  1.0,  1.0,  1.0 },
		Vec3 { -1.0,  1.0,  1.0 },
		Vec3 {  1.0, -1.0,  1.0 },
		Vec3 {  1.0, -1.0,  1.0 },
		Vec3 { -1.0,  1.0,  1.0 },
		Vec3 { -1.0, -1.0,  1.0 },
	}

    var addFace = func(rot Quat, normal Vec3, color Vec3) {

		var round = func(f float32) float32 {
			if f < 0 {
				return float32(math.Ceil(float64(f) - 0.5))
			}
			return float32(math.Floor(float64(f) + 0.5))
		}
        
		for _, v := range vertVecs {
			var temp Mat3
			QMat3(&temp, &rot)
			M3MulV3(&v, &temp, &v)
		
            // Vert
            buf[bufOffset]     = round(v[0])
            buf[bufOffset + 1] = round(v[1])
            buf[bufOffset + 2] = round(v[2])

            // Normal
            buf[bufOffset + 3] = normal[0]
            buf[bufOffset + 4] = normal[1]
            buf[bufOffset + 5] = normal[2]

            // Color
            buf[bufOffset + 6] = color[0]
            buf[bufOffset + 7] = color[1]
            buf[bufOffset + 8] = color[2]

            bufOffset += 9
        }
    }

    // Front
	var q Quat
	QIdent(&q)

	var yAxis = Vec3 { 0.0, 1.0, 0.0 }
	var xAxis = Vec3 { 1.0, 0.0, 0.0 }

    addFace(q, Vec3 {0.0, 0.0, -1.0}, Vec3 {1.0, 0.0, 0.0})

    // Right
	QRotAng(&q, math.Pi / 2, &yAxis)
    addFace(q, Vec3 {1.0, 0.0, 0.0}, Vec3 {1.0, 1.0, 0.0})

    // Back
	QRotAng(&q, math.Pi, &yAxis)
    addFace(q, Vec3 {0.0, 0.0, 1.0}, Vec3 {0.0, 1.0, 0.0})

    // Left
	QRotAng(&q, -math.Pi/2, &yAxis)
    addFace(q, Vec3 {-1.0, 0.0, 0.0}, Vec3 {0.0, 1.0, 1.0})

    // Top
	QRotAng(&q, math.Pi / 2, &xAxis)
    addFace(q, Vec3 {0.0, -1.0, 0.0}, Vec3 {0.0, 0.0, 1.0})

    // Bottom
	QRotAng(&q, -math.Pi / 2, &xAxis)
    addFace(q, Vec3 {0.0, 1.0, 0.0}, Vec3 {1.0, 0.0, 1.0});

    glBuf := gl.GenBuffer()
    glBuf.Bind(gl.ARRAY_BUFFER)
    gl.BufferData(gl.ARRAY_BUFFER, len(buf) * 4, buf, gl.STATIC_DRAW);

    return glBuf;
}
