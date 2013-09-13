package space

import (
	"fmt"
)

type ShipControl struct {
	Ship *Entity
}

func (c *ShipControl) Priority() int {
	return 1
}

func (c *ShipControl) Actions() []string {
	return []string {
		"ship_accel", "ship_decel", "ship_left", "ship_right",
	}
}

func (c *ShipControl) KeyDown(action string) bool {
	fmt.Println("Ship control on:", action)
	return true
}


func (c *ShipControl) KeyUp(action string) {
	fmt.Println("Ship control off:", action)
}
