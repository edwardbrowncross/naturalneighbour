# naturalneighbour
A golang implementation of the Natural Neighbour Interpolation algorithm in 2D

[Wikipedia entry](https://en.wikipedia.org/wiki/Natural_neighbor_interpolation) for algorithm.

Takes a list of 2D input points with attached values. Gets interpolated values from arbitrary locations within the points.

## Usage

```golang
// Create 1000 random points of source data.
dataPoints := make([]*delaunay.Point, 1000)
for i := 0; i < 1000; i++ {
    dataPoints[i] = interpolation.NewPoint(rand.Float64(), rand.Float64(), rand.Float64())
}
// Create an interpolator that will use the source data.
interpolator, err := interpolation.New(dataPoints)
// Interpolate at the point (0.5, 0.5)
result, err := interpolator.Interpolate(0.5, 0.5)
```