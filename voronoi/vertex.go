package voronoi

type Vertex struct {
	X float64
	Y float64
}

func NewVertex(x, y float64) Vertex {
	return Vertex{
		X: x,
		Y: y,
	}
}
