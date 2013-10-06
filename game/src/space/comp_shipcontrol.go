package space

import (
	"math"
)

type ShipControl struct {
	BaseComponent

	Physics *SpacePhysics
	Thrust float32
	Brake float32
	Turn float32
}

func (c *ShipControl) Init() {
	c.Physics = c.Entity.GetComponent("struct_spacephysics").(*SpacePhysics)
}

func (c *ShipControl) Tag() string {
	return "struct_shipcontrol"
}

func (c *ShipControl) TickLogic() {
	var physics = c.Physics

	physics.ApplyRotation(c.Turn * 0.05)
	physics.ApplyThrust(c.Thrust * 0.0005)
	
	if (c.Brake > 0.005) {
		x, y := physics.VelX, physics.VelY
		
		dir := math.Atan2(float64(x), -float64(y))

		force := math.Sqrt(float64(x*x + y*y))
		if force > 0.001 {
			force = math.Min(force, 0.0001)
		}

		physics.ApplyForce(float32(dir), float32(force))
	}
}
