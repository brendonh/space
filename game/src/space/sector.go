package space

import (
	"space/render"
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

func (s *Sector) Render(context *render.Context, alpha float64) {
	s.render.Render(context, alpha)
}

func (s *Sector) updateEntities() {
	s.physics.Update()
	s.logic.Update()
	s.render.Update()
}
