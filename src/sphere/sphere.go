package sphere

import (
	"fmt"
	"math"
)

type NVector [3]float64

//func (v1 *NVector) LessThan(v2 *NVector) bool {
//	return (*v1).LessThan(*v2)
//}

func (v1 *NVector) LessThan(v2 *NVector) bool {
	if v1[0] < v2[0] {
		return true
	} else if v1[0] > v2[0] {
		return false
	} else if v1[1] < v2[1] {
		return true
	} else if v1[1] > v2[1] {
		return false
	} else if v1[2] < v2[2] {
		return true
	}
	return false
}

func NewNVectorFromLatLong(lat, lon float64) (result *NVector) {
	result = new(NVector)
	clat := math.Cos(lat)
	result[0] = clat * math.Cos(lon)
	result[1] = clat * math.Sin(lon)
	result[2] = math.Sin(lat)
	return
}

func NewNVectorFromLatLongDeg(lat, lon float64) (result *NVector) {
	return NewNVectorFromLatLong(DegreeToRadian(lat), DegreeToRadian(lon))
}

func DegreeToRadian(degree float64) float64 {
	return degree * math.Pi / 180.0
}

func (v1 *NVector) AngleBetween(v2 *NVector) float64 {
	return math.Atan2(v1.CrossProduct(v2).Magnitude(), v1.DotProduct(v2))
}

func (v1 *NVector) CrossProduct(v2 *NVector) (r *NVector) {
	r = new(NVector)
	r[0] = v1[1]*v2[2] - v1[2]*v2[1]
	r[1] = v1[2]*v2[0] - v1[0]*v2[2]
	r[2] = v1[0]*v2[1] - v1[1]*v2[0]
	return
}

func (v1 *NVector) DotProduct(v2 *NVector) float64 {
	return v1[0]*v2[0] + v1[1]*v2[1] + v1[2]*v2[2]
}

func (v *NVector) Magnitude() float64 {
	return math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
}

func (v *NVector) String() string {
	return fmt.Sprintf("(%g, %g, %g)", v[0], v[1], v[2])
}

func (v *NVector) ScaleBy(factor float64) (r *NVector) {
	r = new(NVector)
	r[0] = v[0] * factor
	r[1] = v[1] * factor
	r[2] = v[2] * factor
	return
}

func (v *NVector) ScaleTo(mag float64) (r *NVector) {
	factor := mag / v.Magnitude()
	return v.ScaleBy(factor)
}

func (v *NVector) Normalize() (r *NVector) {
	return v.ScaleTo(1.0)
}

func (v1 *NVector) Add(v2 *NVector) (r *NVector) {
	r = new(NVector)
	r[0] = v1[0] + v2[0]
	r[1] = v1[1] + v2[1]
	r[2] = v1[2] + v2[2]
	return
}
