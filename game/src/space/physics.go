package space

type PhysicsSystem struct {
	Active []PhysicsComponent
	toAdd []PhysicsComponent
	toRemove []PhysicsComponent
}

func NewPhysicsSystem() *PhysicsSystem {
	return &PhysicsSystem{}
}

func (s *PhysicsSystem) Add(c PhysicsComponent) {
	s.toAdd = append(s.toAdd, c)
}

func (s *PhysicsSystem) Remove(c PhysicsComponent) {
	s.toRemove = append(s.toRemove, c)
}

func (s *PhysicsSystem) Update() {
	if len(s.toAdd) > 0 {
		s.Active = append(s.Active, s.toAdd...)
		s.toAdd = nil
	}

	for _, c := range s.toRemove {
		for i, el := range s.Active {
			if el == c {
				var a = s.Active
				a[len(a)-1], a[i], a = nil, a[len(a)-1], a[:len(a)-1]
				s.Active = a
				break
			}
		}
	}
	s.toRemove = nil
}

func (s *PhysicsSystem) Tick() {
	for _, c := range s.Active {
		c.TickPhysics()
	}
}
