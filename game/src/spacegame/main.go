package main

import (
	"fmt"
	"runtime"
	"time"
	"math/rand"

	"net/http"
	_ "net/http/pprof"

	"space"
)

func init() {
	runtime.LockOSThread()
}



func main() {
	runtime.GOMAXPROCS(2)

	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	ml := space.NewMainloop()
	ml.MakeGlobal()

	initSector(ml)

	go func() {
		ship := ml.Entities.GetNamedEntity("ship")
		manager := ship.GetComponent("action_manager").(*space.ActionManager)
		c := time.Tick(1 * time.Second)
		for {
			
			ml.Interventions<- func() {
				if manager.HasIdlers() {
					grid := manager.Grid
					var pos space.Vec2i
					for {
						posX := rand.Intn(grid.Extent.X) - grid.Offset.X
						posY := rand.Intn(grid.Extent.Y) - grid.Offset.Y
						pos = space.Vec2i{ posX, posY }
						tile := grid.Get(pos)
						if tile != nil {
							break
						}
					}

					fmt.Println("Triggering:", pos)
					manager.AddAction(&space.Action{
						Location: pos,
					})
				}
			}
			<-c
		}
	}()

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

	ship.AddComponent(space.NewActionManager())

	rooms := &space.RoomsComponent{}
	rooms.AddRoom(space.MakeSquareRoom(0, 0, 3, 5, space.CubeColor{ 1.0, 1.0, 0.3, 1.0 }))
	rooms.AddRoom(space.MakeSquareRoom(3, 0, 4, 2, space.CubeColor{ 1.0, 0.5, 1.0, 1.0 }))
	rooms.AddRoom(space.MakeSquareRoom(-4, 0, 4, 2, space.CubeColor{ 1.0, 1.0, 0.5, 1.0 }))
	rooms.AddRoom(space.MakeSquareRoom(3, 5, 4, 2, space.CubeColor{ 1.0, 1.0, 0.5, 1.0 }))
	rooms.AddRoom(space.MakeSquareRoom(0, 5, 1, 5, space.CubeColor{ 1.0, 1.0, 0.5, 1.0 }))
	rooms.AddRoom(space.MakeSquareRoom(1, 9, 5, 1, space.CubeColor{ 1.0, 1.0, 0.5, 1.0 }))
	rooms.AddRoom(space.MakeSquareRoom(4, 7, 2, 2, space.CubeColor{ 1.0, 1.0, 0.5, 1.0 }))
	ship.AddComponent(rooms)

	ship.InitComponents()

	ml.Sector.AddEntity(ship)
	ml.Camera.FollowEntity(ship)

	guy := ml.Entities.NewEntity()
	guy.Name = "guy"
	ml.Entities.NameEntity(guy)

	pos := &space.AvatarPosition{
		WalkSpeed: 1.0 / 10.0,
	}
	guy.AddComponent(pos)
	guy.AddComponent(space.NewAvatarBehaviour())
	guy.AddComponent(space.NewAvatarRenderer(space.CubeColor{ 1.0, 0.3, 0.3, 1.0 }))
	guy.InitComponents()
	ml.Sector.AddEntity(guy)
	pos.AttachToShipPosition(ship, space.Vec2i{ 0, 1 })


	guy = ml.Entities.NewEntity()
	guy.Name = "guy2"
	ml.Entities.NameEntity(guy)

	pos = &space.AvatarPosition{
		WalkSpeed: 1.0 / 20.0,
	}
	guy.AddComponent(pos)
	guy.AddComponent(space.NewAvatarBehaviour())
	guy.AddComponent(space.NewAvatarRenderer(space.CubeColor{ 0.3, 1.0, 0.3, 1.0 }))
	guy.InitComponents()
	ml.Sector.AddEntity(guy)
	pos.AttachToShipPosition(ship, space.Vec2i{ 0, 2 })



	guy = ml.Entities.NewEntity()
	guy.Name = "guy3"
	ml.Entities.NameEntity(guy)

	pos = &space.AvatarPosition{
		WalkSpeed: 1.0 / 15.0,
	}
	guy.AddComponent(pos)
	guy.AddComponent(space.NewAvatarBehaviour())
	guy.AddComponent(space.NewAvatarRenderer(space.CubeColor{ 0.3, 0.3, 1.0, 1.0 }))
	guy.InitComponents()
	ml.Sector.AddEntity(guy)
	pos.AttachToShipPosition(ship, space.Vec2i{ -1, 0 })

}

