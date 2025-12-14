package common

import (
	"math/rand/v2"
	"strconv"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/geojson"
	mb "github.com/visendi-labs/mapbox-gl-gojs"
)

type KmeansRunner[T1 any] interface {
	GetCentroid([]T1) T1
	Distance(T1, T1) float64
}

func Kmeans[T1 any](k int, spec KmeansRunner[T1], data []T1) (centroids []T1, clusters [][]T1) {
	centroids = make([]T1, k)
	random := make([]T1, len(data))
	copy(random, data)
	rand.Shuffle(len(data), func(i, j int) {
		random[i], random[j] = random[j], random[i]
	})
	copy(centroids, random[:k])
	for range 5 { // 5 passes (skip centroid equality check..)
		clusters = make([][]T1, k)
		for _, l := range data {
			closestIndex, minDistance := 0, spec.Distance(l, centroids[0])
			for j := range k {
				d := spec.Distance(l, centroids[j])
				if d < minDistance {
					closestIndex, minDistance = j, d
				}
			}
			clusters[closestIndex] = append(clusters[closestIndex], l)
		}
		centroids = make([]T1, k)
		for i := range k {
			centroids[i] = spec.GetCentroid(clusters[i])
		}
	}
	return centroids, clusters
}

type KmeansPoint struct{}

func (KmeansPoint) Distance(p1, p2 orb.Point) float64 {
	return geo.DistanceHaversine(p1, p2)
}

func (KmeansPoint) GetCentroid(cluster []orb.Point) orb.Point {
	p1AvgLat, p1AvgLon := 0.0, 0.0
	l := float64(len(cluster))
	for _, p := range cluster {
		p1AvgLat += p.Lat() / l
		p1AvgLon += p.Lon() / l
	}
	return orb.Point{p1AvgLon, p1AvgLat}
}

type KmeansLine struct{}

func (k KmeansLine) Distance(l1, l2 orb.LineString) float64 {
	return geo.DistanceHaversine(l1[0], l2[0]) + geo.DistanceHaversine(l1[1], l2[1])
}

func (KmeansLine) GetCentroid(cluster []orb.LineString) orb.LineString {
	l := float64(len(cluster))
	p1AvgLat, p1AvgLon, p2AvgLat, p2AvgLon := 0.0, 0.0, 0.0, 0.0
	for _, c := range cluster {
		p1AvgLat += c[0].Lat() / l
		p1AvgLon += c[0].Lon() / l
		p2AvgLat += c[1].Lat() / l
		p2AvgLon += c[1].Lon() / l
	}
	return orb.LineString{orb.Point{p1AvgLon, p1AvgLat}, orb.Point{p2AvgLon, p2AvgLat}}
}

var points *geojson.FeatureCollection = geojson.NewFeatureCollection()
var lines *geojson.FeatureCollection = geojson.NewFeatureCollection()

func Init() {
	for range 10 {
		lat, lon := -50+rand.Float64()*100, -50+rand.Float64()*100
		for range 50 {
			p := geojson.NewFeature(orb.Point{lat + rand.Float64()*10, lon + rand.Float64()*10})
			p.Properties["radius"] = 2
			points.Append(p)
		}
		lat1, lon1 := -50+rand.Float64()*100, -50+rand.Float64()*100
		lat2, lon2 := -50+rand.Float64()*100, -50+rand.Float64()*100
		for range 25 {
			l := geojson.NewFeature(orb.LineString{
				orb.Point{lon1 + rand.Float64()*10, lat1 + rand.Float64()*10},
				orb.Point{lon2 + rand.Float64()*10, lat2 + rand.Float64()*10},
			})
			l.Properties["width"] = 1
			lines.Append(l)
		}
	}
}

func Example(token string) string {
	Init()
	return mb.NewGroup(
		mb.NewMap(mb.Map{Container: "map", AccessToken: token}),
		mb.NewMapOnLoad(
			mb.NewMapAddLayer(mb.MapLayer{Id: "my_lines", Type: "line",
				Paint:  mb.MapLayerPaint{LineWidth: []any{"get", "width"}, LineOpacity: 0.5},
				Source: mb.MapSource{Type: "geojson", Data: *lines},
			}),
			mb.NewMapAddLayer(mb.MapLayer{Id: "k_means_lines", Type: "line",
				Paint:  mb.MapLayerPaint{LineWidth: []any{"get", "width"}, LineOpacity: 0.9},
				Source: mb.MapSource{Type: "geojson", Data: *geojson.NewFeatureCollection()},
			}),
			mb.NewMapAddLayer(mb.MapLayer{Id: "my_points", Type: "circle",
				Paint:  mb.MapLayerPaint{CircleRadius: []any{"get", "radius"}, CircleOpacity: 0.5},
				Source: mb.MapSource{Type: "geojson", Data: *points},
			}),
			mb.NewMapAddLayer(mb.MapLayer{Id: "k_means_points", Type: "circle",
				Paint:  mb.MapLayerPaint{CircleRadius: []any{"get", "radius"}, CircleOpacity: 0.5},
				Source: mb.MapSource{Type: "geojson", Data: *geojson.NewFeatureCollection()},
			}),
		),
	).MustRenderDefault()
}

func KmeanClusterLines(kStr string) string {
	k, _ := strconv.Atoi(kStr)
	data := make([]orb.LineString, len(lines.Features))
	for i, l := range lines.Features {
		data[i] = l.Geometry.(orb.LineString)
	}
	centroids, clusters := Kmeans(k, KmeansLine{}, data)
	clusteredLines := geojson.NewFeatureCollection()
	for i, c := range centroids {
		f := geojson.NewFeature(c)
		f.Properties["width"] = 1 + len(clusters[i])/5
		clusteredLines.Append(f)
	}
	return mb.NewScript(mb.NewMapSourceSetData("k_means_lines", *clusteredLines)).MustRenderDefault()
}

func KmeanClusterPoints(kStr string) string {
	k, _ := strconv.Atoi(kStr)
	data := make([]orb.Point, len(points.Features))
	for i, l := range points.Features {
		data[i] = l.Geometry.(orb.Point)
	}
	centroids, clusters := Kmeans(k, KmeansPoint{}, data)
	clusteredLines := geojson.NewFeatureCollection()
	for i, c := range centroids {
		f := geojson.NewFeature(c)
		f.Properties["radius"] = 2 + len(clusters[i])/5
		clusteredLines.Append(f)
	}
	return mb.NewScript(mb.NewMapSourceSetData("k_means_points", *clusteredLines)).MustRenderDefault()
}
