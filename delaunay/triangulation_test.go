package delaunay

import (
	"math"
	"math/rand"
	"testing"
)

func TestTriangulation(t *testing.T) {
	points := make([]*Point, 6)
	for i := 0; i < 5; i++ {
		points[i] = NewPoint(
			10*math.Sin(float64(i)*2*math.Pi/5),
			10*math.Cos(float64(i)*2*math.Pi/5),
			0,
		)
	}
	points[5] = NewPoint(0, 0, 0)
	NewTriangulation(points)
	c := points[5]
	if len(c.GetConnected()) != 5 {
		t.Errorf("epected centre point to have 5 connected points but got %d", len(c.GetConnected()))
	}
}

var result *Triangulation

func benchmarkTriangulation(n int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		points := make([]*Point, n)
		for j := 0; j < n; j++ {
			points[j] = NewPoint(rand.Float64(), rand.Float64(), 0)
		}
		b.StartTimer()
		result, _ = NewTriangulation(points)
	}
}

func BenchmarkTriangulation10(b *testing.B)      { benchmarkTriangulation(10, b) }
func BenchmarkTriangulation100(b *testing.B)     { benchmarkTriangulation(100, b) }
func BenchmarkTriangulation1000(b *testing.B)    { benchmarkTriangulation(1000, b) }
func BenchmarkTriangulation10000(b *testing.B)   { benchmarkTriangulation(10000, b) }
func BenchmarkTriangulation50000(b *testing.B)   { benchmarkTriangulation(50000, b) }
func BenchmarkTriangulation100000(b *testing.B)  { benchmarkTriangulation(100000, b) }
func BenchmarkTriangulation1000000(b *testing.B) { benchmarkTriangulation(1000000, b) }
