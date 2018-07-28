package delaunay

import "errors"

// Point represents a vertex in the delauany triangulation.
type Point struct {
	X         float64     // X location of point.
	Y         float64     // Y location of point.
	Value     float64     // Value associated with point.
	Triangles []*Triangle // The triangles (leaf nodes only) this point is a vertex of.
}

// NewPoint creates a new Point object.
func NewPoint(x, y, value float64) *Point {
	return &Point{
		X:         x,
		Y:         y,
		Value:     value,
		Triangles: []*Triangle{},
	}
}

// addTriangle adds a new triangle to the point's triangle array.
func (p *Point) addTriangle(t *Triangle) {
	p.Triangles = append(p.Triangles, t)
}

// removeTriangle removes a triangle from the point's triangle array.
func (p *Point) removeTriangle(t *Triangle) error {
	triangles := p.Triangles
	lt := len(triangles) - 1
	for i, tr := range triangles {
		if tr == t {
			triangles[lt], triangles[i] = triangles[i], triangles[lt]
			p.Triangles = triangles[:lt]
			return nil
		}
	}
	return errors.New("triangle not found")
}

// GetConnected gets a list of points that this point is connected to by an edge of the delaunay triangulation.
// Equivilently, these are the points whos voronoi cells share an edge with this point's voronoi cell.
func (p *Point) GetConnected() (r []*Point) {
	seen := map[*Point]bool{
		p: true,
	}
	for _, t := range p.Triangles {
		for _, pt := range t.Points {
			if seen[pt] {
				continue
			}
			r = append(r, pt)
			seen[pt] = true
		}
	}
	return
}
