package sphere

import (
	"math"
	"testing"
)

func TestNormalize(t *testing.T) {
	v1 := NewNVectorFromLatLongDeg(8.0, 31.0)
	v2 := v1.Normalize()
	if math.Abs(v2.Magnitude()-1.0) > 0.0001 {
		t.Errorf("sphere.Normalize does not work")
	}
}

func TestScaleTo(t *testing.T) {
	v1 := NewNVectorFromLatLongDeg(8.0, 31.0)
	v2 := v1.ScaleTo(11.0)
	if math.Abs(v2.Magnitude()-11.0) > 0.0001 {
		t.Errorf("sphere.ScaleTo does not work")
	}
}

// Example from Wikipedia Great Circle article (http://en.wikipedia.org/wiki/Great-circle_distance)
func TestRealWorld(t *testing.T) {
	const earthRadiusKm = 6372.8
	const expectedDistanceKm = 2887.26
	bna := NewNVectorFromLatLongDeg(36.12, -86.67)
	lax := NewNVectorFromLatLongDeg(33.94, -118.40)
	angle1 := bna.AngleBetween(lax)
	angle2 := lax.AngleBetween(bna)
	dist1 := earthRadiusKm * angle1
	dist2 := earthRadiusKm * angle2
	if math.Abs(dist1-expectedDistanceKm) >= 0.005 {
		t.Errorf("distance calculation may be broken")
	} else if math.Abs(dist2-expectedDistanceKm) >= 0.005 {
		t.Errorf("distance calculation may be broken")
	}
}

