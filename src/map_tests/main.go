package main

import (
	"fmt"
	gsm "google_static_map"
	"sphere"
)

const earthRadiusMiles = 3959.0

func setPolyLineColors(pl *gsm.PolyLine, color string) {
	pl.SetWeight(1)
	pl.SetColor(color + "c0")
	pl.SetFillColor(color + "10")
}

func makePolyLine(points []*sphere.NVector) *gsm.PolyLine {
	pl := gsm.NewPolyLine()
	for _, point := range points {
		lat, lon := point.ToLatLonDegrees()
		pl.AddPointLatLon(lat, lon)
	}
	return pl
}

func main() {
	m := gsm.NewMap(640, 640, 2)
	zoom := 11
	m.Zoom = &zoom

	annArbor := sphere.NewNVectorFromLatLongDeg(42.281389, -83.748333)
	annandale := sphere.NewNVectorFromLatLongDeg(42.0128, -73.9082)
	radford := sphere.NewNVectorFromLatLongDeg(37.1275, -80.569444)
	pinyonCrest := sphere.NewNVectorFromLatLongDeg(33.603845, -116.44001)
	laQuinta := sphere.NewNVectorFromLatLongDeg(33.659093, -116.3188)

	colors := [...]string{"0xff0000", "0x0000ff", "0x009900", "0xffcc00"}

	var locations []*sphere.NVector

	locations = []*sphere.NVector{annArbor, annandale, radford}
	locations = []*sphere.NVector{pinyonCrest, laQuinta}
	// locations = []*sphere.NVector{pinyonCrest}
	// locations = []*sphere.NVector{laQuinta}

	var radii []float64

	// radii = []float64{1.0, 2.0}
	// radii = []float64{0.9, 1.0, 1.1}
	// radii = []float64{1.45, 1.5, 1.6, 1.4}
	//radii = []float64{0.5, 1.0, 1.5}
	radii = []float64{1.0, 5.0, 7.91, 10.0}
	radii = []float64{10.0}
	// radii = []float64{10.0, 20.0}
	// 	radii = []float64{50.0, 100.0, 150.0, 200.0}
	// radii = []float64{100.0, 200.0, 300.0}

	for li, location := range locations {
		for _, radius := range radii {
			circlePoints := location.CircleOnSphere(earthRadiusMiles, radius, 6)
			circlePath := makePolyLine(circlePoints)
			circlePath.ClosePath = true
			setPolyLineColors(circlePath, colors[li])
			m.AddPath(circlePath)
		}
	}

	fmt.Println(m.Encode(true))
}
