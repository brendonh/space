package space

import (
	"container/heap"
	"math/rand"

	"testing"
)


func TestOpenTileSet(t *testing.T) {
	indexes := newTileIndexes(5, 25)
	ots := newOpenTileSet(Vec2i{ 5, 5 }, indexes)

	heap.Init(ots)

	p := Vec2i{ 0, 0 }

	ots.Add(Vec2i{ 0, 0 }, p, 0)
	ots.Add(Vec2i{ 4, 0 }, p, 0)
	ots.Add(Vec2i{ 0, 3 }, p, 0)
	ots.Add(Vec2i{ 2, 2 }, p, 0)
	ots.Add(Vec2i{ 1, 1 }, p, 0)

	found, ok := ots.Peek(Vec2i{ 1, 1 })
	if !ok || (found.Coords != Vec2i{ 1, 1 }) {
		t.Errorf("Couldn't find 1,1")
		return
	}

	out := heap.Remove(ots, found.Index)
	if out == nil || (out.(openTile).Coords != Vec2i{ 1, 1 }) {
		t.Errorf("Couldn't pop 1,1, got %v instead", out)
	}

	var testPop = func(expect Vec2i) {
		actual := heap.Pop(ots).(openTile).Coords
		if actual != expect {
			t.Errorf("Got %v, expected %v", actual, expect)
		}
	}

	testPop(Vec2i{ 2, 2 })
	testPop(Vec2i{ 4, 0 })
	testPop(Vec2i{ 0, 3 })
	testPop(Vec2i{ 0, 0 })

}

func BenchmarkOpenTileSet(b *testing.B) {
	var DIMENSION = 100

	coords := make([]Vec2i, 0, DIMENSION)
	for i := 0; i < DIMENSION; i++ {
		coords = append(coords, Vec2i{ rand.Intn(DIMENSION), rand.Intn(DIMENSION) })
	}

	indexes := newTileIndexes(DIMENSION, DIMENSION*DIMENSION)
	ots := newOpenTileSet(Vec2i{ 0, 0 }, indexes)
	heap.Init(ots)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, coord := range coords {
			ots.Add(coord, Vec2i{0, 0}, 0)
		}

		for i := 0; i < DIMENSION; i++ {
			heap.Pop(ots)
		}
	}

}