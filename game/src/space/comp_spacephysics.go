package space

import (
	"math"

	. "github.com/brendonh/glvec"
)


type SpaceForce struct {
	Angle float64
	Acc float64
}

type SpacePhysics struct {
	PosX float64
	PosY float64

	VelX float64
	VelY float64

	Angle float64

	Forces []SpaceForce
}

func (c *SpacePhysics) Init() {
}

func (c *SpacePhysics) Tag() string {
	return "struct_spacephysics"
}

func (c *SpacePhysics) SetEntity (e *Entity) {
	// TODO
}

func (c *SpacePhysics) TickPhysics() {
	for _, force := range c.Forces {
		c.VelX += force.Acc * -math.Sin(force.Angle)
		c.VelY += force.Acc * math.Cos(force.Angle)
	}
	c.Forces = nil

	c.PosX += c.VelX
	c.PosY += c.VelY


}

func (c *SpacePhysics) ApplyForce(Angle float64, Acc float64) {
	c.Forces = append(c.Forces, SpaceForce{ Angle, Acc })
}

// TODO: Make relative to reference point
func (c *SpacePhysics) GetModelMatrix() *Mat4 {
	var result Mat4
	M4MakeRotation(&result, float32(c.Angle), &Vec3 { 0.0, 0.0, 1.0 })
	M4SetTransform(&result, &Vec3 { float32(c.PosX), float32(c.PosY), 0.0 })
	return &result
}
