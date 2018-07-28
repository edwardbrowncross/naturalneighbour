package delaunay

import (
	"fmt"
	"math"
)

// Triangulation represents a delaunay triangulation.
type Triangulation struct {
	Root *Triangle
}

// NewTriangulation creates a new triangulation object.
// Given points should to be randomly sorted for optimum efficicency of triangle tree.
func NewTriangulation(points []*Point) (*Triangulation, error) {
	// Create a single root triangle that contains all the given points.
	t := Triangulation{
		Root: getBoundingTriangle(points),
	}
	// Add each point to the triangulation one at a time.
	for _, p := range points {
		if _, err := t.addPoint(p, false); err != nil {
			return nil, err
		}
	}
	return &t, nil
}

// AddPoint adds a new point to the delaunay triangulation and returns a function that will remove said point again.
// If point is not inside the bounding triangle created at the start, returns an error.
func (t *Triangulation) AddPoint(p *Point) (Undo, error) {
	return t.addPoint(p, true)
}

// addPoint adds a new point to the delaunay triangulation and optionally returns a function that will remove said point again.
// http://web.mit.edu/alexmv/Public/6.850-lectures/lecture09.pdf
func (t *Triangulation) addPoint(p *Point, undoable bool) (Undo, error) {
	var ul undoList
	// Find leaf triangle to insert new point into.
	leaf, err := t.Root.Search(p)
	if err != nil {
		return nil, fmt.Errorf("error finding leaf triangle: %v", err)
	}
	if leaf == nil {
		return nil, fmt.Errorf("point (%f,%f) does not lie within bounds", p.X, p.Y)
	}
	// Insert into leaf triangle.
	if err := leaf.Insert(p); err != nil {
		return nil, err
	}
	if undoable {
		ul = undoList{}
		ul.Add(newUninserter(leaf))
	}
	// Check each of the three new triangles for being locally delaunay.
	toCheck := make([]*Triangle, len(p.Triangles))
	copy(toCheck, p.Triangles)
	for i := 0; i < len(toCheck); i++ {
		t1 := toCheck[i]
		t2 := t1.getTriangleOpposite(p)
		if t2 == nil {
			continue
		}
		if t1.IsDelaunayWith(t2) {
			continue
		}
		// Flip any triangles not delaunay.
		if err := t1.FlipWith(t2); err != nil {
			return nil, fmt.Errorf("could not flip triangles: %v", err)
		}
		if undoable {
			ul.Add(newUnflipper(t1, t2))
		}
		// Add the new triangles to the list of triangles to check.
		toCheck = append(toCheck, t1.Children[0], t1.Children[1])
	}
	return ul.Undo, nil
}

// getBounds gets the minimum and maximum x and y coordinates of any points in the given array.
func getBounds(points []*Point) (minX, minY, maxX, maxY float64) {
	minX = math.Inf(+1)
	minY = math.Inf(+1)
	maxX = math.Inf(-1)
	maxY = math.Inf(-1)
	for _, p := range points {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}
	return
}

// getBoundingTriangle returns a triangle that will encompass all the given points.
func getBoundingTriangle(points []*Point) *Triangle {
	minX, minY, maxX, maxY := getBounds(points)
	cx := (minX + maxX) / 2
	cy := (minY + maxY) / 2
	s := math.Max(maxX-minX, maxY-minY) / 2
	return NewTriangle(
		NewPoint(cx, cy+3*s, 0),
		NewPoint(cx+3*s, cy, 0),
		NewPoint(cx-3*s, cy-3*s, 0),
	)
}

// Undo is a function that will remove an added point from the triangulation,
// restoring it to how it was before that point was inserted.
type Undo func() error

type undoer interface {
	Undo() error
}

type unflipper struct {
	t1 *Triangle
	t2 *Triangle
}

func newUnflipper(t1, t2 *Triangle) unflipper {
	return unflipper{
		t1: t1,
		t2: t2,
	}
}
func (u unflipper) Undo() error {
	return u.t1.UnflipWith(u.t2)
}

type uninserter struct {
	t *Triangle
}

func newUninserter(t *Triangle) uninserter {
	return uninserter{
		t: t,
	}
}
func (u uninserter) Undo() error {
	return u.t.Uninsert()
}

type undoList struct {
	list []undoer
}

func (u *undoList) Add(new undoer) {
	u.list = append(u.list, new)
}

func (u undoList) Undo() error {
	for i := len(u.list) - 1; i >= 0; i-- {
		if err := u.list[i].Undo(); err != nil {
			return err
		}
	}
	return nil
}
