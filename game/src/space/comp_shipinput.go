package space

import (
	. "github.com/brendonh/glvec"
)

type ShipInput struct {
	BaseComponent

	ShipControl *ShipControl
}

func (c *ShipInput) Init() {
	c.ShipControl = c.Entity.GetComponent("struct_shipcontrol").(*ShipControl)
}

func (c *ShipInput) Tag() string {
	return "struct_shipinput"
}

func (c *ShipInput) Priority() int {
	return 1
}

func (c *ShipInput) Actions() []string {
	return []string {
		"ship_accel", "ship_decel", "ship_left", "ship_right", 
		"ship_debug_dump",
	}
}

func (c *ShipInput) KeyDown(action string) bool {
	return c.setState(action, 1)
}

func (c *ShipInput) KeyUp(action string) {
	c.setState(action, -1)
}

func (c *ShipInput) setState(action string, onOff float64) bool {
	switch action {
	case "ship_accel":
		c.ShipControl.Thrust += onOff
	case "ship_decel":
		c.ShipControl.Brake += onOff
	case "ship_left":
		c.ShipControl.Turn += onOff
	case "ship_right":
		c.ShipControl.Turn -= onOff
	case "ship_debug_dump":
		// Nada
	default:
		return false
	}
	return true
}


func (c *ShipInput) HandleCursorPosition(x, y float64) bool {
	var ray = mainloop.RenderContext.ScreenToWorld(x, y)
	
	pos, ok := ray.PlaneIntersect(Plane{
		Point: Vec3{ 0, 0, 1 },
		Normal: Vec3{ 0, 0, 1 },
	})

	if !ok {
		// XXX TODO No way to hide things right now
		pos = Vec3{ 0, 0, 10000 }
	}

	var indicator = mainloop.Sector.Entities[2]
	var physics = indicator.GetComponent("struct_spacephysics").(*SpacePhysics)
	physics.Position.PosX = float64(pos[0])
	physics.Position.PosY = float64(pos[1])

	return true
}
