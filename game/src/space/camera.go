package space

import (
	"space/render"

	. "github.com/brendonh/glvec"
)


type Camera struct {
	VCamTranslate Vec3
	FollowPhysics *SpacePhysics
}

func NewCamera() *Camera {
	return &Camera {
		//VCamTranslate: Vec3 { 8.0, 8.0, 8.0 },
		VCamTranslate: Vec3 { 0.0, -20.0, 30.0 },
	}
}

func (cam *Camera) FollowEntity(e *Entity) {
	cam.FollowPhysics = e.GetComponent("struct_spacephysics").(*SpacePhysics)
}

func (cam *Camera) UpdateRenderContext(context *render.Context, alpha float64) {
	var phys = cam.FollowPhysics
	var pos = phys.InterpolatePosition(alpha)
	
	var center = Vec3 { float32(pos.PosX), float32(pos.PosY), 0.0 }

	V3Add(&context.VCamPos, center, cam.VCamTranslate)

	var up = Vec3 { 0.0, 0.0, 1.0 }
	M4LookAt(&context.MView, context.VCamPos, center, up)

	M4RotationMatrix(&context.MCamRotate, &context.MView)
}