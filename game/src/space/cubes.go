package space

import (
	"math"

	. "github.com/brendonh/glvec"
)

const (
	CUBE_TOP uint8 = 1 << iota
	CUBE_LEFT      
	CUBE_BOTTOM    
	CUBE_RIGHT     
	CUBE_FRONT     
	CUBE_BACK      
)

type CubeFaces uint8

func CubeFacesAll() CubeFaces {
	return 255
}

func (cf *CubeFaces) Set(face uint8) {
	*cf |= CubeFaces(face)
}

func (cf *CubeFaces) Unset(face uint8) {
	*cf &= CubeFaces(^face)
}

func (cf *CubeFaces) Has(face uint8) bool {
	return *cf & CubeFaces(face) > 0
}



type CubeColor struct {
	R, G, B float32
}

type Cube struct {
	X, Y int
	Color CubeColor
	Faces CubeFaces
}

func (c Cube) Center() Vec3 {
	return Vec3{ float32(c.X), float32(c.Y), 0 }
}

type CubeSet struct {
	Cubes []Cube
	Center Vec3
}


func addCubeFaces(buffer []float32, cube Cube, center Vec3) []float32 {
	var yAxis = Vec3{ 0.0, 1.0, 0.0 }
	var xAxis = Vec3{ 1.0, 0.0, 0.0 }

	var pos = cube.Center()
	V3Add(&pos, pos, center)
	V3ScalarMul(&pos, pos, 2.0)

	var color = Vec3{ cube.Color.R, cube.Color.G, cube.Color.B }

	var q Quat

	if cube.Faces.Has(CUBE_FRONT) {
		QIdent(&q)
		buffer = addCubeFace(buffer, q, pos, Vec3 {0.0, 0.0, -1.0}, color)
	}

	if cube.Faces.Has(CUBE_BACK) {
		QRotAng(&q, math.Pi, yAxis)
		buffer = addCubeFace(buffer, q, pos, Vec3 {0.0, 0.0, 1.0}, color)
	}

	if cube.Faces.Has(CUBE_TOP) {
		QRotAng(&q, math.Pi / 2, xAxis)
		buffer = addCubeFace(buffer, q, pos, Vec3 {0.0, -1.0, 0.0}, color)
	}

	if cube.Faces.Has(CUBE_LEFT) {
		QRotAng(&q, math.Pi / 2, yAxis)
		buffer = addCubeFace(buffer, q, pos, Vec3 {1.0, 0.0, 0.0}, color)
	}

	if cube.Faces.Has(CUBE_BOTTOM) {
		QRotAng(&q, -math.Pi / 2, xAxis)
		buffer = addCubeFace(buffer, q, pos, Vec3 {0.0, 1.0, 0.0}, color)
	}

	if cube.Faces.Has(CUBE_RIGHT) {
		QRotAng(&q, -math.Pi/2, yAxis)
		buffer = addCubeFace(buffer, q, pos, Vec3 {-1.0, 0.0, 0.0}, color)
	}

	return buffer
}



var vertVecs = []Vec3 {
	Vec3 {  1.0,  1.0,  1.0 },
	Vec3 { -1.0,  1.0,  1.0 },
	Vec3 {  1.0, -1.0,  1.0 },
	Vec3 {  1.0, -1.0,  1.0 },
	Vec3 { -1.0,  1.0,  1.0 },
	Vec3 { -1.0, -1.0,  1.0 },
}

func addCubeFace(buffer []float32, rot Quat, pos Vec3, normal Vec3, color Vec3) []float32 {
	for _, v := range vertVecs {
		var temp Mat3
		QMat3(&temp, rot)
		M3MulV3(&v, &temp, v)

		buffer = append(buffer, 
			round(v[0]) + pos[0], 
			round(v[1]) + pos[1], 
			round(v[2]) + pos[2])

		buffer = append(buffer, normal[0], normal[1], normal[2])
		buffer = append(buffer, color[0], color[1], color[2])
	}
	return buffer
}

var edgeVecs = []Vec3 {
}

func addCubeEdges(buffer []float32, cube Cube, center Vec3) []float32 {
	var pos Vec3 = cube.Center()
	V3Add(&pos, pos, center)
	V3ScalarMul(&pos, pos, 2.0)

	var addEdge = func(verts... Vec3) {
		for _, v := range verts {
			var nv Vec3
			V3Add(&nv, pos, v)
			buffer = append(buffer, nv[0], nv[1], nv[2])
		}
	}

	if cube.Faces.Has(CUBE_TOP) {
		addEdge(Vec3 { -1.0,  1.0, 1.001 }, Vec3 {  1.0,  1.0, 1.001 })
	}
	
	//if cube.Left {
		addEdge(Vec3 {  -1.0,  1.0, 1.001 }, Vec3 {  -1.0, -1.0, 1.001 })
	//}

	//if cube.Bottom {
		addEdge(Vec3 { -1.0, -1.0, 1.001 }, Vec3 { 1.0, -1.0, 1.001 })
	//}

	if cube.Faces.Has(CUBE_RIGHT) {
		addEdge(Vec3 { 1.0,  1.0, 1.001 }, Vec3 { 1.0, -1.0, 1.001 })
	}

	return buffer
}

func round(f float32) float32 {
	if f < 0 {
		return float32(math.Ceil(float64(f) - 0.5))
	}
	return float32(math.Floor(float64(f) + 0.5))
}