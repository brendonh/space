package space

import (
	"fmt"
)

type Action struct {
	Manager *ActionManager
	Avatar *Entity
	Location Vec2i
	Path []Vec2i
}

func (a *Action) Prepare(avatar *Entity) ([]Vec2i, float64, bool) {
	if !a.qualified(avatar) {
		return nil, 0, false
	}

	position := avatar.GetComponent("struct_avatarposition").(*AvatarPosition)
	if position.AttachedTo() != a.Manager.Entity {
		fmt.Println("Avatar not attached to manager entity!", avatar, a.Manager.Entity)
		return nil, 0, false
	}

	path, cost, ok := a.Manager.Grid.FindPath(position.ShipPosition, a.Location)
	if !ok {
		fmt.Println("No path to action location")
		return nil, 0, false
	}

	cost /= position.WalkSpeed

	return path, cost, true
}

func (a *Action) Abandon() {
	a.Avatar = nil
	a.Path = nil
}

func (a *Action) qualified(avatar *Entity) bool {
	return true
}
