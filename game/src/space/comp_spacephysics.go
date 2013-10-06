package space

import (
	"math"

	. "github.com/brendonh/glvec"
)


type SpaceForce struct {
	Angle float32
	Acc float32
}

type SpaceThrust float32

type SpaceRotation float32

type SpacePosition struct {
	Pos Vec3
	Angle float32
}

func (pos SpacePosition) InterpolateFrom(other SpacePosition, alpha float32) SpacePosition {
	return SpacePosition {
		Pos: Vec3{
			other.Pos[0] + ((pos.Pos[0] - other.Pos[0]) * alpha),
			other.Pos[1] + ((pos.Pos[1] - other.Pos[1]) * alpha),
			0,
		},
		Angle: other.Angle + ((pos.Angle - other.Angle) * alpha),
	}
}



type SpacePhysics struct {
	BaseComponent

	Position SpacePosition
	PrevPosition SpacePosition

	VelX float32
	VelY float32

	Rotations []SpaceRotation
	Thrusts []SpaceThrust
	Forces []SpaceForce
}

func (c *SpacePhysics) Tag() string {
	return "struct_spacephysics"
}

func (c *SpacePhysics) TickPhysics() {
	c.PrevPosition = c.Position

	var pos = &c.Position

	for _, rot := range c.Rotations {
		pos.Angle += float32(rot)
	}
	c.Rotations = c.Rotations[:0]

	for _, thrust := range c.Thrusts {
		c.ApplyForce(pos.Angle, float32(thrust))
	}
	c.Thrusts = c.Thrusts[:0] 

	for _, force := range c.Forces {
		c.VelX += force.Acc * -float32(math.Sin(float64(force.Angle)))
		c.VelY += force.Acc * float32(math.Cos(float64(force.Angle)))
	}
	c.Forces = c.Forces[:0]

	c.Position.Pos[0] += c.VelX
	c.Position.Pos[1] += c.VelY
}

func (c *SpacePhysics) ApplyForce(Angle float32, Acc float32) {
	c.Forces = append(c.Forces, SpaceForce{ Angle, Acc })
}

func (c *SpacePhysics) ApplyRotation(Angle float32) {
	c.Rotations = append(c.Rotations, SpaceRotation(Angle))
}

func (c *SpacePhysics) ApplyThrust(Acc float32) {
	c.Thrusts = append(c.Thrusts, SpaceThrust(Acc))
}

func (c *SpacePhysics) InterpolatePosition(alpha float32) SpacePosition {
	return c.Position.InterpolateFrom(c.PrevPosition, alpha)
}

// TODO: Make relative to reference point
func (c *SpacePhysics) GetModelMatrix(alpha float32) Mat4 {
	var result Mat4
	var pos = c.InterpolatePosition(alpha)
	M4MakeRotation(&result, float32(pos.Angle), Vec3{ 0, 0, 1 })
	M4SetTransform(&result, pos.Pos)
	return result
}
