package space

type Component interface {
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
