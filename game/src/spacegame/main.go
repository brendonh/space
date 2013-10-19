package main

import (
	"fmt"
	"runtime"

	"net/http"
	_ "net/http/pprof"

	"space"
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
	ship.Name = "ship"
	ml.Entities.NameEntity(ship)

	ship.AddComponent(&space.SpacePhysics{
		Position: space.SpacePosition {
		},
	})
	cubes := space.NewCubesComponent()
	ship.AddComponent(cubes)
	ship.AddComponent(&space.ShipControl{})
	ship.AddComponent(&space.ShipInput{})

	rooms := &space.RoomsComponent{}
	rooms.AddRoom(space.MakeSquareRoom(0, 0, 3, 5, space.CubeColor{ 1.0, 1.0, 0.3, 1.0 }))
	rooms.AddRoom(space.MakeSquareRoom(3, 0, 4, 2, space.CubeColor{ 1.0, 0.5, 1.0, 1.0 }))
	rooms.AddRoom(space.MakeSquareRoom(-4, 0, 4, 2, space.CubeColor{ 1.0, 1.0, 0.5, 1.0 }))
	ship.AddComponent(rooms)

	ship.InitComponents()
	ml.Sector.AddEntity(ship)
	ml.Camera.FollowEntity(ship)

	guy := ml.Entities.NewEntity()
	guy.Name = "guy"
	ml.Entities.NameEntity(guy)

	worldPos := cubes.TileToModel(2, 1)
	pos := &space.AvatarPosition{
		Position: space.SpacePosition {
			Pos: worldPos,
		},
	}
	guy.AddComponent(pos)
	guy.AddComponent(space.NewAvatarRenderer())
	guy.InitComponents()
	ml.Sector.AddEntity(guy)

	pos.AttachTo(ship)
}

