package interpolation

import "github.com/edwardbrowncross/naturalneighbour/delaunay"

// NewPoint is an alias for delaunay.NewPoint.
func NewPoint(x, y, value float64) *delaunay.Point {
	return delaunay.NewPoint(x, y, value)
}
