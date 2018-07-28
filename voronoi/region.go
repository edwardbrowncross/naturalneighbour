package voronoi

import (
	"sort"

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
	for _, t := range p.Triangles {
		v := NewVertex(t.GetCircumcenter())
		verts = append(verts, v)
	}
	r := Region{
		Center: p,
		Verts:  verts,
	}
	// Sort vertices in clockwise orders (consecutive verts should be from neighbouring triangles).
	sort.Sort(r)
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

// Sorting logic.
// Clockwise sort: https://stackoverflow.com/questions/6989100/sort-points-in-clockwise-order
func (r Region) Len() int {
	return len(r.Verts)
}
func (r Region) Swap(i, j int) {
	r.Verts[i], r.Verts[j] = r.Verts[j], r.Verts[i]
}
func (r Region) Less(i, j int) bool {
	return (r.Verts[i].X-r.Center.X)*(r.Verts[j].Y-r.Center.Y) < (r.Verts[j].X-r.Center.X)*(r.Verts[i].Y-r.Center.Y)
}
