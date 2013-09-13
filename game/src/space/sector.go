package space

import (
	"fmt"
)

type Sector struct {
	Entities map[EntityID]*Entity
	physics *PhysicsSystem
	logic *LogicSystem
	render *RenderSystem
	Input *InputSystem
}

func NewSector() *Sector {
	sector := &Sector {
		Entities: make(map[EntityID]*Entity),
		physics: NewPhysicsSystem(),
		logic: NewLogicSystem(),
		render: NewRenderSystem(),
		Input: NewInputSystem(),
	}
	sector.RegisterComponent(NewStarfield())
	return sector
}

func (s *Sector) AddEntity(e *Entity) {
	e.Sector = s

	fmt.Println("Adding entity", e)
	s.Entities[e.ID] = e
	
	for _, c := range e.Components {
		s.RegisterComponent(c)
	}
}

func (s *Sector) RegisterComponent(c Component) {
	if c, ok := c.(PhysicsComponent); ok {
		s.physics.Add(c)
	}

	if c, ok := c.(LogicComponent); ok {
		s.logic.Add(c)
	}

	if c, ok := c.(RenderComponent); ok {
		s.render.Add(c)
	}

	if c, ok := c.(InputComponent); ok {
		s.Input.Add(c)
	}
}

func (s *Sector) Tick() {
	s.logic.Tick()
	s.physics.Tick()
	s.updateEntities()
}

func (s *Sector) Render(context *RenderContext) {
	s.render.Render(context)
}

func (s *Sector) updateEntities() {
	s.physics.Update()
	s.logic.Update()
	s.render.Update()
}
