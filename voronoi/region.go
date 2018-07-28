package voronoi

import (
	"github.com/edwardbrowncross/naturalneighbour/delaunay"
	"github.com/edwardbrowncross/naturalneighbour/geom"
)

// Region represents a voronoi cell.
type Region struct {
	Center *delaunay.Point // The point in the delaunay triangulation this cell is associated with.
	Verts  []Vertex        // The vertices of the polygon representing the bounds of this cell.
}

// NewRegion creates a new veronoi region for the given delaunay point.
func NewRegion(p *delaunay.Point) Region {
	verts := []Vertex{}
	// Vertices are the circumcenters of the delaunay triangles surrounding the point.
	// https://stackoverflow.com/questions/85275/how-do-i-derive-a-voronoi-diagram-given-its-point-set-and-its-delaunay-triangula#comment619809_85359
	// Consecutive verts should be from neighbouring triangles.
	t0 := p.Triangles[0]
	curt := t0
	curp := t0.Points[0]
	if curp == p {
		curp = t0.Points[1]
	}
	for true {
		v := NewVertex(curt.GetCircumcenter())
		verts = append(verts, v)
		newt := curt.GetAdjacentTo(p, curp)
		if newt == t0 || newt == nil {
			break
		}
		curp = newt.GetPointOpposite(curt)
		curt = newt
	}
	r := Region{
		Center: p,
		Verts:  verts,
	}
	return r
}

// GetArea returns the area of a voronoi cell.
func (r Region) GetArea() float64 {
	lv := len(r.Verts)
	// Decompose polygon into triangles radiating from the central point and sum areas.
	sum := 0.0
	for i, v1 := range r.Verts {
		v2 := r.Verts[(i+1)%lv]
		sum += geom.GetArea(r.Center.X, r.Center.Y, v1.X, v1.Y, v2.X, v2.Y)
	}
	return sum
}
