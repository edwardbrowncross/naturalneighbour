package geom

import "math"

// GetArea returns the area of the triangle produced by the three points.
func GetArea(p1x, p1y, p2x, p2y, p3x, p3y float64) float64 {
	return 0.5 * math.Abs(Det3s(
		p1x, p1y,
		p2x, p2y,
		p3x, p3y,
	))
}

// IsClockwise returns whether the given points are defined in clockwise order.
func IsClockwise(p1x, p1y, p2x, p2y, p3x, p3y float64) bool {
	return Det3s(
		p1x, p1y,
		p2x, p2y,
		p3x, p3y,
	) < 0
}
