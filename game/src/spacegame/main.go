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
			PosX: 0.0,
			PosY: 0.0,
		},
	})
	cubes := space.NewCubesComponent()
	ship.AddComponent(cubes)
	ship.AddComponent(&space.ShipControl{})
	ship.AddComponent(&space.ShipInput{})

	rooms := &space.RoomsComponent{}
	rooms.AddRoom(space.MakeSquareRoom(3, 5, space.CubeColor{ 1.0, 1.0, 0.3 }))
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
			PosX: float64(worldPos[0]),
			PosY: float64(worldPos[1]),
		},
	}
	guy.AddComponent(pos)
	guy.AddComponent(space.NewAvatarRenderer())
	guy.InitComponents()
	ml.Sector.AddEntity(guy)

	pos.AttachTo(ship)
}

