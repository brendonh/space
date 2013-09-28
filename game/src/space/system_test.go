package space

import (
	"strconv"
	"math/rand"
	"reflect"
	"testing"
)

type TestComponent struct {
	BaseComponent
	Name string
}

func (c *TestComponent) Tag() string {
	return "test"
}

func (c *TestComponent) TickPhysics() {
}

func (c *TestComponent) String() string {
	return c.Name
}


func TestSystemAdd(t *testing.T) {
	s := &PhysicsSystem{}

	s.Add(&TestComponent{ Name: "one" })
	s.Add(&TestComponent{ Name: "two" })

	if len(s.Active) != 0 {
		t.Error("Not empty", s.Active)
	}

	s.Update()

	if len(s.Active) != 2 {
		t.Error("Wrong element count", s.Active)
	}
}


func TestSystemRemove(t *testing.T) {
	candidate := &TestComponent{ Name: "candidate" }
	s := &PhysicsSystem {
		Active: []PhysicsComponent {
			&TestComponent{ Name: "one" },
			&TestComponent{ Name: "candidate" },
			&TestComponent{ Name: "two" },
			candidate,
			&TestComponent{ Name: "three" },
		},
	}

	s.Remove(candidate)

	if len(s.Active) != 5 {
		t.Error("Wrong elements", s.Active)
	}

	s.Update()

	var names []string
	for _, c := range s.Active {
		names = append(names, c.(*TestComponent).Name)
	}

	if !reflect.DeepEqual(names, []string{ "one", "candidate", "two", "three" }) {
		t.Error("Wrong elements", s.Active)
	}
}


func BenchmarkSystemRemove(b *testing.B) {
	
	s := &PhysicsSystem{}
	var comps []PhysicsComponent

	for i := 0; i < b.N; i++ {
		comp := &TestComponent{ Name: strconv.Itoa(i) }
		s.Add(comp)
		comps = append(comps, comp)
	}

	for i := range comps {
		j := rand.Intn(i + 1)
		comps[i], comps[j] = comps[j], comps[i]
	}

	s.Update()

	b.ResetTimer()

	for _, comp := range comps {
		s.Remove(comp)
		s.Update()
	}

	if len(s.Active) != 0 {
		b.Error("Wrong elements:", s.Active)
	}
	
}
