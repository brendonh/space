package space

import (
	"fmt"
	"math"

	. "github.com/brendonh/glvec"
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

		b.getManager().AddAvatar(b.Entity)
		return
	}

	if b.Move != nil && b.Move.Done() {
		fmt.Println("Move finished")
		b.Position.SetShipPosition(b.Move.ToPos)
		b.Move = nil
	}
	
	if b.Move == nil {
		b.assignNextMove()
	}

	if b.Move == nil {
		fmt.Println("No move")
		return
	}

	b.Move.Tick(b.Position)
}

func (b *AvatarBehaviour) assignNextMove() {
	if len(b.Action.Path) == 0 {
		fmt.Println("Action pathing finished")
		b.Action = nil
		return
	}

	nextTile := b.Action.Path[0]	
	currentTile := b.Position.ShipPosition
	distance := nextTile.Distance(currentTile)
	ticks := int(math.Floor(distance / b.Position.WalkSpeed))

	fmt.Println("Move:", currentTile, nextTile, distance, ticks)
	b.Move = NewAvatarMove(currentTile, nextTile, ticks)
	b.Action.Path = b.Action.Path[1:]
}

func (b *AvatarBehaviour) getManager() *ActionManager {
	ship := b.Position.AttachedTo()
	if ship == nil {
		return nil
	}
	manager := ship.GetComponent("action_manager")
	if manager == nil {
		return nil
	}
	return manager.(*ActionManager)
}


type AvatarMove struct {
	FromPos Vec2i
	ToPos Vec2i
	Step Vec3
	TotalTicks int
	CurrentTicks int
}

func NewAvatarMove(from, to Vec2i, ticks int) *AvatarMove {
	am := &AvatarMove {
		FromPos: from,
		ToPos: to,
		TotalTicks: ticks,
		CurrentTicks: 0,
	}

	step := to.Vec3()
	V3Sub(&step, step, from.Vec3())
	V3ScalarDiv(&am.Step, step, float32(ticks) / CUBE_SCALE)
	fmt.Println("Step:", am.Step)
	return am
}

func (m *AvatarMove) Tick(pos *AvatarPosition) {
	pos.AddMove(m.Step)
	m.CurrentTicks++
}

func (m *AvatarMove) Done() bool {
	return m.CurrentTicks == m.TotalTicks
}
