package space

import (
	"fmt"

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
		if onOff > 0 {
			c.DebugDump(true)
		}
	default:
		return false
	}
	return true
}


func (c *ShipInput) HandleCursorPosition(x, y float64) bool {
	c.DebugDump(false)
	return true
}


func (c *ShipInput) DebugDump(print bool) bool {
	var context = mainloop.RenderContext
	width, height := context.Window.GetSize()
	x, y := context.Window.GetCursorPosition()

	var pInv Mat4
	M4Inverse(&pInv, &context.MPerspective)

	var vInv Mat4
	M4Inverse(&vInv, &context.MView)

	ndx := (float32(x) / float32(width)) * 2.0 - 1.0
	ndy := 1.0 - (float32(y) / float32(height)) * 2.0

	var ray = Vec4 { ndx, ndy, 0.0, 0.0 }
	M4MulV4(&ray, &pInv, ray)

	if print {
		fmt.Println("~~~~~~~~~~~~~~~~~~~~")
		fmt.Println("ND", ndx, ndy)
		fmt.Println("VRay", ray)

		var back Vec4
		M4MulV4(&back, &context.MPerspective, ray)
		fmt.Println("Back to ND", back)
	}

	var rayDir = Vec3 { ray[0], ray[1], ray[2] }

	var camRotate Mat3
	M4RotationMatrix(&camRotate, &vInv)
	M3MulV3(&rayDir, &camRotate, rayDir)
	V3Normalize(&rayDir, rayDir)

	if print {
		fmt.Println("Rotated ray", rayDir)
	}

	var rayOrig = context.VCamPos

	if print {
		fmt.Println("Ray", rayOrig, rayDir)	
	}

	var camToPlane = rayOrig
	V3Sub(&camToPlane, camToPlane, Vec3{ 0.0, 0.0, 1.0 })

	var normal = Vec3 { 0.0, 0.0, 1.0 }

	var dist = -(V3Dot(camToPlane, normal)) / V3Dot(rayDir, normal)
	
	if print {
		fmt.Println("Dist", dist)
	}

	var wPos Vec3
	V3ScalarMul(&wPos, rayDir, dist)
	V3Add(&wPos, wPos, rayOrig)

	if print {
		fmt.Println("World", wPos)
	}

	var indicator = mainloop.Sector.Entities[2]
	var physics = indicator.GetComponent("struct_spacephysics").(*SpacePhysics)
	physics.Position.PosX = float64(wPos[0])
	physics.Position.PosY = float64(wPos[1])
	
	var prev = wPos
	for i := 0; i < 10; i++ {
		V3Add(&prev, prev, rayDir)
	}

	return true
}