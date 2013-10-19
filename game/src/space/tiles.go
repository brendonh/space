package space

type Tile struct {
	Room *Room
	X, Y int
	Color CubeColor
}

func (t *Tile) GetShipPos() (x, y int) {
	x = t.Room.X + t.X
	y = t.Room.Y + t.Y
	return
}


// --------------------------------------------------------


type Room struct {
	X, Y int
	Tiles []*Tile
}

func MakeSquareRoom(x, y, width, height int, color CubeColor) *Room {
	room := &Room{
		X: x, 
		Y: y,
	}

	var tiles = make([]*Tile, 0, width * height)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			tiles = append(tiles, &Tile{ room, i, j, color })
		}
	}
	room.Tiles = tiles
	return room
}


// --------------------------------------------------------


type TileGrid struct {
	Grid []*Tile
	ExtentX, ExtentY int
	OffsetX, OffsetY int
}

func (g *TileGrid) idx(x, y int) int {
	x += g.OffsetX
	y += g.OffsetY
	return (y * g.ExtentX) + x
}

func (g *TileGrid) Get(x, y int) *Tile {
	idx := g.idx(x, y)
	if idx < 0 || idx > len(g.Grid) {
		return nil
	}
	return g.Grid[idx]
}

func (g *TileGrid) Set(x, y int, t *Tile) {
	idx := g.idx(x, y)
	g.Grid[idx] = t
}

func (g *TileGrid) SetRooms(rooms []*Room) {
	var minX, maxX, minY, maxY int
	for _, room := range rooms {
		for _, tile := range room.Tiles {
			x, y := tile.GetShipPos()
			if x < minX { minX = x }
			if x > maxX { maxX = x }
			if y < minY { minY = y }
			if y > maxY { maxY = y }
		}
	}
	g.ExtentX = (maxX - minX) + 1
	g.ExtentY = (maxY - minY) + 1

	g.OffsetX = -minX
	g.OffsetY = -minY	

	g.Grid = make([]*Tile, g.ExtentX * g.ExtentY)
	for _, room := range rooms {
		for _, tile := range room.Tiles {
			x, y := tile.GetShipPos()
			g.Set(x, y, tile)
		}
	}
}
