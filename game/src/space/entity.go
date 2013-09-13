package space

import (
	"fmt"
	"math"
)

type EntityID uint64

type Entity struct {
	ID EntityID
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
}

func (e *Entity) GetComponent(tag string) Component {
	for _, c := range e.Components {
		if c.Tag() == tag {
			return c
		}
	}
	panic(fmt.Sprintf("No such component: %s", tag))
}

func (e *Entity) InitComponents() {
	for _, c := range e.Components {
		c.Init()
	}
}

type EntityManager struct {
	Entities map[EntityID]*Entity
	lastID EntityID
}

func NewEntityManager() *EntityManager {
	return &EntityManager {
		Entities: make(map[EntityID]*Entity),
		lastID: 0,
	}
}

func (em *EntityManager) NewEntity() *Entity {
	e := NewEntity(em.getID())
	em.Entities[e.ID] = e
	return e
}

func (em *EntityManager) getID() EntityID {
	if (em.lastID == math.MaxUint64) {
		panic("Ran out of IDs!")
	}
	em.lastID++
	return em.lastID
}


