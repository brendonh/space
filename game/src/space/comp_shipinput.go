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
	var worldPos = c.ScreenToWorld(x, y)

	var indicator = mainloop.Sector.Entities[2]
	var physics = indicator.GetComponent("struct_spacephysics").(*SpacePhysics)
	physics.Position.PosX = float64(worldPos[0])
	physics.Position.PosY = float64(worldPos[1])

	return true
}


func (c *ShipInput) ScreenToWorld(x, y float64) Vec3 {
	var context = mainloop.RenderContext
	width, height := context.Window.GetSize()

	ndx := (float32(x) / float32(width)) * 2.0 - 1.0
	ndy := 1.0 - (float32(y) / float32(height)) * 2.0

	var ray = Vec4 { ndx, ndy, 0.0, 0.0 }
	M4MulV4(&ray, &context.MPerspectiveInverse, ray)

	var rayDir = Vec3 { ray[0], ray[1], ray[2] }

	var camRotateInverse Mat3
	M3Inverse(&camRotateInverse, &context.MCamRotate)
	M3MulV3(&rayDir, &camRotateInverse, rayDir)
	V3Normalize(&rayDir, rayDir)

	var rayOrig = context.VCamPos
	var planePoint = Vec3{ 0.0, 0.0, 1.0 }
	var normal = Vec3 { 0.0, 0.0, 1.0 }

	var camToPlane Vec3
	V3Sub(&camToPlane, rayOrig, planePoint)

	var dist = -(V3Dot(camToPlane, normal)) / V3Dot(rayDir, normal)

	var wPos Vec3
	V3ScalarMul(&wPos, rayDir, dist)
	V3Add(&wPos, wPos, rayOrig)

	return wPos
}