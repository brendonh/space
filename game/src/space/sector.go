package space

import (
	"space/render"
)

type Sector struct {
	Entities map[EntityID]*Entity
	InputSystem *InputSystem
	ManagementSystem *ManagementSystem
	LogicSystem *LogicSystem
	PhysicsSystem *PhysicsSystem
	RenderSystem *RenderSystem
}

func NewSector() *Sector {
	sector := &Sector {
		Entities: make(map[EntityID]*Entity),
		InputSystem: NewInputSystem(),
		ManagementSystem: NewManagementSystem(),
		LogicSystem: NewLogicSystem(),
		PhysicsSystem: NewPhysicsSystem(),
		RenderSystem: NewRenderSystem(),

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

	if c, ok := c.(InputComponent); ok {
		s.InputSystem.Add(c)
	}

	if c, ok := c.(ManagementComponent); ok {
		s.ManagementSystem.Add(c)
	}

	if c, ok := c.(LogicComponent); ok {
		s.LogicSystem.Add(c)
	}

	if c, ok := c.(PhysicsComponent); ok {
		s.PhysicsSystem.Add(c)
	}

	if c, ok := c.(RenderComponent); ok {
		s.RenderSystem.Add(c)
	}
}

func (s *Sector) UnregisterComponent(c Component) {

	if c, ok := c.(InputComponent); ok {
		s.InputSystem.Remove(c)
	}

	if c, ok := c.(ManagementComponent); ok {
		s.ManagementSystem.Remove(c)
	}

	if c, ok := c.(LogicComponent); ok {
		s.LogicSystem.Remove(c)
	}

	if c, ok := c.(PhysicsComponent); ok {
		s.PhysicsSystem.Remove(c)
	}

	if c, ok := c.(RenderComponent); ok {
		s.RenderSystem.Remove(c)
	}
}

func (s *Sector) Tick() {
	s.ManagementSystem.Tick()
	s.LogicSystem.Tick()
	s.updateEntities()

	s.PhysicsSystem.Tick()
}

func (s *Sector) Render(context *render.Context, alpha float32) {
	s.RenderSystem.Render(context, alpha)
}

func (s *Sector) updateEntities() {
	s.ManagementSystem.Update()
	s.LogicSystem.Update()
	s.PhysicsSystem.Update()
	s.RenderSystem.Update()
}
