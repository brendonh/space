package space

import (
	"math"

	"space/render"
	. "github.com/brendonh/glvec"
)

type Component interface {
	ID() ComponentID
	Tag() string
	Init()
	SetEntity(*Entity)
	GetEntity() *Entity
	Event(tag string, args interface{})
}

type PhysicsComponent interface {
	Component
	TickPhysics()
}

type LogicComponent interface {
	Component
	TickLogic()
}

type RenderComponent interface {
	Component
	Render(context *render.Context, alpha float64)
	HandleMouse(Ray) bool
}

type PhysicalComponent interface {
	Component
	Weight() float64
}

// --------------------------------------------------

type BaseComponent struct {
	id ComponentID
	Entity *Entity
}

func NewBaseComponent() BaseComponent {
	return BaseComponent{
		id: getComponentID(),
	}
}

func (c *BaseComponent) ID() ComponentID {
	return c.id
}

func (c *BaseComponent) SetEntity(e *Entity) {
	c.Entity = e
}

func (c *BaseComponent) GetEntity() *Entity {
	return c.Entity
}

func (c *BaseComponent) Init() {
}

func (c *BaseComponent) Event(tag string, args interface{}) {
}

// --------------------------------------------------

type ComponentID uint64

var lastComponentID ComponentID

func getComponentID() ComponentID {
	if (lastComponentID == math.MaxUint64) {
		panic("Ran out of IDs!")
	}
	lastComponentID++
	return lastComponentID
}


