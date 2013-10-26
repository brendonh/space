package space

import (
	"container/heap"
	"math"
)

// --------------------------------------------------------

type RoomTile struct {
	Pos Vec2i
	Color CubeColor
}

type Room struct {
	Pos Vec2i
	Tiles []RoomTile
}

func MakeSquareRoom(x, y, width, height int, color CubeColor) *Room {
	room := &Room{ Pos: Vec2i{ x, y } }

	var tiles = make([]RoomTile, 0, width * height)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			tiles = append(tiles, RoomTile{ Vec2i{ i, j }, color })
		}
	}
	room.Tiles = tiles
	return room
}


// --------------------------------------------------------


type Tile struct {
	Pos Vec2i
	Color CubeColor
	Valid bool
}

// --------------------------------------------------------

type TileGrid struct {
	Grid []Tile
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

func (g *TileGrid) Valid(pos Vec2i) bool {
	tile := g.Get(pos)
	return tile != nil && tile.Valid
}

func (g *TileGrid) Get(pos Vec2i) *Tile {
	idx := g.shipIndex(pos)
	if idx < 0 || idx >= len(g.Grid) {
		return nil
	}
	return &g.Grid[idx]
}

func (g *TileGrid) SetRooms(rooms []*Room) {
	var minX, maxX, minY, maxY int
	for _, room := range rooms {
		for _, roomTile := range room.Tiles {
			shipPos := room.Pos.Add(roomTile.Pos)
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

	g.Grid = make([]Tile, g.Extent.X * g.Extent.Y)
	for _, room := range rooms {
		for _, roomTile := range room.Tiles {
			roomPos := room.Pos.Add(roomTile.Pos)
			g.Grid[g.shipIndex(roomPos)] = Tile{
				Pos: roomPos,
				Color: roomTile.Color,
				Valid: true,
			}
		}
	}
}

func (g *TileGrid) FindPath(startTile, endTile Vec2i) ([]Vec2i, bool) {
	start := startTile.Add(g.Offset)
	end := endTile.Add(g.Offset)

	indexes := newTileIndexes(g.Extent.X, g.Extent.X * g.Extent.Y)
	open := newOpenTileSet(end, indexes)
	heap.Init(open)
	heap.Push(open, openTile{
		Coords: start,
		Cost: 0,
		Heuristic: open.DiagonalDistance(start),
		Valid: true,
	})

	closed := make([]openTile, (g.Extent.X * g.Extent.Y))

	neighbourBuf := make([]neighbourTile, 0, 8)

	var current openTile
	var success bool

	for open.Len() > 0 {
		current = heap.Pop(open).(openTile)

		if current.Coords == end {
			success = true
			break
		}

		closed[g.localIndex(current.Coords)] = current

		for _, neighbour := range g.neighbours(current.Coords, neighbourBuf) {
			cost := current.Cost + neighbour.MoveCost

			existingOpen, inOpen := open.Peek(neighbour.Coords)
			if inOpen && cost < existingOpen.Cost {
				heap.Remove(open, existingOpen.Index)
				inOpen = false
			}

			existingClosed := closed[g.localIndex(neighbour.Coords)]
			inClosed := existingClosed.Valid
			if inClosed && cost < existingClosed.Cost {
				closed[g.localIndex(neighbour.Coords)] = openTile{}
				inClosed = false
			}

			if !inOpen && !inClosed {
				open.Add(neighbour.Coords, current.Coords, cost)
			}
		}
	}

	if !success {
		return nil, false
	}

	reversePath := make([]Vec2i, 0, 10)
	for current.Coords != start {
		reversePath = append(reversePath, current.Coords)
		next := closed[g.localIndex(current.Parent)]
		if next == current {
			panic("BROKEN")
		}
		current = next
	}

	path := make([]Vec2i, 0, len(reversePath))
	negOffset := g.Offset.Neg()
	for i := len(reversePath) - 1; i >= 0; i-- {
		path = append(path, reversePath[i].Add(negOffset))
	}

	return path, true
}

func (g *TileGrid) localPosAvailable(coord Vec2i) bool {
	if coord.X < 0 || coord.Y < 0 || coord.X >= g.Extent.X || coord.Y >= g.Extent.Y {
		return false
	}
	idx := g.localIndex(coord)
	if idx >= len(g.Grid) {
		return false
	}

	tile := g.Grid[idx]

	return tile.Valid
}

type neighbourTile struct {
	Coords Vec2i
	MoveCost float64
}

func (g *TileGrid) neighbours(coords Vec2i, buf []neighbourTile) []neighbourTile {
	buf = buf[:0]

	push := func(neighbour Vec2i, cost float64) {
		buf = append(buf, neighbourTile{ neighbour, cost })
	}

	left  := coords.Add(Vec2i{ -1,  0 })
	right := coords.Add(Vec2i{  1,  0 })
	up    := coords.Add(Vec2i{  0,  1 })
	down  := coords.Add(Vec2i{  0, -1 })

	leftOK  := g.localPosAvailable(left)
	rightOK := g.localPosAvailable(right)
	upOK    := g.localPosAvailable(up)
	downOK  := g.localPosAvailable(down)

	if leftOK  { push(left, 1) }
	if rightOK { push(right, 1) }
	if upOK    { push(up, 1) }
	if downOK  { push(down, 1) }

	if leftOK && upOK {
		leftUp := coords.Add(Vec2i{ -1,  1 })
		if g.localPosAvailable(leftUp) { push(leftUp, DIAG_COST) }
	}

	if leftOK && downOK {
		leftDown := coords.Add(Vec2i{ -1, -1 })
		if g.localPosAvailable(leftDown) { push(leftDown, DIAG_COST) }
	}

	if rightOK && upOK {
		rightUp := coords.Add(Vec2i{  1,  1 })
		if g.localPosAvailable(rightUp) { push(rightUp, DIAG_COST) }
	}

	if rightOK && downOK {
		rightDown := coords.Add(Vec2i{  1, -1 })
		if g.localPosAvailable(rightDown) { push(rightDown, DIAG_COST) }
	}

	return buf
}


// -----------------------------------------

type tileIndexes struct {
	ExtentX int
	indexes []int
}

func newTileIndexes(extentX int, maxIndexes int) *tileIndexes {
	return &tileIndexes{
		ExtentX: extentX,
		indexes: make([]int, maxIndexes + 1),
	}
}

func (ti *tileIndexes) localIndex(pos Vec2i) int {
	return (pos.Y * ti.ExtentX) + pos.X
}

func (ti *tileIndexes) Set(coord Vec2i, i int) {
	ti.indexes[ti.localIndex(coord)] = i + 1
}

func (ti *tileIndexes) Get(coord Vec2i) int {
	idx := ti.indexes[ti.localIndex(coord)]
	if idx == 0 {
		return -1
	}
	return idx - 1
}

func (ti *tileIndexes) Clear(coord Vec2i) {
	ti.indexes[ti.localIndex(coord)] = 0
}

// ------------------------------------------

type openTile struct {
	Coords Vec2i
	Parent Vec2i
	Cost float64
	Heuristic float64
	Index int
	Valid bool
}

type openTileSet struct {
	Target Vec2i
	indexes *tileIndexes
	queue []openTile
}

func newOpenTileSet(target Vec2i, tileIndexes *tileIndexes) *openTileSet {
	return &openTileSet{
		Target: target,
		indexes: tileIndexes,
		queue: make([]openTile, 0, 10),
	}
}

func (ots *openTileSet) Add(coord Vec2i, parent Vec2i, cost float64) {
	heap.Push(ots, openTile{
		Coords: coord,
		Parent: parent,
		Cost: cost,
		Heuristic: ots.DiagonalDistance(coord),
		Valid: true,
	})
}

func (ots *openTileSet) Peek(coord Vec2i) (openTile, bool) {
	idx := ots.indexes.Get(coord)
	if idx == -1 {
		return openTile{}, false
	}
	return ots.queue[idx], true
}

func (ots *openTileSet) Len() int {
	return len(ots.queue)
}

func (ots *openTileSet) Less(i, j int) bool {
	ghi := ots.queue[i].Cost + ots.queue[i].Heuristic
	ghj := ots.queue[j].Cost + ots.queue[j].Heuristic
	return ghi < ghj
}

func (ots *openTileSet) Swap(i, j int) {
	ots.queue[i], ots.queue[j] = ots.queue[j], ots.queue[i]
	ots.queue[i].Index = i
	ots.queue[j].Index = j
	ots.indexes.Set(ots.queue[i].Coords, i)
	ots.indexes.Set(ots.queue[j].Coords, j)
}

func (ots *openTileSet) Push(x interface{}) {
	ot := x.(openTile)
	ot.Index = len(ots.queue)
	ots.queue = append(ots.queue, ot)
	ots.indexes.Set(ot.Coords, ot.Index)
}

func (ots *openTileSet) Pop() interface{} {
	idx := len(ots.queue) - 1
	ot := ots.queue[idx]
	ots.queue = ots.queue[:idx]
	ots.indexes.Clear(ot.Coords)
	ot.Index = -1
	return ot
}

var DIAG_COST = math.Sqrt(2)

func (ots *openTileSet) DiagonalDistance(pos Vec2i) float64 {
	dx := math.Abs(float64(ots.Target.X - pos.X))
    dy := math.Abs(float64(ots.Target.Y - pos.Y))
    return (dx + dy) + (DIAG_COST - 2) * math.Min(dx, dy)
}
