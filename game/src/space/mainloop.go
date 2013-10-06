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
	context.Window.SetMouseButtonCallback(m.OnMouse)
}

func (m *Mainloop) SetSector(sector *Sector) {
	m.Sector = sector
}

func (m *Mainloop) Loop() {
	var ticksPerSecond float64 = 60.0
	secondsPerTick := 1.0 / ticksPerSecond

	var context = m.RenderContext

	m.Sector.RegisterComponent(&GameControl{ Mainloop: m })

	var prevTime = glfw.GetTime()
	var tickAcc float64 = secondsPerTick

	m.Sector.Tick()

	var fpsCounter = 0
	var fpsAcc = 0.0

	for !context.Window.ShouldClose() {
		now := glfw.GetTime()

		acc := (now - prevTime)
		tickAcc += acc
		fpsAcc += acc
		fpsCounter++

		if fpsAcc > 5 {
			var fps = float64(fpsCounter) / fpsAcc
			fmt.Println("FPS:", fps)
			fpsAcc = 0.0
			fpsCounter = 0
		}

		prevTime = now

		m.Sector.InputSystem.UpdateMouse()

		glfw.PollEvents()
		
		for ; tickAcc >= secondsPerTick; tickAcc -= secondsPerTick {
			m.Sector.Tick()
		}

		
		var alpha = float32(tickAcc / secondsPerTick)
		
		m.RenderContext.StartFrame()
		m.Camera.UpdateRenderContext(m.RenderContext, alpha)
		m.Sector.Render(m.RenderContext, alpha)
		m.RenderContext.FlushQueue()
	}

}

func (m *Mainloop) OnKey(w *glfw.Window, key glfw.Key, scancode int, 
	action glfw.Action, mods glfw.ModifierKey) {
	
	m.Sector.InputSystem.HandleKey(key, action, mods)
}

func (m *Mainloop) OnMouse(w *glfw.Window, button glfw.MouseButton, 
	action glfw.Action, mods glfw.ModifierKey) {
	m.Sector.InputSystem.HandleMouse(button, action, mods)
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