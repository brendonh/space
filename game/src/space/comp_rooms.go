package space

import (
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
	Entity *Entity
	Rooms []*Room
	Cubes *CubesComponent
}

func (r *RoomsComponent) Tag() string {
	return "rooms"
}

func (r *RoomsComponent) SetEntity(e *Entity) {
	r.Entity = e
}

func (r *RoomsComponent) Init() {
	r.Cubes = r.Entity.GetComponent("cubes").(*CubesComponent)
	r.Update()
}

func (r *RoomsComponent) AddRoom(room *Room) {
	r.Rooms = append(r.Rooms, room)
}

func (r *RoomsComponent) Update() {
	var cubes []Cube
	var cogX, cogY, mass int

	for _, room := range r.Rooms {
		for _, tile := range room.Tiles {
			cubes = append(cubes, Cube{ 
				tile.X, tile.Y, tile.Color })

			cogX += tile.X
			cogY += tile.Y
			mass += 1
		}
	}

	cogX /= mass
	cogY /= mass

	for i := range cubes {
		cubes[i].X -= cogX
		cubes[i].Y -= cogY
	}

	r.Cubes.SetCubes(cubes)
}