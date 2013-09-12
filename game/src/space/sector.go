package space

import (
	"fmt"
)

type Sector struct {
	Entities map[EntityID]*Entity
	physics *PhysicsSystem
	logic *LogicSystem
	render *RenderSystem
}

func NewSector() *Sector {
	return &Sector {
		Entities: make(map[EntityID]*Entity),
		physics: NewPhysicsSystem(),
		logic: NewLogicSystem(),
		render: NewRenderSystem(),
	}
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
		fmt.Println("Found physics component:", c)		
		s.physics.Add(c)
	}

	if c, ok := c.(LogicComponent); ok {
		fmt.Println("Found logic component:", c)		
		s.logic.Add(c)
	}

	if c, ok := c.(RenderComponent); ok {
		fmt.Println("Found render component:", c)		
		s.render.Add(c)
	}

}

func (s *Sector) Tick() {
	s.physics.Tick()
	s.logic.Tick()
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
