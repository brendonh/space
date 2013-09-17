package space

import (
	"math"

	. "github.com/brendonh/glvec"
)


type SpaceForce struct {
	Angle float64
	Acc float64
}

type SpaceThrust float64

type SpaceRotation float64

type SpacePosition struct {
	PosX float64
	PosY float64

	Angle float64
}

type SpacePhysics struct {
	Entity *Entity

	Position SpacePosition
	PrevPosition SpacePosition

	VelX float64
	VelY float64

	Rotations []SpaceRotation
	Thrusts []SpaceThrust
	Forces []SpaceForce
}

func (c *SpacePhysics) Init() {
}

func (c *SpacePhysics) Tag() string {
	return "struct_spacephysics"
}

func (c *SpacePhysics) SetEntity (e *Entity) {
	c.Entity = e
}

func (c *SpacePhysics) TickPhysics() {
	c.PrevPosition = c.Position

	var pos = &c.Position

	for _, rot := range c.Rotations {
		pos.Angle += float64(rot)
	}
	c.Rotations = c.Rotations[:0]

	for _, thrust := range c.Thrusts {
		c.ApplyForce(pos.Angle, float64(thrust))
	}
	c.Thrusts = c.Thrusts[:0] 

	for _, force := range c.Forces {
		c.VelX += force.Acc * -math.Sin(force.Angle)
		c.VelY += force.Acc * math.Cos(force.Angle)
	}
	c.Forces = c.Forces[:0]

	c.Position.PosX += c.VelX
	c.Position.PosY += c.VelY
}

func (c *SpacePhysics) ApplyForce(Angle float64, Acc float64) {
	c.Forces = append(c.Forces, SpaceForce{ Angle, Acc })
}

func (c *SpacePhysics) ApplyRotation(Angle float64) {
	c.Rotations = append(c.Rotations, SpaceRotation(Angle))
}

func (c *SpacePhysics) ApplyThrust(Acc float64) {
	c.Thrusts = append(c.Thrusts, SpaceThrust(Acc))
}

func (c *SpacePhysics) InterpolatePosition(delta float64) SpacePosition {
	prev, pos := c.PrevPosition, c.Position
	return SpacePosition {
		PosX: prev.PosX + ((pos.PosX - prev.PosX) * delta),
		PosY: prev.PosY + ((pos.PosY - prev.PosY) * delta),
		Angle: prev.Angle + ((pos.Angle - prev.Angle) * delta),
	}
}

// TODO: Make relative to reference point
func (c *SpacePhysics) GetModelMatrix(delta float64) *Mat4 {
	var result Mat4
	var pos = c.InterpolatePosition(delta)
	M4MakeRotation(&result, float32(pos.Angle), Vec3 { 0.0, 0.0, 1.0 })
	M4SetTransform(&result, Vec3 { float32(pos.PosX), float32(pos.PosY), 0.0 })
	return &result
}
