package space

import (
	"fmt"
)

type CrewComponent struct {
	BaseComponent
	Avatars []*Entity
}

func (c *CrewComponent) Tag() string {
	return "crew"
}

func (c *CrewComponent) Add(avatar *Entity) {
	c.Avatars = append(c.Avatars, avatar)
}

func (c *CrewComponent) Remove(avatar *Entity) {
	for i, av := range c.Avatars {
		if av == avatar {
			var a = c.Avatars
			a[len(a)-1], a[i], a = nil, a[len(a)-1], a[:len(a)-1]
			c.Avatars = a
			break
		}
	}
}

func (c *CrewComponent) Event(tag string, args interface{}) {
	switch(tag) {
	case "trigger_tile":
		var target = args.(*Tile)
		for _, crew := range c.Avatars {
			crew.BroadcastEvent("move_to", target)
			break
		}
	}
}