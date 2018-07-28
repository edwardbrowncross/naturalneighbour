package interpolation

import (
	"math"
	"math/rand"
	"testing"

	"github.com/edwardbrowncross/naturalneighbour/delaunay"
)

var Epsilon float64 = 0.00000001

func TestInterpolator(t *testing.T) {
	points := make([]*delaunay.Point, 5)
	values := [5]float64{2.0, 3.0, 5.0, 7.0, 11.0}
	for i := 0; i < 5; i++ {
		points[i] = NewPoint(
			10*math.Sin(float64(i)*2*math.Pi/5),
			10*math.Cos(float64(i)*2*math.Pi/5),
			values[i],
		)
	}
	interpolator, err := New(points)
	if err != nil {
		t.Fatalf("error creating interpolator: %v", err)
	}
	result, err := interpolator.Interpolate(0, 0)
	if err != nil {
		t.Fatalf("error interpolating point: %v", err)
	}
	expected := (2.0 + 3.0 + 5.0 + 7.0 + 11.0) / 5.0
	if math.Abs(result-expected) > Epsilon {
		t.Errorf("expected result of %v but got %v", expected, result)
	}
}

var result float64

func benchmarkInterpolation(n int, b *testing.B) {
	b.StopTimer()
	dataPoints := make([]*delaunay.Point, n)
	dataPoints[0] = NewPoint(1, 1, 0)
	dataPoints[1] = NewPoint(1, -1, 0)
	dataPoints[2] = NewPoint(-1, -1, 0)
	dataPoints[3] = NewPoint(-1, 1, 0)
	for j := 4; j < n; j++ {
		dataPoints[j] = NewPoint(rand.Float64(), rand.Float64(), rand.Float64())
	}
	interpolator, err := New(dataPoints)
	if err != nil {
		b.Fatalf("error creating interpolator: %v", err)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		result, err = interpolator.Interpolate(rand.Float64(), rand.Float64())
		if err != nil {
			b.Errorf("error interpolating point: %v", err)
		}
	}
}

func BenchmarkInterpolation10(b *testing.B)      { benchmarkInterpolation(10, b) }
func BenchmarkInterpolation100(b *testing.B)     { benchmarkInterpolation(100, b) }
func BenchmarkInterpolation1000(b *testing.B)    { benchmarkInterpolation(1000, b) }
func BenchmarkInterpolation10000(b *testing.B)   { benchmarkInterpolation(10000, b) }
func BenchmarkInterpolation50000(b *testing.B)   { benchmarkInterpolation(50000, b) }
func BenchmarkInterpolation100000(b *testing.B)  { benchmarkInterpolation(100000, b) }
func BenchmarkInterpolation1000000(b *testing.B) { benchmarkInterpolation(1000000, b) }
