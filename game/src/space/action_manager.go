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

func (m *ActionManager) AddAction(action *Action) {
	action.Manager = m
	if action.Avatar != nil {
		m.readyActions.PushBack(action)
		return
	}

	for e := m.idleAvatars.Front(); e != nil; e = e.Next() {
		av := e.Value.(*Entity)
		if !action.Prepare(av) {
			continue
		}

		m.idleAvatars.Remove(e)
		m.readyActions.PushBack(action)
		return
	}

	m.pendingActions.PushBack(action)
}


func (m *ActionManager) AddAvatar(avatar *Entity) {
	for e := m.pendingActions.Front(); e != nil; e = e.Next() {
		action := e.Value.(*Action)
		if !action.Prepare(avatar) {
			continue
		}
		m.pendingActions.Remove(e)
		m.readyActions.PushBack(action)
		return
	}

	fmt.Println("Idle Avatar:", avatar)
	m.idleAvatars.PushBack(avatar)
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
