package space

import (
	"fmt"
	"encoding/json"
	"io/ioutil"

	"space/render"

	glfw "github.com/go-gl/glfw3"
)

type Mainloop struct {
	Window *glfw.Window
	Entities *EntityManager

	Camera *Camera

	RenderContext *render.Context

	Sector *Sector

	stopping bool
}

func NewMainloop(window *glfw.Window) *Mainloop {
	return &Mainloop {
		Window: window,
		RenderContext: render.NewContext(),
		Entities: NewEntityManager(),
		Camera: NewCamera(),
		stopping: false,
	}
}

func (m *Mainloop) SetSector(sector *Sector) {
	m.Sector = sector
}

func (m *Mainloop) Loop() {
	var ticksPerSecond float64 = 60.0
	secondsPerTick := 1.0 / ticksPerSecond

	width, height := m.Window.GetSize()
	m.RenderContext.Resize(width, height)

	m.Window.SetFramebufferSizeCallback(m.OnResize)
	m.Window.SetKeyCallback(m.OnKey)

	m.Sector.RegisterComponent(&GameControl{ Mainloop: m })

	prevTime := glfw.GetTime()
	var tickAcc float64 = secondsPerTick

	var tickTimes []float64
	var frameTimes []float64
	m.Sector.Tick()
	
	for !m.Window.ShouldClose() {
		now := glfw.GetTime()
		tickAcc += (now - prevTime)
		prevTime = now

		glfw.PollEvents()

		if (m.stopping) {
			break
		}

		for ; tickAcc >= secondsPerTick; tickAcc -= secondsPerTick {
			m.Sector.Tick()
		}

		var alpha = tickAcc / secondsPerTick

		m.RenderContext.StartFrame()
		m.Camera.UpdateRenderContext(m.RenderContext, alpha)
		m.Sector.Render(m.RenderContext, alpha)
		m.RenderContext.FlushQueue()

		tickTimes = append(tickTimes, glfw.GetTime() - now)

		m.Window.SwapBuffers()

		frameTimes = append(frameTimes, glfw.GetTime() - now)
	}

	m.DumpData(tickTimes, "ticks.json")
	m.DumpData(frameTimes, "frametimes.json")

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