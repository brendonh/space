package space

type LogicSystem struct {
	Active []LogicComponent
	toAdd []LogicComponent
	toRemove []LogicComponent
}

func NewLogicSystem() *LogicSystem {
	return &LogicSystem{}
}

func (s *LogicSystem) Add(c LogicComponent) {
	s.toAdd = append(s.toAdd, c)
}

func (s *LogicSystem) Remove(c LogicComponent) {
	s.toRemove = append(s.toRemove, c)
}

func (s *LogicSystem) Update() {
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

func (s *LogicSystem) Tick() {
	for _, c := range s.Active {
		c.TickLogic()
	}
}
