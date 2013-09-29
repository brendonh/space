package space

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	glfw "github.com/go-gl/glfw3"
)

const (
	INPUT_KEY   uint8 = iota
	INPUT_MOUSE
)

type Input struct {
	Type uint8
	Input int
}


type InputComponent interface {
	Component
	Priority() int
	Actions() []string
	KeyDown(string) bool
	KeyUp(string)
}


type InputHandler struct {
	Handler InputComponent
	Action string
}

type InputHandlerSet struct {
	Input Input
	handlers []*InputHandler
	active *InputHandler
}

func (s *InputHandlerSet) Add(handler *InputHandler) {
	s.handlers = append(s.handlers, handler)
	sort.Sort(s)
}

func (s *InputHandlerSet) RemoveComponent(c InputComponent) {
	for i, h := range s.handlers {
		if c == h.Handler {
			var hs = s.handlers
			hs[len(hs)-1], hs[i], hs = nil, hs[len(hs)-1], hs[:len(hs)-1]
			s.handlers = hs
		}
	}
}

func (s *InputHandlerSet) Dispatch(action glfw.Action) {
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

func (s *InputHandlerSet) Len() int {
	return len(s.handlers)
}

func (s *InputHandlerSet) Swap(i, j int) {
	s.handlers[i], s.handlers[j] = s.handlers[j], s.handlers[i]
}

func (s *InputHandlerSet) Less(i, j int) bool {
	return s.handlers[i].Handler.Priority() > s.handlers[j].Handler.Priority()
}


// ------------------

type InputSystem struct {
	bindings map[string][]Input
	handlers map[Input]*InputHandlerSet
	active map[Input]InputComponent
}

func NewInputSystem() *InputSystem {
	return &InputSystem {
		bindings: readInputBindings(),
		handlers: make(map[Input]*InputHandlerSet),
		active: make(map[Input]InputComponent),
	}
}

func (is *InputSystem) Add(c InputComponent) {
	for _, action := range c.Actions() {
		inputs, ok := is.bindings[action]
		if !ok {
			panic("Registering unknown input action: " + action)
		}

		for _, input := range inputs {
			handlerSet, ok := is.handlers[input]
			if !ok {
				handlerSet = &InputHandlerSet { Input: input }
				is.handlers[input] = handlerSet
			}
			handlerSet.Add(&InputHandler {
				Handler: c,
				Action: action,
			})
		}
	}
}

func (is *InputSystem) Remove(c InputComponent) {
	for _, action := range c.Actions() {
		inputs, ok := is.bindings[action]
		if !ok {
			panic("Unregistering unknown input action: " + action)
		}

		for _, input := range inputs {
			handlerSet, ok := is.handlers[input]
			if !ok {
				continue
			}
			handlerSet.RemoveComponent(c)
		}
	}
}

func (is *InputSystem) HandleKey(key glfw.Key, action glfw.Action, mods glfw.ModifierKey) {
	input := Input{ INPUT_KEY, int(key) }
	is.handlers[input].Dispatch(action)
}

func (is *InputSystem) HandleMouse(button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	input := Input{ INPUT_MOUSE, int(button) }
	is.handlers[input].Dispatch(action)
}


// XXX TODO: Handle out-of-window via CursorEnterCallback
func (is *InputSystem) UpdateMouse() {
	x, y := mainloop.RenderContext.Window.GetCursorPosition()
	ray := mainloop.RenderContext.ScreenToWorld(x, y)

	mainloop.Sector.RenderSystem.Iterate(func(c RenderComponent) bool {
		if c.HandleMouse(ray) {
			return true
		}
		return false
	})
}


func readInputBindings() map[string][]Input {
	bytes, err := ioutil.ReadFile("data/bindings.json")
	if err != nil {
		panic(fmt.Sprintf("Couldn't read bindings: %v", err))
	}

	var config map[string][]string
	err = json.Unmarshal(bytes, &config)

	if err != nil {
		panic(fmt.Sprintf("Couldn't decode bindings: %v", err))
	}

	bindings := make(map[string][]Input)

	for action, inputNames := range config {
		var inputs []Input
		for _, inputName := range inputNames {
			inputName = strings.ToLower(inputName)

			key, ok := glfwKeyNames[inputName]
			if ok {
				inputs = append(inputs, Input{ INPUT_KEY, int(key) })
				continue
			}

			button, ok := glfwMouseNames[inputName]
			if ok {
				inputs = append(inputs, Input{ INPUT_MOUSE, int(button) })
				continue
			}

			panic(fmt.Sprintf("Unknown input name: %v", inputName))
		}

		bindings[action] = inputs
	}

	return bindings
}

var glfwMouseNames = map[string]glfw.MouseButton {
	"mouse1": glfw.MouseButton1,
	"mouse2": glfw.MouseButton2,
	"mouse3": glfw.MouseButton3,
	"mouse4": glfw.MouseButton4,
	"mouse5": glfw.MouseButton5,
	"mouse6": glfw.MouseButton6,
	"mouse7": glfw.MouseButton7,
	"mouse8": glfw.MouseButton8,
	"mouselast": glfw.MouseButtonLast,
	"mouseleft": glfw.MouseButtonLeft,
	"mouseright": glfw.MouseButtonRight,
	"mousemiddle": glfw.MouseButtonMiddle,
}

var glfwKeyNames = map[string]glfw.Key {
	"unknown": glfw.KeyUnknown,
	"space": glfw.KeySpace,
	"apostrophe": glfw.KeyApostrophe,
	"comma": glfw.KeyComma,
	"minus": glfw.KeyMinus,
	"period": glfw.KeyPeriod,
	"slash": glfw.KeySlash,
	"0": glfw.Key0,
	"1": glfw.Key1,
	"2": glfw.Key2,
	"3": glfw.Key3,
	"4": glfw.Key4,
	"5": glfw.Key5,
	"6": glfw.Key6,
	"7": glfw.Key7,
	"8": glfw.Key8,
	"9": glfw.Key9,
	"semicolon": glfw.KeySemicolon,
	"equal": glfw.KeyEqual,
	"a": glfw.KeyA,
	"b": glfw.KeyB,
	"c": glfw.KeyC,
	"d": glfw.KeyD,
	"e": glfw.KeyE,
	"f": glfw.KeyF,
	"g": glfw.KeyG,
	"h": glfw.KeyH,
	"i": glfw.KeyI,
	"j": glfw.KeyJ,
	"k": glfw.KeyK,
	"l": glfw.KeyL,
	"m": glfw.KeyM,
	"n": glfw.KeyN,
	"o": glfw.KeyO,
	"p": glfw.KeyP,
	"q": glfw.KeyQ,
	"r": glfw.KeyR,
	"s": glfw.KeyS,
	"t": glfw.KeyT,
	"u": glfw.KeyU,
	"v": glfw.KeyV,
	"w": glfw.KeyW,
	"x": glfw.KeyX,
	"y": glfw.KeyY,
	"z": glfw.KeyZ,
	"leftbracket": glfw.KeyLeftBracket,
	"backslash": glfw.KeyBackslash,
	"bracket": glfw.KeyBracket,
	"rightbracket": glfw.KeyRightBracket,
	"graveaccent": glfw.KeyGraveAccent,
	"world1": glfw.KeyWorld1,
	"world2": glfw.KeyWorld2,
	"escape": glfw.KeyEscape,
	"enter": glfw.KeyEnter,
	"tab": glfw.KeyTab,
	"backspace": glfw.KeyBackspace,
	"insert": glfw.KeyInsert,
	"delete": glfw.KeyDelete,
	"right": glfw.KeyRight,
	"left": glfw.KeyLeft,
	"down": glfw.KeyDown,
	"up": glfw.KeyUp,
	"pageup": glfw.KeyPageUp,
	"pagedown": glfw.KeyPageDown,
	"home": glfw.KeyHome,
	"end": glfw.KeyEnd,
	"capslock": glfw.KeyCapsLock,
	"scrolllock": glfw.KeyScrollLock,
	"numlock": glfw.KeyNumLock,
	"printscreen": glfw.KeyPrintScreen,
	"pause": glfw.KeyPause,
	"f1": glfw.KeyF1,
	"f2": glfw.KeyF2,
	"f3": glfw.KeyF3,
	"f4": glfw.KeyF4,
	"f5": glfw.KeyF5,
	"f6": glfw.KeyF6,
	"f7": glfw.KeyF7,
	"f8": glfw.KeyF8,
	"f9": glfw.KeyF9,
	"f10": glfw.KeyF10,
	"f11": glfw.KeyF11,
	"f12": glfw.KeyF12,
	"f13": glfw.KeyF13,
	"f14": glfw.KeyF14,
	"f15": glfw.KeyF15,
	"f16": glfw.KeyF16,
	"f17": glfw.KeyF17,
	"f18": glfw.KeyF18,
	"f19": glfw.KeyF19,
	"f20": glfw.KeyF20,
	"f21": glfw.KeyF21,
	"f22": glfw.KeyF22,
	"f23": glfw.KeyF23,
	"f24": glfw.KeyF24,
	"f25": glfw.KeyF25,
	"kp0": glfw.KeyKp0,
	"kp1": glfw.KeyKp1,
	"kp2": glfw.KeyKp2,
	"kp3": glfw.KeyKp3,
	"kp4": glfw.KeyKp4,
	"kp5": glfw.KeyKp5,
	"kp6": glfw.KeyKp6,
	"kp7": glfw.KeyKp7,
	"kp8": glfw.KeyKp8,
	"kp9": glfw.KeyKp9,
	"kpdecimal": glfw.KeyKpDecimal,
	"kpdivide": glfw.KeyKpDivide,
	"kpmultiply": glfw.KeyKpMultiply,
	"kpsubtract": glfw.KeyKpSubtract,
	"kpadd": glfw.KeyKpAdd,
	"kpenter": glfw.KeyKpEnter,
	"kpequal": glfw.KeyKpEqual,
	"leftshift": glfw.KeyLeftShift,
	"leftcontrol": glfw.KeyLeftControl,
	"leftalt": glfw.KeyLeftAlt,
	"leftsuper": glfw.KeyLeftSuper,
	"rightshift": glfw.KeyRightShift,
	"rightcontrol": glfw.KeyRightControl,
	"rightalt": glfw.KeyRightAlt,
	"rightsuper": glfw.KeyRightSuper,
	"menu": glfw.KeyMenu,
	"last": glfw.KeyLast,
}