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

func (cam *Camera) UpdateRenderContext(context *render.Context, alpha float32) {
	var phys = cam.FollowPhysics
	var pos = phys.InterpolatePosition(alpha)
	
	V3Add(&context.VCamPos, pos.Pos, cam.VCamTranslate)

	var up = Vec3 { 0.0, 0.0, 1.0 }
	M4LookAt(&context.MView, context.VCamPos, pos.Pos, up)

	M4RotationMatrix(&context.MCamRotate, &context.MView)
}