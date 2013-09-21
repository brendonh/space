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
	
	// for i := 1; i < 100; i++ {

	// 	var ship2 = ml.Entities.NewEntity()
	// 	ship2.AddComponent(&space.SpacePhysics{
	// 		Position: space.SpacePosition {
	// 			PosX: float64(i),
	// 			PosY: float64(i),
	// 		},
	// 	})
	// 	ship2.AddComponent(space.NewCubesComponent())
	// 	ship2.AddComponent(&space.ShipControl{})
	// 	ship2.InitComponents()
	// 	ml.Sector.AddEntity(ship2)
	// }

}

