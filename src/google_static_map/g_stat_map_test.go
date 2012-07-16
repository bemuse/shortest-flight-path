package google_static_map

import (
	"fmt"
	"testing"
)

type testFloat struct {
	in  float64
	out string
}

type testUint struct {
	in  uint32
	out string
}

var testSignedFloats = [...]testFloat{{-179.9832104, "`~oia@"}, {38.5, "_p~iF"}, {-120.2, "~ps|U"}, {2.2, "_ulL"}, {-0.75, "nnqC"}, {2.552, "_mqN"}, {-5.503, "vxq`@"}}
var testUnsignedInts = [...]testUint{{174, "mD"}}

func TestSignedFloats(t *testing.T) {
	for _, pair := range testSignedFloats {
		s := EncodeSignedFloat(pair.in)
		if s != pair.out {
			t.Error(fmt.Sprintf("was expecting \"%s\", but got \"%s\"", pair.out, s))
		}
	}
}

func TestUnsignedInts(t *testing.T) {
	for _, pair := range testUnsignedInts {
		s := EncodeUnsignedInt(pair.in)
		if s != pair.out {
			t.Error(fmt.Sprintf("was expecting \"%s\", but got \"%s\"", pair.out, s))
		}
	}
}

func TestPolyLine(t *testing.T) {
	const expected = "_p~iF~ps|U_ulLnnqC_mqNvxq`@"
	pl := NewPolyLine()
	pl.AddPointLatLon(38.5, -120.2)
	pl.AddPointLatLon(40.7, -120.95)
	pl.AddPointLatLon(43.252, -126.453)
	
	s1 := pl.Encode(false)
	if s1 != expected {
		t.Error(fmt.Sprintf("was expecting \"%s\", but got \"%s\"", expected, s1))
	}
	
	s2 := pl.Encode(true)
	if s2 == expected {
		t.Error(fmt.Sprintf("was not expecting \"%s\", yet got \"%s\"", expected, s2))
	}
}
