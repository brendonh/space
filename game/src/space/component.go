package space

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
	Render(*RenderContext)
}

type PhysicalComponent interface {
	Component
	Weight() float64
}
