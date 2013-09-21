package space

import (
	"space/render"
)

type Dustfield struct {
	Entity *Entity
	Physics *SpacePhysics
	MaterialID render.MaterialID
}


func NewDustfield() *Dustfield {
	return &Dustfield {
		MaterialID: render.GetDustMaterialID(),
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
	context.Enqueue(s.MaterialID, context.VCamPos)
}
