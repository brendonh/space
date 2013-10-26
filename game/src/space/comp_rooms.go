package space

import (
	"math"

	. "github.com/brendonh/glvec"
)


type TileSelection struct {
	Pos Vec2i
	Valid bool
}

type RoomsComponent struct {
	BaseComponent
	Manager *ActionManager
	Rooms []*Room
	SelectedTile TileSelection
	Center Vec3
}

func (r *RoomsComponent) Tag() string {
	return "rooms"
}

func (r *RoomsComponent) Init() {
	r.Manager = r.Entity.GetComponent("action_manager").(*ActionManager)
	r.Refresh()
}

func (r *RoomsComponent) AddRoom(room *Room) {
	r.Rooms = append(r.Rooms, room)
}

func (r *RoomsComponent) Refresh() {
	r.Manager.Grid.SetRooms(r.Rooms)
	r.updateCubes()
}


func (r *RoomsComponent) SetSelectedTile(pos Vec2i) bool {
	var selection TileSelection

	if r.Manager.Grid.Valid(pos) {
		selection = TileSelection{ 
			Pos: pos, 
			Valid: true,
		}
	}

	if !selection.Valid {
		r.ClearSelectedTile()
		return false
	} 

	if selection != r.SelectedTile {
		r.SelectedTile = selection
	}

	return true
}

func (r *RoomsComponent) ClearSelectedTile() {
	if r.SelectedTile.Valid {
		r.SelectedTile = TileSelection{}
		r.Entity.BroadcastEvent("selected_tile", r.SelectedTile.Pos)
	}
}

func (r *RoomsComponent) TriggerSelectedTile() {
	manager := r.Entity.GetComponent("action_manager").(*ActionManager)
	manager.AddAction(&Action{
		Location: r.SelectedTile.Pos, 
	})
}

func (r *RoomsComponent) TileToModel(pos Vec2i) (worldPos Vec3) {
	worldPos = Vec3 { float32(pos.X), float32(pos.Y), 0 }
	V3Add(&worldPos, worldPos, r.Center)
	V3ScalarMul(&worldPos, worldPos, CUBE_SCALE) // Oof
	return
}

func (r *RoomsComponent) ModelToTile(worldPos Vec3) Vec2i {
	V3ScalarDiv(&worldPos, worldPos, CUBE_SCALE)
	
	V3Sub(&worldPos, worldPos, r.Center)
	
	return Vec2i{
		int(math.Floor(float64(worldPos[0]))),
		int(math.Floor(float64(worldPos[1]))),
	}
}


func (r *RoomsComponent) updateCubes() {
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