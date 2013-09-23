package space

import (
	"space/render"
)

type RenderSystem struct {
	Active []RenderComponent
	toAdd []RenderComponent
	toRemove []RenderComponent
}

func NewRenderSystem() *RenderSystem {
	return &RenderSystem{}
}

func (s *RenderSystem) Add(c RenderComponent) {
	s.toAdd = append(s.toAdd, c)
}

func (s *RenderSystem) Remove(c RenderComponent) {
	s.toRemove = append(s.toRemove, c)
}

func (s *RenderSystem) Update() {
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

func (s *RenderSystem) Render(context *render.Context, alpha float64) {
	for _, c := range s.Active {
		c.Render(context, alpha)
	}
}

func (s *RenderSystem) Iterate(cb func(RenderComponent) bool) {
	for _, c := range s.Active {
		if cb(c) {
			break
		}
	}
}