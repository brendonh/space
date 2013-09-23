package space

import (
	"space/render"
)

type Sector struct {
	Entities map[EntityID]*Entity
	PhysicsSystem *PhysicsSystem
	LogicSystem *LogicSystem
	RenderSystem *RenderSystem
	InputSystem *InputSystem
}

func NewSector() *Sector {
	sector := &Sector {
		Entities: make(map[EntityID]*Entity),
		PhysicsSystem: NewPhysicsSystem(),
		LogicSystem: NewLogicSystem(),
		RenderSystem: NewRenderSystem(),
		InputSystem: NewInputSystem(),
	}
	sector.RegisterComponent(NewDustfield())
	return sector
}

func (s *Sector) AddEntity(e *Entity) {
	e.Sector = s

	s.Entities[e.ID] = e
	
	for _, c := range e.Components {
		s.RegisterComponent(c)
	}
}

func (s *Sector) RegisterComponent(c Component) {
	if c, ok := c.(PhysicsComponent); ok {
		s.PhysicsSystem.Add(c)
	}

	if c, ok := c.(LogicComponent); ok {
		s.LogicSystem.Add(c)
	}

	if c, ok := c.(RenderComponent); ok {
		s.RenderSystem.Add(c)
	}

	if c, ok := c.(InputComponent); ok {
		s.InputSystem.Add(c)
	}
}

func (s *Sector) Tick() {
	s.LogicSystem.Tick()
	s.PhysicsSystem.Tick()
	s.updateEntities()
}

func (s *Sector) Render(context *render.Context, alpha float64) {
	s.RenderSystem.Render(context, alpha)
}

func (s *Sector) updateEntities() {
	s.PhysicsSystem.Update()
	s.LogicSystem.Update()
	s.RenderSystem.Update()
}
