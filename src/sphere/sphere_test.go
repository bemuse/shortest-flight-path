package sphere

import (
	"testing"
	"math"
	// "fmt"
	// "os"
)

const (
	earthRadiusKm = 6372.8
	floatEpsilon  = 0.000000001
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

/* Place circles around a variety of places and test each point's distance. */
func TestCircles(t *testing.T) {
	annArbor := NewNVectorFromLatLongDeg(42.281389, -83.748333)
	melbourne := NewNVectorFromLatLongDeg(-37.813611, 144.963056)
	mcmurdoStation := NewNVectorFromLatLongDeg(-77.85, 166.666667) // 750

	f := func(place *NVector, radius float64, pointCount int) {
		points := place.CircleOnSphere(earthRadiusKm, radius, pointCount)
		if len(points) != pointCount {
			t.Errorf("wrong number of points (%d instead of %d)", len(points), pointCount)
		}
		for _, point := range points {
			d := place.AngleBetween(point) * earthRadiusKm
			if math.Abs(radius-d) > floatEpsilon {
				t.Errorf("distance is %f rather than %d", d, radius)
			}
		}
	}

	f(annArbor, 60, 20)
	f(melbourne, 250, 8)
	f(mcmurdoStation, 1750, 13)
}

func TestIsBetween(t *testing.T) {
	stLouis := NewNVectorFromLatLongDeg(38.627222, -90.197778)
	sanJose := NewNVectorFromLatLongDeg(37.335278, -121.891944)
	sanFrancisco := NewNVectorFromLatLongDeg(37.7793, -122.4192)
	sacramento := NewNVectorFromLatLongDeg(38.555556, -121.468889)
	saltLakeCity := NewNVectorFromLatLongDeg(40.75, -111.883333)
	seattle := NewNVectorFromLatLongDeg(47.609722, -122.333056)
	stpaul := NewNVectorFromLatLongDeg(44.9441, -93.0852)
	saline := NewNVectorFromLatLongDeg(42.170833, -83.779722)
	syracuse := NewNVectorFromLatLongDeg(43.046944, -76.144167)
	saratogaSprings := NewNVectorFromLatLongDeg(43.075278, -73.7825)
	stamford := NewNVectorFromLatLongDeg(41.096667, -73.552222)
	silverSpring := NewNVectorFromLatLongDeg(39.004242, -77.019004)
	savannah := NewNVectorFromLatLongDeg(32.081111, -81.091111)
	sarasota := NewNVectorFromLatLongDeg(27.337222, -82.535278)
	sanAntonio := NewNVectorFromLatLongDeg(29.416667, -98.5)
	sanPedro := NewNVectorFromLatLongDeg(33.73583, -118.29139)

	places := []*NVector{sanJose, sanFrancisco, sacramento, saltLakeCity, seattle, stpaul, saline, syracuse, saratogaSprings, stamford, silverSpring,
		savannah, sarasota, sanAntonio, sanPedro}

	for i, _ := range places {
		a := i
		b := (i + 1) % len(places)
		c := (i + 2) % len(places)

		if !stLouis.IsWithin(places[a], places[c], places[b]) {
			t.Errorf("error when i is %d permutation %d", i, 1)
		}
		if !stLouis.IsWithin(places[c], places[a], places[b]) {
			t.Errorf("error when i is %d permutation %d", i, 2)
		}

		if stLouis.IsWithin(places[b], places[a], places[c]) {
			t.Errorf("error when i is %d permutation %d", i, 3)
		}
		if stLouis.IsWithin(places[b], places[c], places[a]) {
			t.Errorf("error when i is %d permutation %d", i, 4)
		}
		if stLouis.IsWithin(places[a], places[b], places[c]) {
			t.Errorf("error when i is %d permutation %d", i, 5)
		}
		if stLouis.IsWithin(places[c], places[b], places[a]) {
			t.Errorf("error when i is %d permutation %d", i, 6)
		}
	}

	for i, _ := range places {
		a := i
		b := (i + 1) % len(places)
		c := (i + 2) % len(places)

		if !stLouis.IsWithinEpsilon(places[a], places[c], places[b], floatEpsilon) {
			t.Errorf("error when i is %d permutation %d", i, 7)
		}
		if !stLouis.IsWithinEpsilon(places[c], places[a], places[b], floatEpsilon) {
			t.Errorf("error when i is %d permutation %d", i, 8)
		}

		if stLouis.IsWithinEpsilon(places[b], places[a], places[c], floatEpsilon) {
			t.Errorf("error when i is %d permutation %d", i, 9)
		}
		if stLouis.IsWithinEpsilon(places[b], places[c], places[a], floatEpsilon) {
			t.Errorf("error when i is %d permutation %d", i, 10)
		}
		if stLouis.IsWithinEpsilon(places[a], places[b], places[c], floatEpsilon) {
			t.Errorf("error when i is %d permutation %d", i, 11)
		}
		if stLouis.IsWithinEpsilon(places[c], places[b], places[a], floatEpsilon) {
			t.Errorf("error when i is %d permutation %d", i, 12)
		}
	}

	if stLouis.IsWithin(saline, sanAntonio, seattle) {
		t.Errorf("opposite test expected failure")
	}

	if stLouis.IsWithinEpsilon(saline, sanAntonio, seattle, floatEpsilon) {
		t.Errorf("opposite test expected success")
	}
}
