package space

import (
	"sort"

	glfw "github.com/go-gl/glfw3"
)


var KEYS = map[string][]glfw.Key {
	"quit_game": []glfw.Key { glfw.KeyEscape },

	"ship_accel": []glfw.Key { glfw.KeyUp, glfw.KeyW },
	"ship_decel": []glfw.Key { glfw.KeyDown, glfw.KeyS },
	"ship_left": []glfw.Key { glfw.KeyLeft, glfw.KeyA },
	"ship_right": []glfw.Key { glfw.KeyRight, glfw.KeyD },
}


type InputComponent interface {
	Component
	Priority() int
	Actions() []string
	KeyDown(string) bool
	KeyUp(string)
}


type KeyHandler struct {
	Handler InputComponent
	Action string
}

type KeyHandlerSet struct {
	Key glfw.Key
	handlers []*KeyHandler
	active *KeyHandler
}

func (s *KeyHandlerSet) Add(handler *KeyHandler) {
	s.handlers = append(s.handlers, handler)
	sort.Sort(s)
}

func (s *KeyHandlerSet) Dispatch(action glfw.Action) {
	if s == nil {
		return
	}

	switch (action) {
	case glfw.Release:
		if s.active != nil {
			s.active.Handler.KeyUp(s.active.Action)
			s.active = nil
		}
		return
	case glfw.Press:
		for _, handler := range s.handlers {
			if handler.Handler.KeyDown(handler.Action) {
				s.active = handler
				break
			}
		}
	}
}

func (s *KeyHandlerSet) Len() int {
	return len(s.handlers)
}

func (s *KeyHandlerSet) Swap(i, j int) {
	s.handlers[i], s.handlers[j] = s.handlers[j], s.handlers[i]
}

func (s *KeyHandlerSet) Less(i, j int) bool {
	return s.handlers[i].Handler.Priority() > s.handlers[j].Handler.Priority()
}


// ------------------

type InputSystem struct {
	handlers map[glfw.Key]*KeyHandlerSet
	active map[glfw.Key]InputComponent
}

func NewInputSystem() *InputSystem {
	return &InputSystem {
		handlers: make(map[glfw.Key]*KeyHandlerSet),
		active: make(map[glfw.Key]InputComponent),
	}
}

func (is *InputSystem) Add(c InputComponent) {
	for _, action := range c.Actions() {
		keys, ok := KEYS[action]
		if !ok {
			panic("Registering unknown input action: " + action)
		}

		for _, key := range keys {
			handlerSet, ok := is.handlers[key]
			if !ok {
				handlerSet = &KeyHandlerSet { Key: key }
				is.handlers[key] = handlerSet
			}
			handlerSet.Add(&KeyHandler {
				Handler: c,
				Action: action,
			})
		}
	}
}

func (is *InputSystem) HandleKey(key glfw.Key, action glfw.Action, mods glfw.ModifierKey) {
	is.handlers[key].Dispatch(action)
}
