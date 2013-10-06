package space

import (
	. "github.com/brendonh/glvec"
)

type Tile struct {
	X, Y int
	Color CubeColor
}

type Room struct {
	X, Y int
	Tiles []Tile
}

func MakeSquareRoom(width, height int, color CubeColor) *Room {
	var tiles = make([]Tile, 0, width * height)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			tiles = append(tiles, Tile{ i, j, color })
		}
	}
	return &Room {
		Tiles: tiles,
	}
}


type RoomsComponent struct {
	BaseComponent
	Rooms []*Room
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
}


func (r *RoomsComponent) SetSelectedTile(x, y int) bool {
	// TODO: Something other than this	
	for _, room := range r.Rooms {
		for i := range room.Tiles {
			tile := &room.Tiles[i]
			if tile.X == x && tile.Y == y {
				r.SelectedTile = tile
				return true
			}
		}
	}
	r.ClearSelectedTile()
	return false
}

func (r *RoomsComponent) ClearSelectedTile() {
	r.SelectedTile = nil
}

func (r *RoomsComponent) update() {
	var cubes []Cube
	var cogX, cogY, mass float32

	for _, room := range r.Rooms {
		for _, tile := range room.Tiles {
			cubes = append(cubes, Cube{ 
				X: tile.X, 
				Y: tile.Y, 
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