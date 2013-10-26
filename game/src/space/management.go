package space

type ManagementSystem struct {
	Active []ManagementComponent
	toAdd []ManagementComponent
	toRemove []ManagementComponent
}

func NewManagementSystem() *ManagementSystem {
	return &ManagementSystem{}
}

func (s *ManagementSystem) Add(c ManagementComponent) {
	s.toAdd = append(s.toAdd, c)
}

func (s *ManagementSystem) Remove(c ManagementComponent) {
	s.toRemove = append(s.toRemove, c)
}

func (s *ManagementSystem) Update() {
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

func (s *ManagementSystem) Tick() {
	for _, c := range s.Active {
		c.TickManagement()
	}
}
