package space

import (
	"fmt"
	"math"

	. "github.com/brendonh/glvec"
)

type AvatarPosition struct {
	BaseComponent
	Physics *SpacePhysics

	ShipPosition Vec2i
	Position SpacePosition
	PrevPosition SpacePosition
}

func (p *AvatarPosition) Tag() string {
	return "struct_avatarposition"
}

func (p *AvatarPosition) TickPhysics() {
	p.PrevPosition = p.Position
	// ...
}

func (p *AvatarPosition) Event(tag string, args interface{}) {
	switch(tag) {
	case "move_to":
		p.MoveTo(args.(Vec2i))
	}
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
	}
	p.Physics = newPhysics
}

func (p *AvatarPosition) AttachToShipPosition(e *Entity, pos Vec2i) {
	p.Physics = e.GetComponent("struct_spacephysics").(*SpacePhysics)

	p.ShipPosition = pos

	rooms := e.GetComponent("rooms").(*RoomsComponent)

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

func (p *AvatarPosition) MoveTo(tilePos Vec2i) {
	fmt.Println("Moveto", tilePos)
	//var currentTile = p.Rooms.Grid.Get(p.ShipPosition)
	// path, success := p.Rooms.Grid.FindPath(p.ShipPosition, tilePos)

	// if !success {
	// 	p.Rooms.Entity.BroadcastEvent("update_colors", []CubeColorOverride{})		
	// 	return
	// }

	// overrides := make([]CubeColorOverride, 0, len(path))
	// for _, tilePos := range path {
	// 	//shipPos := tile.GetShipPos()
	// 	overrides = append(overrides, CubeColorOverride{ 
	// 		tilePos.X, tilePos.Y, CubeColor{ 1.0, 0.0, 0.0, 0.5 } })
	// }

	// p.Rooms.Entity.BroadcastEvent("update_colors", overrides)
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