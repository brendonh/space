package space

import (
	"fmt"
	"math"
)

type EntityID uint64

type Entity struct {
	ID EntityID
	Name string
	Sector *Sector
	Components []Component
}

func NewEntity(id EntityID) *Entity {
	return &Entity{
		ID: id,
	}
}

func (e *Entity) AddComponent(c Component) {
	c.SetEntity(e)
	e.Components = append(e.Components, c)
	if e.Sector != nil {
		e.Sector.RegisterComponent(c)
	}
}

func (e *Entity) GetComponent(tag string) Component {
	var component = e.FindComponent(tag)
	if component == nil {
		panic(fmt.Sprintf("No such component: %s", tag))
	}
	return component
}

func (e *Entity) FindComponent(tag string) Component {
	for _, c := range e.Components {
		if c.Tag() == tag {
			return c
		}
	}
	return nil
}

func (e *Entity) RemoveComponent(comp Component) {
	for i, c := range e.Components {
		if c == comp {
			e.Sector.UnregisterComponent(comp)
			var cs = e.Components
			cs[len(cs)-1], cs[i], cs = nil, cs[len(cs)-1], cs[:len(cs)-1]
			e.Components = cs
		}
	}
}

func (e *Entity) InitComponents() {
	for _, c := range e.Components {
		c.Init()
	}
}

func (e *Entity) BroadcastEvent(tag string, args interface{}) {
	for _, c := range e.Components {
		c.Event(tag, args)
	}
}

// ------------------------------------------

type EntityManager struct {
	Entities map[EntityID]*Entity
	NamedEntities map[string]EntityID
	lastID EntityID
}

func NewEntityManager() *EntityManager {
	return &EntityManager {
		Entities: make(map[EntityID]*Entity),
		NamedEntities: make(map[string]EntityID),
		lastID: 0,
	}
}

func (em *EntityManager) NewEntity() *Entity {
	e := NewEntity(em.getID())
	em.Entities[e.ID] = e
	return e
}

func (em *EntityManager) NameEntity(e *Entity) {
	name := e.Name
	existing, ok := em.NamedEntities[name]
	if ok {
		if existing == e.ID {
			return
		}
		panic(fmt.Sprintf("Duplicate entity name: %s (%v, %v)", name, existing, e))
	}
	em.NamedEntities[name] = e.ID
}

func (em *EntityManager) GetNamedEntity(name string) *Entity {
	eID, ok := em.NamedEntities[name]
	if !ok {
		return nil
	}
	return em.Entities[eID]
}

func (em *EntityManager) getID() EntityID {
	if (em.lastID == math.MaxUint64) {
		panic("Ran out of IDs!")
	}
	em.lastID++
	return em.lastID
}


