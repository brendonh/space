package space

import (
	. "github.com/brendonh/glvec"
)


type RoomsComponent struct {
	BaseComponent
	Rooms []*Room
	Grid TileGrid
	SelectedTile Tile
	Center Vec3
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
	var selection Tile

	var tile = r.Grid.Get(pos)
	if tile.Valid {
		selection = Tile{ 
			Pos: pos, 
			Valid: true,
		}
	}

	if !selection.Valid {
		r.ClearSelectedTile()
		return false
	} 

	if (selection != r.SelectedTile) {
		r.SelectedTile = selection
	}

	return true
}

func (r *RoomsComponent) ClearSelectedTile() {
	if r.SelectedTile.Valid {
		r.SelectedTile = Tile{}
		r.Entity.BroadcastEvent("selected_tile", r.SelectedTile.Pos)
	}
}

func (r *RoomsComponent) TriggerSelectedTile() {
	r.Entity.BroadcastEvent("trigger_tile", r.SelectedTile.Pos)
}

func (r *RoomsComponent) TileToModel(pos Vec2i) (worldPos Vec3) {
	worldPos = Vec3 { float32(pos.X), float32(pos.Y), 0 }
	V3Add(&worldPos, worldPos, r.Center)
	V3ScalarMul(&worldPos, worldPos, 2) // Oof
	return
}


func (r *RoomsComponent) update() {
	var cubes []Cube
	var cogX, cogY, mass float32

	// XXX TODO: Should this use the grid?

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

	r.Center = Vec3{ -cogX, -cogY, 0 }

	r.Entity.BroadcastEvent("update_cubes", &CubeSet{
		Cubes: cubes,
		Center: r.Center,
	})
}