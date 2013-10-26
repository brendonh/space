package space

import (
	"fmt"
)

type AvatarBehaviour struct {
	BaseComponent
	Position *AvatarPosition

	Action *Action
	Idle bool
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
}