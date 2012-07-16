package google_static_map

import (
	"bytes"
)

type Point struct {
	lat, lon float64
}

type PolyLine     []Point

func round(in float64) int32 {
	if in >= 0 {
		in += 0.5
	} else {
		in -= 0.5
	}
	return int32(in)
}

func NewPolyLine() *PolyLine {
	result := PolyLine(make([]Point, 0, 10))
	return &result
}

func (pl *PolyLine) AddPoint(pt Point) {
	*pl = append(*pl, pt)
}

func (pl *PolyLine) AddPointLatLon(lat, lon float64) {
	pl.AddPoint(Point{lat, lon})
}

/*
 * Implements the Polyline algorithm as documented at:
 * https://developers.google.com/maps/documentation/utilities/polylinealgorithm
 */
func EncodeSignedFloat(v float64) string {
	i := round(v * 100000) << 1
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

func (pl *PolyLine) Encode(closeLoop bool) string {
	buffer := new(bytes.Buffer)
	prevLat, prevLon := 0.0, 0.0
	for _, p := range []Point(*pl) {
		deltaLat, deltaLon := p.lat - prevLat, p.lon - prevLon
		prevLat, prevLon = p.lat, p.lon
		buffer.WriteString(EncodeSignedFloat(deltaLat))
		buffer.WriteString(EncodeSignedFloat(deltaLon))
	}
	if closeLoop {
		deltaLat, deltaLon := (*pl)[0].lat - prevLat, (*pl)[0].lon - prevLon
		buffer.WriteString(EncodeSignedFloat(deltaLat))
		buffer.WriteString(EncodeSignedFloat(deltaLon))
	}
	return buffer.String()
}