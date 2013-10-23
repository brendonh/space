package space

type Vec2i struct {
	X, Y int
}

func (v Vec2i) Add(o Vec2i) Vec2i {
	return Vec2i{
		v.X + o.X,
		v.Y + o.Y,
	}
}

func (v Vec2i) Neg() Vec2i {
	return Vec2i { -v.X, -v.Y }
}