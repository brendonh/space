package space

import (
	"fmt"
	"time"

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
	var ticksPerSecond int64 = 120
	nanosPerTick := int64(time.Second / time.Nanosecond) / ticksPerSecond
	fmt.Println("Nanos per tick:", nanosPerTick)

	prevTime := time.Now()
	var tickAcc int64 = 0 

	width, height := m.Window.GetSize()
	m.RenderContext.Resize(width, height)

	m.Window.SetFramebufferSizeCallback(m.OnResize)
	m.Window.SetKeyCallback(m.OnKey)

	m.Sector.RegisterComponent(&GameControl{ Mainloop: m })

	for !m.Window.ShouldClose() {
		glfw.PollEvents()

		if (m.stopping) {
			break
		}

		now := time.Now()
		tickAcc += int64(now.Sub(prevTime))
		prevTime = now

		for ; tickAcc > nanosPerTick; tickAcc -= nanosPerTick {
			m.Sector.Tick()
		}

		m.RenderContext.StartFrame()
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
