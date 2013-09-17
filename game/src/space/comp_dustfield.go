package space

import (
	"space/render"

)

var DUST_COUNT = 10

type Dustfield struct {
	Entity *Entity
	Physics *SpacePhysics

	material *render.DustfieldMaterial
	//stars []float32
}


func NewDustfield() *Dustfield {
	return &Dustfield {
		material: render.NewDustfieldMaterial(),
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
	s.material.Render(
		&context.MPerspective, &context.MView, context.VCamPos)
}
