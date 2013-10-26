package space

import (
	"fmt"
	"container/list"
)

type ActionManager struct {
	BaseComponent
	Grid TileGrid

	readyActions *list.List
	pendingActions *list.List

	idleAvatars *list.List
}

func NewActionManager() *ActionManager {
	return &ActionManager{
		readyActions: list.New(),
		pendingActions: list.New(),
		idleAvatars: list.New(),
	}
}

func (m *ActionManager) Tag() string {
	return "action_manager"
}

func (m *ActionManager) TickManagement() {
	for e := m.readyActions.Front(); e != nil; e = e.Next() {
		action := m.readyActions.Remove(e).(*Action)
		behaviour := action.Avatar.GetComponent("behaviour").(*AvatarBehaviour)
		behaviour.SetAction(action)
	}
}

func (m *ActionManager) HasIdlers() bool {
	return m.idleAvatars.Len() > 0
}

func (m *ActionManager) AddAction(action *Action) {
	action.Manager = m
	if action.Avatar != nil {
		m.readyActions.PushBack(action)
		return
	}

	var candidate *actionCandidate
	
	for e := m.idleAvatars.Front(); e != nil; e = e.Next() {
		av := e.Value.(*Entity)
		path, cost, ok := action.Prepare(av)
		if !ok { continue }
		candidate = candidate.Best(&actionCandidate{e, path, cost})
	}

	if candidate != nil {
		action.Avatar = candidate.Element.Value.(*Entity)
		action.Path = candidate.Path
		m.idleAvatars.Remove(candidate.Element)
		m.readyActions.PushBack(action)
		return
	}

	m.pendingActions.PushBack(action)
}

func (m *ActionManager) AddAvatar(avatar *Entity) {
	var candidate *actionCandidate

	for e := m.pendingActions.Front(); e != nil; e = e.Next() {
		action := e.Value.(*Action)
		path, cost, ok := action.Prepare(avatar)
		if !ok { continue }
		candidate = candidate.Best(&actionCandidate{e, path, cost})
	}

	if candidate != nil {
		action := candidate.Element.Value.(*Action)
		action.Avatar = avatar
		action.Path = candidate.Path
		m.pendingActions.Remove(candidate.Element)
		m.readyActions.PushBack(action)
		return
	}

	fmt.Println("Idle Avatar:", avatar)
	m.idleAvatars.PushBack(avatar)
}



type actionCandidate struct {
	Element *list.Element
	Path []Vec2i
	Cost float64
}

func (c *actionCandidate) Best(other *actionCandidate) *actionCandidate {
	if c == nil || c.Cost > other.Cost {
		return other
	}
	return c
}
