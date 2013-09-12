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
}

func NewMainloop(window *glfw.Window) *Mainloop {
	return &Mainloop {
		Window: window,
		RenderContext: NewRenderContext(),
		Entities: NewEntityManager(),
		Sector: NewSector(),
	}
}

func (m *Mainloop) Loop() {
	var ticksPerSecond int64 = 3
	nanosPerTick := int64(time.Second / time.Nanosecond) / ticksPerSecond
	fmt.Println("Nanos per tick:", nanosPerTick)

	prevTime := time.Now()
	var tickAcc int64 = 0 

	for !m.Window.ShouldClose() {
		glfw.PollEvents()

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
