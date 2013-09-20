package space

import (
	"space/render"
)

type Component interface {
	Tag() string
	SetEntity(*Entity)
	Init()
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
}

type PhysicalComponent interface {
	Component
	Weight() float64
}
