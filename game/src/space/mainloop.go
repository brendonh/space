package space

import (
	"fmt"

	"encoding/json"
	"io/ioutil"

	glfw "github.com/go-gl/glfw3"
)

type Mainloop struct {
	Window *glfw.Window
	Entities *EntityManager
	RenderContext *RenderContext
	Sector *Sector

	stopping bool
}

func NewMainloop(window *glfw.Window) *Mainloop {
	return &Mainloop {
		Window: window,
		RenderContext: NewRenderContext(),
		Entities: NewEntityManager(),
		stopping: false,
	}
}

func (m *Mainloop) SetSector(sector *Sector) {
	m.Sector = sector
}

func (m *Mainloop) Loop() {
	var ticksPerSecond float64 = 60.0
	secondsPerTick := 1 / ticksPerSecond
	fmt.Println("Seconds per tick:", secondsPerTick)

	width, height := m.Window.GetSize()
	m.RenderContext.Resize(width, height)

	m.Window.SetFramebufferSizeCallback(m.OnResize)
	m.Window.SetKeyCallback(m.OnKey)

	m.Sector.RegisterComponent(&GameControl{ Mainloop: m })

	prevTime := glfw.GetTime()
	var tickAcc float64 = 0 
	
	for !m.Window.ShouldClose() {

		glfw.PollEvents()

		if (m.stopping) {
			break
		}

		now := glfw.GetTime()

		tickAcc += (now - prevTime)

		prevTime = now

		for ; tickAcc >= secondsPerTick; tickAcc -= secondsPerTick {
			m.Sector.Tick()
		}

		var delta = tickAcc / secondsPerTick

		m.RenderContext.StartFrame(delta)
		m.Sector.Render(m.RenderContext)
		m.Window.SwapBuffers()

	}

}

func (m *Mainloop) OnKey(w *glfw.Window, key glfw.Key, scancode int, 
	action glfw.Action, mods glfw.ModifierKey) {
	
	m.Sector.Input.HandleKey(key, action, mods)
}

func (m *Mainloop) OnResize(w *glfw.Window, width int, height int) {
	m.RenderContext.Resize(width, height)
}

func (m *Mainloop) DumpData(points interface{}, filename string) {
	bytes, err := json.Marshal(points)
	if err != nil {
		fmt.Println("Couldn't encode points:", err)
		return
	}

	err = ioutil.WriteFile(filename, bytes, 0644)
}