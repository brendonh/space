package space

import (
	"math"

	. "github.com/brendonh/glvec"
)

type Vec2i struct {
	X, Y int
}

func (v Vec2i) Add(o Vec2i) Vec2i {
	return Vec2i{
		v.X + o.X,
		v.Y + o.Y,
	}
}

func (v Vec2i) Neg() Vec2i {
	return Vec2i { -v.X, -v.Y }
}

func (v Vec2i) Distance(o Vec2i) float64 {
	dx := float64(o.X - v.X)
	dy := float64(o.Y - v.Y)
	return math.Sqrt((dx*dx) + (dy*dy))
}

func (v Vec2i) Vec3() Vec3 {
	return Vec3{ float32(v.X), float32(v.Y), 0 }
}