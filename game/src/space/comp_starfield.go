package space

import (
	"math/rand"

	"space/render"

	. "github.com/brendonh/glvec"
)

var STAR_COUNT = 100

type Starfield struct {
	material *render.StarfieldMaterial
	stars []Vec3
}


func NewStarfield() *Starfield {
	var stars = make([]Vec3, 0, STAR_COUNT)
	for i := 0; i < STAR_COUNT; i++ {
		stars = append(stars, Vec3 {
			(rand.Float32() * 10)- 5,
			(rand.Float32() * 6) - 3,
			(rand.Float32() * -5) + 1,
		})
	}
	
	return &Starfield {
		material: render.NewStarfieldMaterial(),
		stars: stars,
	}
}

func (s *Starfield) Render(context *RenderContext) {
	s.material.Render(&context.MPerspective, &context.MView, s.stars)
}
