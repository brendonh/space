package main

import (
	"fmt"
	"runtime"

	"space"

	glfw "github.com/go-gl/glfw3"

)

func init() {
	runtime.LockOSThread()
}

const (
	Title  = "SPACE"
	Width  = 800
	Height = 600
	Fullscreen = false
)


func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}


func main() {

	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("Can't init glfw!")
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Samples, 2);

	var monitor *glfw.Monitor
	var width, height int
	if Fullscreen {
		monitor, err := glfw.GetPrimaryMonitor()
		if err != nil {
			panic(err)
		}
		
		mode, err := monitor.GetVideoMode()
		width = mode.Width
		height = mode.Height
	} else {
		width = Width
		height = Height
	}

	window, err := glfw.CreateWindow(width, height, Title, monitor, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	ml := space.NewMainloop(window)
	ml.RenderContext.Init()

	initSector(ml)

	ml.Loop()
	
}


func initSector(ml *space.Mainloop) {
	ml.Sector = space.NewSector()
	ship := ml.Entities.NewEntity()
	ship.AddComponent(&space.SpacePhysics{})
	ship.AddComponent(space.NewCubesComponent())
	ship.AddComponent(&space.ShipControl{})
	ship.AddComponent(&space.ShipInput{})
	ship.InitComponents()
	ml.Sector.AddEntity(ship)
	ml.RenderContext.FollowEntity(ship)

	
	for i := 1; i < 100; i++ {

		var ship2 = ml.Entities.NewEntity()
		ship2.AddComponent(&space.SpacePhysics{
			Position: space.SpacePosition {
				PosX: float64(i),
				PosY: float64(i),
			},
		})
		ship2.AddComponent(space.NewCubesComponent())
		ship2.AddComponent(&space.ShipControl{})
		ship2.InitComponents()
		ml.Sector.AddEntity(ship2)
	}

}

