package space

import (
	"fmt"
	. "github.com/brendonh/glvec"
)


type RoomsComponent struct {
	BaseComponent
	Rooms []*Room
	Grid TileGrid
	SelectedTile *Tile
}

func (r *RoomsComponent) Tag() string {
	return "rooms"
}

func (r *RoomsComponent) Init() {
	r.update()
}

func (r *RoomsComponent) AddRoom(room *Room) {
	r.Rooms = append(r.Rooms, room)
	r.Grid.SetRooms(r.Rooms)
}


func (r *RoomsComponent) SetSelectedTile(x, y int) bool {
	var selection *Tile

	var tile = r.Grid.Get(x, y)
	if tile != nil {
		selection = tile
	}

	if (selection != r.SelectedTile) {
		fmt.Println("Selection update", x, y, selection)
		r.SelectedTile = selection
	}
	return selection != nil
}

func (r *RoomsComponent) ClearSelectedTile() {
	if r.SelectedTile != nil {
		fmt.Println("Selection clear")
		r.SelectedTile = nil
	}
}


func (r *RoomsComponent) update() {
	var cubes []Cube
	var cogX, cogY, mass float32

	for _, room := range r.Rooms {
		for _, tile := range room.Tiles {
			cubes = append(cubes, Cube{ 
				X: room.X + tile.X, 
				Y: room.Y + tile.Y, 
				Color: tile.Color,
			})

			cogX += float32(tile.X)
			cogY += float32(tile.Y)
			mass += 1.0
		}
	}

	cogX /= mass
	cogY /= mass

	r.Entity.BroadcastEvent("update_cubes", &CubeSet{
		Cubes: cubes,
		Center: Vec3{ -cogX, -cogY, 0 },
	})
}