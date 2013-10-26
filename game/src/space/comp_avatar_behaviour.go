package space

import (
	"fmt"
)

type AvatarBehaviour struct {
	BaseComponent
	Position *AvatarPosition
	Idle bool

	Action *Action
	Move *AvatarMove
}

func NewAvatarBehaviour() *AvatarBehaviour {
	return &AvatarBehaviour{}
}

func (b *AvatarBehaviour) Tag() string {
	return "behaviour"
}

func (b *AvatarBehaviour) Init() {
	b.Position = b.Entity.GetComponent("struct_avatarposition").(*AvatarPosition)
}

func (b *AvatarBehaviour) SetAction(action *Action) {
	fmt.Println("Starting action", b.Entity.Name, action)
	b.Action = action
	b.Idle = false
}

func (b *AvatarBehaviour) TickLogic() {
	if b.Idle {
		return
	}

	if b.Action == nil {
		b.Idle = true
		ship := b.Position.AttachedTo()
		if ship == nil {
			return
		}
		manager := ship.GetComponent("action_manager").(*ActionManager)
		manager.AddAvatar(b.Entity)
		return
	}

	if b.Move == nil {
		nextTile := b.Action.Path[0]
		fmt.Println("Getting next move", nextTile)
		
	}
}


type AvatarMove struct {
	FromPos Vec2i
	ToPos Vec2i
	TotalTicks int
	CurrentTicks int
}