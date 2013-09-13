package space

import (
	"math/rand"
	"time"

	"space/render"

	//. "github.com/brendonh/glvec"
)

var DUST_COUNT = 100

type Starfield struct {
	material *render.StarfieldMaterial
	stars []float32
}


func NewStarfield() *Starfield {

	// Temporary
	rand.Seed(time.Now().Unix())

	var stars = make([]float32, 0, DUST_COUNT * 4)
	for i := 0; i < DUST_COUNT; i++ {
		stars = append(stars,
			(rand.Float32() * 10) - 5,
			(rand.Float32() * 10) - 5,
			(rand.Float32() * 10) - 8,
			rand.Float32() + 1,
		)
	}
	
	return &Starfield {
		material: render.NewStarfieldMaterial(),
		stars: stars,
	}
}

func (s *Starfield) Render(context *RenderContext) {
	s.material.Render(&context.MPerspective, &context.MView, s.stars)
}
