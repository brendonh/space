package space

import (
	"math"
)

type ShipControl struct {
	Ship *Entity
	Physics *SpacePhysics
	Thrust float64
	Brake float64
	Turn float64
}

func (c *ShipControl) Init() {
	c.Physics = c.Ship.GetComponent("struct_spacephysics").(*SpacePhysics)
}

func (c *ShipControl) Tag() string {
	return "struct_shipcontrol"
}

func (c *ShipControl) SetEntity(e *Entity) {
	c.Ship = e
}

func (c *ShipControl) TickLogic() {
	var physics = c.Physics

	physics.Angle += c.Turn * 0.05
	physics.ApplyForce(
		physics.Angle,
		c.Thrust * 0.003,
	)
	
	if (c.Brake > 0.005) {
		x, y := physics.VelX, physics.VelY
		
		dir := math.Atan2(x, -y)

		force := math.Sqrt(x*x + y*y)
		if force > 0.001 {
			force = math.Min(force, 0.0005)
		}

		physics.ApplyForce(dir, force)
	}		
}
