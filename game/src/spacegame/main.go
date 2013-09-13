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

	window, err := glfw.CreateWindow(Width, Height, Title, nil, nil)
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
}

