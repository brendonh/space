package main

import (
	"fmt"
	"runtime"

	"net/http"
	_ "net/http/pprof"

	"space"

	. "github.com/brendonh/glvec"
)

func init() {
	runtime.LockOSThread()
}



func main() {

	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	ml := space.NewMainloop()
	ml.MakeGlobal()

	initSector(ml)

	ml.Loop()
	
}


func initSector(ml *space.Mainloop) {
	ml.Sector = space.NewSector()

	ship := ml.Entities.NewEntity()
	ship.AddComponent(&space.SpacePhysics{
		Position: space.SpacePosition {
			PosX: 0.0,
			PosY: 0.0,
		},
	})
	ship.AddComponent(space.NewCubesComponent())
	ship.AddComponent(&space.ShipControl{})
	ship.AddComponent(&space.ShipInput{})

	rooms := &space.RoomsComponent{}
	rooms.AddRoom(space.MakeSquareRoom(3, 5, space.CubeColor{ 1.0, 1.0, 0.3 }))
	ship.AddComponent(rooms)

	ship.InitComponents()
	ml.Sector.AddEntity(ship)
	ml.Camera.FollowEntity(ship)


	ship = ml.Entities.NewEntity()
	ship.AddComponent(&space.SpacePhysics{
		Position: space.SpacePosition {
			PosX: 10.0,
			PosY: 10.0,
		},
	})
	cubes := space.NewCubesComponent()
	M4MakeScale(&cubes.MModel, 0.1)
	M4SetTransform(&cubes.MModel, Vec3 { 0.0, 0.0, 1.0 })
	cubes.ShowEdges = false
	ship.AddComponent(cubes)
	rooms = &space.RoomsComponent{}
	rooms.AddRoom(space.MakeSquareRoom(1, 1, space.CubeColor{ 1.0, 0.0, 0.0 }))
	ship.AddComponent(rooms)
	ship.InitComponents()
	ml.Sector.AddEntity(ship)

}

