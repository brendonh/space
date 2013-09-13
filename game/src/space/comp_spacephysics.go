package space

//import "fmt"

type SpacePhysicsComponent struct {
	PosX float64
	PosY float64

	VelX float64
	VelY float64

	AccX float64
	AccY float64
}

func (c *SpacePhysicsComponent) TickPhysics() {
	//fmt.Println("Ticking physics")
}
