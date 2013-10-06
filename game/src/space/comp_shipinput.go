package space

import (
	"fmt"
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
		"mouse_activate",
		"ship_debug_dump",
	}
}

func (c *ShipInput) KeyDown(action string) bool {
	switch action {
	case "mouse_activate":
		fmt.Println("Click!")

	case "ship_debug_dump":
		guy := mainloop.Entities.GetNamedEntity("guy")
		pos := guy.GetComponent("struct_avatarposition").(*AvatarPosition)
		if pos.Attached() {
			pos.Detach()
		} else {
			ship := mainloop.Entities.GetNamedEntity("ship")
			pos.AttachTo(ship)
		}
	default:
		return c.setState(action, 1)
	}

	return true
}

func (c *ShipInput) KeyUp(action string) {
	c.setState(action, -1)
}

func (c *ShipInput) setState(action string, onOff float32) bool {
	switch action {
	case "ship_accel":
		c.ShipControl.Thrust += onOff
	case "ship_decel":
		c.ShipControl.Brake += onOff
	case "ship_left":
		c.ShipControl.Turn += onOff
	case "ship_right":
		c.ShipControl.Turn -= onOff
	default:
		return false
	}
	return true
}