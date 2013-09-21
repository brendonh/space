package space

import (
	"space/render"

	. "github.com/brendonh/glvec"
)

type Dustfield struct {
	BaseComponent

	Physics *SpacePhysics
	MaterialID render.MaterialID
	Density int
}


func NewDustfield() *Dustfield {
	return &Dustfield {
		MaterialID: render.GetDustMaterialID(),
		Density: 20,
	}
}

func (s *Dustfield) Init() {
	s.Physics = s.Entity.GetComponent("struct_spacephysics").(*SpacePhysics)
}

func (s *Dustfield) Tag() string {
	return "dust"
}

func (s *Dustfield) Render(context *render.Context, alpha float64) {
	startX := floorMod(context.VCamPos[0], 5.0)
	startY := floorMod(context.VCamPos[1], 5.0)

	for x := startX - 5; x <= startX + 5; x += 5 {
		for y := startY - 5; y <= startY + 5; y += 5 {
			var corner = Vec3 { float32(x), float32(y), -2.5 }
			context.Enqueue(s.MaterialID, render.DustArguments{corner, s.Density})
		}
	}
}

func floorMod (val, quot float32) int {
	q := val / quot
	if val <= 0 {
		q -= 1
	}
	return int(q) * int(quot)
}