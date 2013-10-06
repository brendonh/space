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

var DUST_BOX = 20

func (s *Dustfield) Render(context *render.Context, alpha float32) {
	startX := floorMod(context.VCamPos[0], float32(DUST_BOX))
	startY := floorMod(context.VCamPos[1], float32(DUST_BOX))

	for x := startX - DUST_BOX; x <= startX + DUST_BOX; x += DUST_BOX {
		for y := startY - DUST_BOX; y <= startY + DUST_BOX; y += DUST_BOX {
			var corner = Vec3 { float32(x), float32(y), float32(DUST_BOX) / 2 }
			context.Enqueue(s.MaterialID, 
				render.DustArguments{float32(DUST_BOX), corner, s.Density})
		}
	}
}

func (s *Dustfield) HandleMouse(ray Ray) bool {
	return false
}

func floorMod (val, quot float32) int {
	q := val / quot
	if val <= 0 {
		q -= 1
	}
	return int(q) * int(quot)
}