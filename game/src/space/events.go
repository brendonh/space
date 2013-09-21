package space


type EventHandler func(interface{})

type ListenerSet map[ComponentID]EventHandler


type Dispatcher struct {
	events map[string]ListenerSet
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		events: make(map[string]ListenerSet),
	}
}

func (d *Dispatcher) Listen(component Component, event string, handler EventHandler) {
	set, ok := d.events[event]
	if !ok {
		set = make(ListenerSet)
		d.events[event] = set
	}
	set[component.ID()] = handler
}

func (d *Dispatcher) Remove(component Component, event string) {
	delete(d.events[event], component.ID())
}

func (d *Dispatcher) Fire(event string, args interface{}) {
	for _, handler := range d.events[event] {
		handler(args)
	}
}

var globalDispatch = NewDispatcher()