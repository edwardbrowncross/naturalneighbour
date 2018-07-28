package delaunay

import (
	"errors"
	"fmt"

	"github.com/edwardbrowncross/naturalneighbour/geom"
)

// Triangle is a triangle within the delaunay triangulation.
type Triangle struct {
	Points   [3]*Point   // Vertices of triangle.
	Children []*Triangle // Triangles created by splitting of this triangle.
}

// NewTriangle creates a new triangle object with the three points as vertices.
func NewTriangle(p1, p2, p3 *Point) *Triangle {
	// Points should always be defined in clockwise order. Some things break if not.
	if !geom.IsClockwise(p1.X, p1.Y, p2.X, p2.Y, p3.X, p3.Y) {
		p1, p2 = p2, p1
	}
	t := &Triangle{
		Points:   [3]*Point{p1, p2, p3},
		Children: []*Triangle{},
	}
	// Register with the points that they are attached to this triangle.
	p1.addTriangle(t)
	p2.addTriangle(t)
	p3.addTriangle(t)
	return t
}

// getPoints returns the vertices of the triangle as separate return values.
func (t *Triangle) getPoints() (*Point, *Point, *Point) {
	return t.Points[0], t.Points[1], t.Points[2]
}

// Contains tests whether the given point's X and Y values lie inside the bounds of this triangle.
func (t *Triangle) Contains(p *Point) (bool, error) {
	if p == nil {
		return false, errors.New("unable to test triangle for containing nil point")
	}
	t1, t2, t3 := t.getPoints()
	// Test that the line between every vertex and the given point is clockwise
	// from the line from said vertex to the next vertex in the triangle.
	// (Assumes triangle defined in clockwise order).
	return geom.CurlZ(p.X, p.Y, t1.X, t1.Y, t2.X, t2.Y) <= 0 &&
		geom.CurlZ(p.X, p.Y, t2.X, t2.Y, t3.X, t3.Y) <= 0 &&
		geom.CurlZ(p.X, p.Y, t3.X, t3.Y, t1.X, t1.Y) <= 0, nil
}

// GetCircumcenter returns the coordinates of the circumcenter of this triangle.
// Adapted from https://gist.github.com/mutoo/5617691.
func (t *Triangle) GetCircumcenter() (x, y float64) {
	p1, p2, p3 := t.getPoints()
	m1 := p1.X*p1.X + p1.Y*p1.Y
	m2 := p2.X*p2.X + p2.Y*p2.Y
	m3 := p3.X*p3.X + p3.Y*p3.Y
	f := 1 / (2 * geom.Det3s(p1.X, p1.Y, p2.X, p2.Y, p3.X, p3.Y))
	x = f * geom.Det3s(m1, p1.Y, m2, p2.Y, m3, p3.Y)
	y = -f * geom.Det3s(m1, p1.X, m2, p2.X, m3, p3.X)
	return
}

// getChildContaining tests each of this triangles children and returns the one that contains the given point.
// If the triangle has no children, it returns nil.
func (t *Triangle) getChildContaining(p *Point) (*Triangle, error) {
	if p == nil {
		return nil, errors.New("unable to test children for containing nil point")
	}
	for _, c := range t.Children {
		if pInC, _ := c.Contains(p); pInC {
			return c, nil
		}
	}
	return nil, nil
}

// Search finds the leaf node of the triangle tree starting from this triangle that contains the given point.
// If the point is not in this triangle, returns nil.
func (t *Triangle) Search(p *Point) (leaf *Triangle, err error) {
	if p == nil {
		return nil, errors.New("unable to test children for containing nil point")
	}
	if pInT, _ := t.Contains(p); !pInT {
		return nil, nil
	}
	leaf = t
	for len(leaf.Children) != 0 {
		leaf, err = leaf.getChildContaining(p)
		if err != nil || leaf == nil {
			return
		}
	}
	return
}

// Insert splits this triangle into three new triangles, with the given point as a central vertex.
func (t *Triangle) Insert(p *Point) error {
	// Update triangles points are linked to.
	for _, p := range t.Points {
		if err := p.removeTriangle(t); err != nil {
			return err
		}
	}
	t.Children = []*Triangle{
		NewTriangle(t.Points[0], t.Points[1], p),
		NewTriangle(t.Points[1], t.Points[2], p),
		NewTriangle(t.Points[2], t.Points[0], p),
	}
	return nil
}

// Uninsert undoes an Insert operation, removing the child triangles created and removing the point from the triangulation.
func (t *Triangle) Uninsert() error {
	if len(t.Children) != 3 {
		return errors.New("can only uninsert from a triangle previously inserted into")
	}
	c1 := t.Children[0]
	c2 := t.Children[1]
	c3 := t.Children[2]
	if len(c1.Children) != 0 || len(c2.Children) != 0 || len(c3.Children) != 0 {
		return errors.New("cannot uninsert triangle whos children have been split")
	}
	// Update triangles points are linked to.
	for _, p := range c1.Points {
		p.removeTriangle(c1)
	}
	for _, p := range c2.Points {
		p.removeTriangle(c2)
	}
	for _, p := range c3.Points {
		p.removeTriangle(c3)
	}
	for _, p := range t.Points {
		p.addTriangle(t)
	}
	// Remove children.
	t.Children = []*Triangle{}
	return nil
}

// FlipWith takes two triangles that share a common edge and creates two new triangles, which together
// form the same quadrilateral, but whos common edge stretches between the two preiously unconnected points.
func (t1 *Triangle) FlipWith(t2 *Triangle) error {
	if t2 == nil {
		return errors.New("cannot flip with nil triangle")
	}
	// Determine which vertices are shared between the triangles.
	common := []*Point{}
	unique := []*Point{}
	for _, p1 := range t1.Points {
		found := false
		for _, p2 := range t2.Points {
			if p2 == p1 {
				common = append(common, p1)
				found = true
				break
			}
		}
		if !found {
			unique = append(unique, p1)
		}
	}
	for _, p1 := range t2.Points {
		found := false
		for _, p2 := range t1.Points {
			if p2 == p1 {
				found = true
				break
			}
		}
		if !found {
			unique = append(unique, p1)
		}
	}
	if len(unique) != 2 || len(common) != 2 {
		return fmt.Errorf("cannot flip triangles that do not share an edge (%d, %d)", len(unique), len(common))
	}
	// Update point -> triangle reverences.
	for _, p := range t1.Points {
		if err := p.removeTriangle(t1); err != nil {
			return fmt.Errorf("Failed to remove triangle from point: %v", err)
		}
	}
	for _, p := range t2.Points {
		if err := p.removeTriangle(t2); err != nil {
			return fmt.Errorf("Failed to remove triangle from point: %v", err)
		}
	}
	// Create new triangles.
	var t3, t4 *Triangle
	if geom.IsClockwise(unique[0].X, unique[0].Y, unique[1].X, unique[1].Y, common[0].X, common[0].Y) {
		t3 = NewTriangle(unique[0], unique[1], common[0])
		t4 = NewTriangle(unique[1], unique[0], common[1])
	} else {
		t3 = NewTriangle(unique[0], unique[1], common[1])
		t4 = NewTriangle(unique[1], unique[0], common[0])
	}
	t1.Children = []*Triangle{t3, t4}
	t2.Children = []*Triangle{t3, t4}
	return nil
}

// UnflipWith reverses a Flip operation, deleting the created child triangles from t1 and t2.
func (t1 *Triangle) UnflipWith(t2 *Triangle) error {
	if t2 == nil {
		return errors.New("cannot unflip with nil triangle")
	}
	if t1.Children[0] != t2.Children[0] && t1.Children[0] != t2.Children[1] {
		return errors.New("cannot unflip with triangle that was not created in the same flip operation")
	}
	c1 := t1.Children[0]
	c2 := t1.Children[1]
	if len(c1.Children) != 0 || len(c2.Children) != 0 {
		return errors.New("cannot unflip triangles whos children have been split")
	}
	// Update point -> triangle links.
	for _, p := range c1.Points {
		p.removeTriangle(c1)
	}
	for _, p := range c2.Points {
		p.removeTriangle(c2)
	}
	for _, p := range t1.Points {
		p.addTriangle(t1)
	}
	for _, p := range t2.Points {
		p.addTriangle(t2)
	}
	// Erase children.
	t1.Children = []*Triangle{}
	t2.Children = []*Triangle{}
	return nil
}

// IsDelaunayWith tests wheth the two involved triangles are locally delaunay.
// It does this by checking whether the opposing point of one triangle lies within the circumradius of the other triangle.
func (t1 *Triangle) IsDelaunayWith(t2 *Triangle) bool {
	p := t2.getPointOpposite(t1)
	if p == nil {
		return true
	}
	// http://www.cs.utah.edu/~csilva/courses/cpsc7960/pdf/boulos-DT.pdf (slide 7).
	a, b, c := t1.getPoints()
	lenP := p.X*p.X + p.Y*p.Y
	return geom.Det3(a.X-p.X, a.Y-p.Y, a.X*a.X+a.Y*a.Y-lenP,
		b.X-p.X, b.Y-p.Y, b.X*b.X+b.Y*b.Y-lenP,
		c.X-p.X, c.Y-p.Y, c.X*c.X+c.Y*c.Y-lenP) >= 0
}

// getTriangleOpposite takes one of the vertices of this triangle and returns the triangle bordering the edge of the
// triangle that does not contain that point.
// If p is not in t, returns nil. If t has no neighbouring triangle, returns nil.
func (t *Triangle) getTriangleOpposite(p *Point) *Triangle {
	if t.Points[0] == p {
		return t.getAdjacentTo(t.Points[1], t.Points[2])
	} else if t.Points[1] == p {
		return t.getAdjacentTo(t.Points[0], t.Points[2])
	} else if t.Points[2] == p {
		return t.getAdjacentTo(t.Points[0], t.Points[1])
	}
	return nil
}

// getPointOpposite takes a triangle, t2, that borders this triangle and returns the point in this triangle that
// does not form part of triangle t2.
func (t1 *Triangle) getPointOpposite(t2 *Triangle) *Point {
	for _, p := range t1.Points {
		found := false
		for _, t := range p.Triangles {
			if t == t2 {
				found = true
				break
			}
		}
		if !found {
			return p
		}
	}
	return nil
}

// getAdjacentTo takes two vertices from this triangle and returns the triangle adjoining this triangle along that edge.
func (t *Triangle) getAdjacentTo(p1, p2 *Point) *Triangle {
	for _, t1 := range p1.Triangles {
		for _, t2 := range p2.Triangles {
			if t1 == t2 && t1 != t {
				return t1
			}
		}
	}
	return nil
}
