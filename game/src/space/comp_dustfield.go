package space

import (
	"math/rand"
	"time"

	"space/render"
)

var DUST_COUNT = 100

type Dustfield struct {
	Entity *Entity
	Physics *SpacePhysics

	material *render.DustfieldMaterial
	stars []float32
}


func NewDustfield() *Dustfield {

	// Temporary
	rand.Seed(time.Now().Unix())

	var stars = make([]float32, 0, DUST_COUNT * 4)
	for i := 0; i < DUST_COUNT; i++ {
		stars = append(stars,
			(rand.Float32() * 20) - 10,
			(rand.Float32() * 20) - 10,
			(rand.Float32() * 10) - 8,
			rand.Float32() + 1,
		)
	}

	return &Dustfield {
		material: render.NewDustfieldMaterial(),
		stars: stars,
	}
}

func (s *Dustfield) Init() {
}

func (s *Dustfield) Tag() string {
	return ""
}

func (s *Dustfield) SetEntity(e *Entity) {
	s.Entity = e
	s.Physics = e.GetComponent("struct_spacephysics").(*SpacePhysics)
}

func (s *Dustfield) Render(context *RenderContext) {
	s.material.Render(&context.MPerspective, &context.MView, s.stars)
}
