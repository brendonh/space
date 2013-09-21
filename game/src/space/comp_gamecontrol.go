package space

type GameControl struct {
	BaseComponent

	Mainloop *Mainloop
}

func (c *GameControl) Tag() string {
	return "struct_gamecontrol"
}

func (c *GameControl) Priority() int {
	return 0
}

func (c *GameControl) Actions() []string {
	return []string { "fullscreen", "quit_game" }
}

func (c *GameControl) KeyDown(action string) bool {
	switch (action) {
	case "quit_game":
		c.Mainloop.RenderContext.Window.SetShouldClose(true)
		return true
	case "fullscreen":
		c.Mainloop.ToggleFullscreen()
		return true
	}
	return false
}


func (c *GameControl) KeyUp(action string) {
}
