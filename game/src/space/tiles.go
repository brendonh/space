package space

import (
	"fmt"
)

type Vec2i struct {
	X, Y int
}

type Tile struct {
	Room *Room
	Pos Vec2i
	Color CubeColor
}

func (t *Tile) GetShipPos() Vec2i {
	return Vec2i{
		t.Room.Pos.X + t.Pos.X,
		t.Room.Pos.Y + t.Pos.Y,
	}
}


// --------------------------------------------------------


type Room struct {
	Pos Vec2i
	Tiles []*Tile
}

func MakeSquareRoom(x, y, width, height int, color CubeColor) *Room {
	room := &Room{
		Pos: Vec2i{
			X: x, 
			Y: y,
		},
	}

	var tiles = make([]*Tile, 0, width * height)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			tiles = append(tiles, &Tile{ room, Vec2i{ i, j }, color })
		}
	}
	room.Tiles = tiles
	return room
}


// --------------------------------------------------------


type TileGrid struct {
	Grid []*Tile
	Extent Vec2i
	Offset Vec2i
}

func (g *TileGrid) localIndex(pos Vec2i) int {
	return (pos.Y * g.Extent.X) + pos.X
}

func (g *TileGrid) shipIndex(pos Vec2i) int {
	pos.X += g.Offset.X
	pos.Y += g.Offset.Y
	if pos.X < 0 || pos.Y < 0 {
		return -1
	}
	return g.localIndex(pos)
}

func (g *TileGrid) Get(pos Vec2i) *Tile {
	idx := g.shipIndex(pos)
	if idx < 0 || idx >= len(g.Grid) {
		return nil
	}
	return g.Grid[idx]
}

func (g *TileGrid) Set(pos Vec2i, t *Tile) {
	idx := g.shipIndex(pos)
	g.Grid[idx] = t
}

func (g *TileGrid) SetRooms(rooms []*Room) {
	var minX, maxX, minY, maxY int
	for _, room := range rooms {
		for _, tile := range room.Tiles {
			shipPos := tile.GetShipPos()
			if shipPos.X < minX { minX = shipPos.X }
			if shipPos.X > maxX { maxX = shipPos.X }
			if shipPos.Y < minY { minY = shipPos.Y }
			if shipPos.Y > maxY { maxY = shipPos.Y }
		}
	}
	g.Extent.X = (maxX - minX) + 1
	g.Extent.Y = (maxY - minY) + 1

	g.Offset.X = -minX
	g.Offset.Y = -minY	

	g.Grid = make([]*Tile, g.Extent.X * g.Extent.Y)
	for _, room := range rooms {
		for _, tile := range room.Tiles {
			g.Set(tile.GetShipPos(), tile)
		}
	}
}

func (g *TileGrid) FindPath(startTile, endTile *Tile) {
	start := startTile.GetShipPos()
	end := endTile.GetShipPos()
	fmt.Println("Pathing", start, end)
}

//func (g *TileGrid) DiagonalDistance(

