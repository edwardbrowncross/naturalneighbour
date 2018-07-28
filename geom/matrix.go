package geom

// CurlZ returns the z coordinate of the curl of two vectors in the 2D plane.
func CurlZ(p1x, p1y, p2x, p2y, p3x, p3y float64) float64 {
	return (p1x-p3x)*(p2y-p3y) - (p2x-p3x)*(p1y-p3y)
}

// Det2 returns the determinant of a 2x2 matrix.
func Det2(a, b, c, d float64) float64 {
	return a*d - b*c
}

// Det3 returns the determinant of a 3x3 matrix.
func Det3(a, b, c, d, e, f, g, h, i float64) float64 {
	return a*(e*i-f*h) - b*(d*i-f*g) + c*(d*h-e*g)
}

// Det3s returns the determinant of a 3x3 matrix that has a column of all 1's.
// Equivalent to Det3(a, b, 1, d, e, 1, f, g, 1.)
func Det3s(a, b, d, e, f, g float64) float64 {
	return (d-a)*(g-b) + (b-e)*(f-a)
}
