package space

import (
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


func (r *RoomsComponent) SetSelectedTile(pos Vec2i) bool {
	var selection *Tile

	var tile = r.Grid.Get(pos)
	if tile != nil {
		selection = tile
	}

	if selection == nil {
		r.ClearSelectedTile()
		return false
	} 

	if (selection != r.SelectedTile) {
		r.SelectedTile = selection
		shipPos := tile.GetShipPos()
		r.Entity.BroadcastEvent("selected_tile", r.SelectedTile)
		r.Entity.BroadcastEvent("update_colors", []CubeColorOverride {
			CubeColorOverride{ shipPos.X, shipPos.Y, CubeColor{ 1.0, 0.0, 0.0, 0.5 } },
		})
	}

	return true
}

func (r *RoomsComponent) ClearSelectedTile() {
	if r.SelectedTile != nil {
		r.SelectedTile = nil
		r.Entity.BroadcastEvent("selected_tile", r.SelectedTile)
		r.Entity.BroadcastEvent("update_colors", []CubeColorOverride{})
	}
}


func (r *RoomsComponent) update() {
	var cubes []Cube
	var cogX, cogY, mass float32

	for _, room := range r.Rooms {
		for _, tile := range room.Tiles {
			cubes = append(cubes, Cube{ 
				Pos: Vec2i{
					room.Pos.X + tile.Pos.X, 
					room.Pos.Y + tile.Pos.Y, 
				},
				Color: tile.Color,
			})

			cogX += float32(tile.Pos.X)
			cogY += float32(tile.Pos.Y)
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