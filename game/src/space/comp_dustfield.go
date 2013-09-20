package space

import (
	"space/render"
)

type Dustfield struct {
	Entity *Entity
	Physics *SpacePhysics

	//material *render.DustMaterial
}


func NewDustfield() *Dustfield {
	return &Dustfield {
		//material: render.NewDustfieldMaterial(),
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

func (s *Dustfield) Render(context *render.Context, alpha float64) {
	
	// Temp
	mat := render.NewDustMaterial()
	mat.Prepare(context)
	mat.Render(context.VCamPos)
	mat.Cleanup()

	// s.material.Render(
	// 	&context.MPerspective, &context.MView, context.VCamPos)
}
