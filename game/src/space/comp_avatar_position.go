package space

import (
	"math"

	. "github.com/brendonh/glvec"
)

type AvatarPosition struct {
	BaseComponent
	Physics *SpacePhysics

	ShipPosition Vec2i
	Moves []Vec3

	Position SpacePosition
	PrevPosition SpacePosition

	WalkSpeed float64
}

func (p *AvatarPosition) Tag() string {
	return "struct_avatarposition"
}

func (p *AvatarPosition) TickPhysics() {
	p.PrevPosition = p.Position

	pos := &p.Position

	if len(p.Moves) > 0 {
		for _, move := range p.Moves {
			V3Add(&pos.Pos, pos.Pos, move)
		}
		
		p.Moves = p.Moves[:0]
	}
}

func (p *AvatarPosition) AddMove(move Vec3) {
	p.Moves = append(p.Moves, move)
}

func (p *AvatarPosition) Attached() bool {
	return p.Physics.Entity != p.Entity
}

func (p *AvatarPosition) AttachedTo() *Entity {
	if !p.Attached() {
		return nil
	}
	return p.Physics.Entity
}

func (p *AvatarPosition) AttachTo(e *Entity) {
	oldPhysics := p.Physics
	newPhysics := e.GetComponent("struct_spacephysics").(*SpacePhysics)

	if oldPhysics != nil {
		if oldPhysics.Entity == p.Entity {
			p.Entity.RemoveComponent(oldPhysics)
		}

		p.physicsToPosition(oldPhysics, newPhysics, &p.PrevPosition, 0)
		p.physicsToPosition(oldPhysics, newPhysics, &p.Position, 1)

		// ??
		// rooms := e.GetComponent("rooms").(*RoomsComponent)
		// p.ShipPosition = rooms.ModelToTile(

	}
	p.Physics = newPhysics
}

func (p *AvatarPosition) AttachToShipPosition(e *Entity, pos Vec2i) bool {
	p.Physics = e.GetComponent("struct_spacephysics").(*SpacePhysics)
	p.SetShipPosition(pos)
	return true
}

func (p *AvatarPosition) SetShipPosition(pos Vec2i) {
	p.ShipPosition = pos

	rooms := p.Physics.Entity.GetComponent("rooms").(*RoomsComponent)
	modelPos := rooms.TileToModel(pos)
	p.Position = SpacePosition{ modelPos, 0 }
	p.PrevPosition = p.Position
}

func (p *AvatarPosition) Detach() {
	var phys SpacePhysics
	if p.Physics != nil {
		if p.Physics.Entity == p.Entity {
			return
		}

		phys = *p.Physics

		p.positionToPhysics(&phys.PrevPosition, 0)
		p.positionToPhysics(&phys.Position, 1)

		p.Position = SpacePosition{}
		p.PrevPosition = SpacePosition{}
	}
	p.Entity.AddComponent(&phys)
	p.Physics = &phys
}

func (p *AvatarPosition) GetModelMatrix(alpha float32) Mat4 {
	spacePos := p.Physics.GetModelMatrix(alpha)

	pos := p.Position.InterpolateFrom(p.PrevPosition, alpha)

	var groundPos Mat4
	M4MakeRotation(&groundPos, float32(pos.Angle), Vec3{ 0, 0, 1 })
	M4SetTransform(&groundPos, pos.Pos)

	M4MulM4(&groundPos, &spacePos, &groundPos)

	return groundPos
}


func (p *AvatarPosition) positionToPhysics(position *SpacePosition, alpha float32) {
	var spaceMat = p.GetModelMatrix(alpha)
	position.Pos = M4GetTransform(&spaceMat)
	position.Angle = float32(math.Atan2(float64(spaceMat[1]), float64(spaceMat[0])))
}


func (p *AvatarPosition) physicsToPosition(
	oldPhysics, newPhysics *SpacePhysics,
	position *SpacePosition,
	alpha float32) {

	var oldMat = oldPhysics.GetModelMatrix(alpha)
	var oldPos = M4GetTransform(&oldMat)

	var newMat = newPhysics.GetModelMatrix(alpha)
	var newPos = M4GetTransform(&newMat)

	var offset = Vec3{ oldPos[0] - newPos[0], oldPos[1] - newPos[1], 0 }

	var invRotation Mat3
	M4RotationMatrix(&invRotation, &newMat)
	M3Inverse(&invRotation, &invRotation)

	M3MulV3(&offset, &invRotation, offset)

	V3Add(&position.Pos, position.Pos, offset)

	oldAngle := math.Atan2(float64(oldMat[1]), float64(oldMat[0]))
	newAngle := math.Atan2(float64(newMat[1]), float64(newMat[0]))
	position.Angle = float32(oldAngle - newAngle)
}