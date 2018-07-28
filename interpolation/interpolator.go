package interpolation

import (
	"github.com/edwardbrowncross/naturalneighbour/delaunay"
	"github.com/edwardbrowncross/naturalneighbour/voronoi"
)

// Interpolator provides natural neighbour interpolation within a set of points.
type Interpolator struct {
	t *delaunay.Triangulation
}

// New creates a new Interpolator using the given points.
func New(points []*delaunay.Point) (*Interpolator, error) {
	t, err := delaunay.NewTriangulation(points)
	return &Interpolator{t: t}, err
}

// Interpolate returns the interpolated value at the given x and y coordinates using natural neighbour interpolation.
// https://pdfs.semanticscholar.org/52ca/255573eded0e4371fe2ced980b196636718d.pdf
func (i *Interpolator) Interpolate(x, y float64) (float64, error) {
	// Create a new point and add it to the triangulation.
	p := delaunay.NewPoint(x, y, 0)
	undo, err := i.t.AddPoint(p)
	if err != nil {
		return 0, err
	}
	// Calculate the area of the voronoi cells of the points linked to the new point.
	neighbours := p.GetConnected()
	areasAfter := make([]float64, len(neighbours))
	// Calculate the area of the new test point's voronoi cells.
	for i, n := range neighbours {
		areasAfter[i] = voronoi.NewRegion(n).GetArea()
	}
	totalArea := voronoi.NewRegion(p).GetArea()
	undo()
	// Calculate the area of the same points without the new point in the triangulation.
	areasBefore := make([]float64, len(neighbours))
	for i, n := range neighbours {
		areasBefore[i] = voronoi.NewRegion(n).GetArea()
	}
	// Take a weighted average of the values of the points the new test point was connected to.
	// Weighting is the percentage of the test point's voronoi cell that was stolen from each neighbour point.
	total := 0.0
	for i, n := range neighbours {
		total += n.Value * (areasBefore[i] - areasAfter[i])
	}
	return total / totalArea, nil
}
