package space

import (
	"fmt"
	"encoding/json"
	"io/ioutil"

	"space/render"

	glfw "github.com/go-gl/glfw3"
)

type Mainloop struct {
	Entities *EntityManager
	Camera *Camera
	RenderContext *render.Context
	Sector *Sector

	restart bool
}

var mainloop *Mainloop

func NewMainloop() *Mainloop {
	loop := &Mainloop {
		RenderContext: render.NewContext(),
		Entities: NewEntityManager(),
		Camera: NewCamera(),

		restart: false,
	}

	loop.RenderContext.Init()
	loop.PrepareWindow()
	return loop
}

func (m *Mainloop) MakeGlobal() {
	mainloop = m
}

func (m *Mainloop) PrepareWindow() {
	var context = m.RenderContext
	context.Window.SetFramebufferSizeCallback(m.OnResize)
	context.Window.SetKeyCallback(m.OnKey)
}

func (m *Mainloop) SetSector(sector *Sector) {
	m.Sector = sector
}

func (m *Mainloop) Loop() {
	var ticksPerSecond float64 = 60.0
	secondsPerTick := 1.0 / ticksPerSecond

	var context = m.RenderContext

	m.Sector.RegisterComponent(&GameControl{ Mainloop: m })

	prevTime := glfw.GetTime()
	var tickAcc float64 = secondsPerTick

	m.Sector.Tick()

	for !context.Window.ShouldClose() {
		now := glfw.GetTime()
		tickAcc += (now - prevTime)
		prevTime = now
		
		glfw.PollEvents()

		mx, my := m.RenderContext.Window.GetCursorPosition()
		m.Sector.Input.TickCursor(mx, my)
		
		for ; tickAcc >= secondsPerTick; tickAcc -= secondsPerTick {
			m.Sector.Tick()
		}
		
		var alpha = tickAcc / secondsPerTick
		
		m.RenderContext.StartFrame()
		m.Camera.UpdateRenderContext(m.RenderContext, alpha)
		m.Sector.Render(m.RenderContext, alpha)
		m.RenderContext.FlushQueue()
	}

}

func (m *Mainloop) OnKey(w *glfw.Window, key glfw.Key, scancode int, 
	action glfw.Action, mods glfw.ModifierKey) {
	
	m.Sector.Input.HandleKey(key, action, mods)
}

func (m *Mainloop) OnResize(w *glfw.Window, width int, height int) {
	m.RenderContext.Resize(width, height)
}

func (m *Mainloop) ToggleFullscreen() {
	m.RenderContext.ToggleFullscreen()
	m.PrepareWindow()
	render.ClearShaderCache()
	globalDispatch.Fire("gl_init", nil)
}

func (m *Mainloop) DumpData(points interface{}, filename string) {
	bytes, err := json.Marshal(points)
	if err != nil {
		fmt.Println("Couldn't encode points:", err)
		return
	}

	err = ioutil.WriteFile(filename, bytes, 0644)
}