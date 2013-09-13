package space

import "fmt"

type GameControl struct {
	Mainloop *Mainloop
}

func (c *GameControl) Init() {
}

func (c *GameControl) Tag() string {
	return ""
}

func (c *GameControl) SetEntity(e *Entity) {
}

func (c *GameControl) Priority() int {
	return 0
}

func (c *GameControl) Actions() []string {
	return []string { "quit_game" }
}

func (c *GameControl) KeyDown(action string) bool {
	switch (action) {
	case "quit_game":
		fmt.Println("Quitting")
		c.Mainloop.stopping = true
		return true
	}
	return false
}


func (c *GameControl) KeyUp(action string) {
}
