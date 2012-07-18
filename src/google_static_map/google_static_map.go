package google_static_map

import (
	"bytes"
	"fmt"
	"net/url"
)

type Location interface {
	GetLocation() string
}

type Point struct {
	lat, lon float64
}

type Address string

type PolyLine struct {
	locations         []Location
	allPointLocations bool
	ClosePath         bool
	Weight            *int
	Color             *string
	FillColor         *string
}

const URL_HEAD = "http://maps.googleapis.com/maps/api/staticmap?"

type Map struct {
	sensor        bool
	markers       []Location
	paths         []*PolyLine
	width, height int
	scale         *int
}

var precedeChar map[bool]string

func init() {
	precedeChar = make(map[bool]string)
	precedeChar[true] = "="
	precedeChar[false] = url.QueryEscape("|")
}

func NewMap(width, height int, scale int) *Map {
	result := new(Map)
	result.sensor = false
	result.width = width
	result.height = height
	if scale != 0 {
		result.scale = &scale
	}
	result.markers = make([]Location, 0, 10)
	result.paths = make([]*PolyLine, 0, 10)
	return result
}

func (m *Map) AddMarker(l Location) {
	m.markers = append(m.markers, l)
}

func (m *Map) AddPath(p *PolyLine) {
	m.paths = append(m.paths, p)
}

func (m *Map) Encode() string {
	buffer := new(bytes.Buffer)

	buffer.WriteString(URL_HEAD)
	buffer.WriteString(fmt.Sprint("sensor=", m.sensor))
	buffer.WriteString(fmt.Sprint("&size=", m.width, "x", m.height))
	if m.scale != nil {
		buffer.WriteString(fmt.Sprint("&scale=", *m.scale))
	}

	for _, path := range m.paths {
		buffer.WriteString("&")
		buffer.WriteString(path.Encode())
	}

	if len(m.markers) > 0 {
		buffer.WriteString("&markers")
		firstParam := true

		for _, marker := range m.markers {
			buffer.WriteString(precedeChar[firstParam])
			firstParam = false
			buffer.WriteString(marker.GetLocation())
		}
	}

	return buffer.String()
}

func round(in float64) int32 {
	if in >= 0 {
		in += 0.5
	} else {
		in -= 0.5
	}
	return int32(in)
}

func NewPoint(lat, lon float64) Point {
	return Point{lat, lon}
}

func (p Point) GetLocation() string {
	return fmt.Sprintf("%0.5f,%0.5f", p.lat, p.lon)
}

func (a Address) GetLocation() string {
	return string(a)
}

func NewPolyLine() *PolyLine {
	result := new(PolyLine)
	result.ClosePath = false
	result.locations = make([]Location, 0, 10)
	result.allPointLocations = true
	return result
}

func (pl *PolyLine) AddPoint(l Location) {
	pl.locations = append(pl.locations, l)
	if _, isPoint := l.(Point); !isPoint {
		pl.allPointLocations = false
	}
}

func (pl *PolyLine) AddPointLatLon(lat, lon float64) {
	pl.AddPoint(Point{lat, lon})
}

/*
 * Implements the Polyline algorithm as documented at:
 * https://developers.google.com/maps/documentation/utilities/polylinealgorithm
 */
func EncodeSignedFloat(v float64) string {
	i := round(v*100000) << 1
	if i < 0 {
		i = -i - 1
	}

	buffer := new(bytes.Buffer)
	for i > 0 {
		b := uint8(i & 0x1F)
		i = i >> 5
		if i > 0 {
			b |= 0x20
		}
		b += 63
		buffer.WriteByte(b)
	}

	return buffer.String()
}

func EncodeUnsignedInt(v uint32) string {
	buffer := new(bytes.Buffer)
	for v > 0 {
		b := uint8(v & 0x1F)
		v = v >> 5
		if v > 0 {
			b |= 0x20
		}
		b += 63
		buffer.WriteByte(b)
	}

	return buffer.String()
}

func (pl *PolyLine) EncodeLocations() string {
	buffer := new(bytes.Buffer)
	prevLat, prevLon := 0.0, 0.0
	for _, loc := range pl.locations {
		pt := loc.(Point)
		deltaLat, deltaLon := pt.lat-prevLat, pt.lon-prevLon
		prevLat, prevLon = pt.lat, pt.lon
		buffer.WriteString(EncodeSignedFloat(deltaLat))
		buffer.WriteString(EncodeSignedFloat(deltaLon))
	}
	if pl.ClosePath {
		pt := pl.locations[0].(Point)
		deltaLat, deltaLon := pt.lat-prevLat, pt.lon-prevLon
		buffer.WriteString(EncodeSignedFloat(deltaLat))
		buffer.WriteString(EncodeSignedFloat(deltaLon))
	}
	return buffer.String()
}

func (pl *PolyLine) Encode() string {
	buffer := new(bytes.Buffer)
	buffer.WriteString("path")
	firstParam := true

	if pl.Weight != nil {
		buffer.WriteString(precedeChar[firstParam])
		firstParam = false
		buffer.WriteString("weight:")
		buffer.WriteString(fmt.Sprint(*pl.Weight))
	}

	if pl.Color != nil {
		buffer.WriteString(precedeChar[firstParam])
		firstParam = false
		buffer.WriteString("color:")
		buffer.WriteString(*pl.Color)
	}

	if pl.FillColor != nil {
		buffer.WriteString(precedeChar[firstParam])
		firstParam = false
		buffer.WriteString("fillcolor:")
		buffer.WriteString(*pl.FillColor)
	}

	if pl.allPointLocations {
		buffer.WriteString(precedeChar[firstParam])
		firstParam = false
		buffer.WriteString("enc:")
		buffer.WriteString(url.QueryEscape(pl.EncodeLocations()))
	} else {
		for _, loc := range pl.locations {
			buffer.WriteString(precedeChar[firstParam])
			firstParam = false
			buffer.WriteString(url.QueryEscape(loc.GetLocation()))
		}
	}

	return buffer.String()
}
