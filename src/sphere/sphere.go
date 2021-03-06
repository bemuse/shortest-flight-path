package sphere

import (
	"bytes"
	"fmt"
	"math"
)

type NVector [3]float64

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

func NewNVector(x, y, z float64) (result *NVector) {
	result = new(NVector)
	result[0] = x
	result[1] = y
	result[2] = z
	return
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
	return NewNVectorFromLatLong(DegreesToRadians(lat), DegreesToRadians(lon))
}

func DegreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

func RadiansToDegrees(radians float64) float64 {
	return radians * 180 / math.Pi
}

func (v *NVector) ToLatLon() (lat, lon float64) {
	lat = math.Atan2(v[2], math.Sqrt(v[0]*v[0]+v[1]*v[1]))
	lon = math.Atan2(v[1], v[0])
	return
}

func (v *NVector) ToLatLonDegrees() (lat, lon float64) {
	tLat, tLon := v.ToLatLon()
	lat = RadiansToDegrees(tLat)
	lon = RadiansToDegrees(tLon)
	return
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
	return fmt.Sprintf("(%0.5f, %0.5f, %0.5f)", v[0], v[1], v[2])
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

func (v1 *NVector) Subtract(v2 *NVector) (r *NVector) {
	r = new(NVector)
	r[0] = v1[0] - v2[0]
	r[1] = v1[1] - v2[1]
	r[2] = v1[2] - v2[2]
	return
}

func (v1 *NVector) ProjectOnto(onto *NVector) *NVector {
	ontoUnit := onto.Normalize()
	scalarProjection := v1.DotProduct(ontoUnit)
	return ontoUnit.ScaleBy(scalarProjection)
}

func (v1 *NVector) ProjectOntoPlane(normal *NVector) *NVector {
	projectOnto := v1.ProjectOnto(normal)
	return v1.Subtract(projectOnto)
}

/* Assumes vb, v1, and v2 are all on the same plane; if not the results will not be accurate */
func (vb *NVector) IsBetween(v1, v2 *NVector) bool {
	big := v1.AngleBetween(v2)
	small1 := vb.AngleBetween(v1)
	if small1 > big {
		return false
	}
	small2 := vb.AngleBetween(v2)
	return small2 <= big && small1 + small2 <= math.Pi
}

/* Assumes vb, v1, and v2 are all on the same plane; if not the results will not be accurate */
func (vb *NVector) IsBetweenEpsilon(v1, v2 *NVector, epsilon float64) bool {
	big := v1.AngleBetween(v2)
	small1 := vb.AngleBetween(v1)
	small2 := vb.AngleBetween(v2)
	zero := math.Abs(big - small1 - small2)
	return zero <= epsilon
}

func (start *NVector) IsWithin(extreme1, extreme2, between *NVector) bool {
	toExtreme1 := extreme1.Subtract(start)
	toExtreme2 := extreme2.Subtract(start)
	toBetween := between.Subtract(start)
	
	toExtreme1Plane := toExtreme1.ProjectOntoPlane(start)
	toExtreme2Plane := toExtreme2.ProjectOntoPlane(start)
	toBetweenPlane := toBetween.ProjectOntoPlane(start)
	
	return toBetweenPlane.IsBetween(toExtreme1Plane, toExtreme2Plane)
}

func (start *NVector) IsWithinEpsilon(extreme1, extreme2, between *NVector, epsilon float64) bool {
	toExtreme1 := extreme1.Subtract(start)
	toExtreme2 := extreme2.Subtract(start)
	toBetween := between.Subtract(start)
	
	toExtreme1Plane := toExtreme1.ProjectOntoPlane(start)
	toExtreme2Plane := toExtreme2.ProjectOntoPlane(start)
	toBetweenPlane := toBetween.ProjectOntoPlane(start)
	
	return toBetweenPlane.IsBetweenEpsilon(toExtreme1Plane, toExtreme2Plane, epsilon)
}

type Transformation [3][3]float64

func (v *NVector) TransformationMatrix() (result *Transformation) {
	lat, lon := v.ToLatLon()

	sinLon := math.Sin(lon)
	cosLon := math.Cos(lon)
	sinLat := math.Sin(lat)
	cosLat := math.Cos(lat)

	result = new(Transformation)

	result[0][0] = cosLat * cosLon
	result[0][1] = -sinLon
	result[0][2] = -sinLat * cosLon

	result[1][0] = cosLat * sinLon
	result[1][1] = cosLon
	result[1][2] = -sinLat * sinLon

	result[2][0] = sinLat
	result[2][1] = 0.0
	result[2][2] = cosLat

	return
}

func (t *Transformation) Transform(v *NVector) (result *NVector) {
	result = new(NVector)
	for ri, row := range t {
		for ci, col := range row {
			result[ri] += v[ci] * col
		}
	}
	return
}

func (t *Transformation) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString("[")
	for ri, row := range t {
		if ri != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString("[")
		for ci, col := range row {
			if ci != 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(fmt.Sprintf("%0.3f", col))
		}
		buf.WriteString("]")
	}
	buf.WriteString("]")
	return buf.String()
}

func (v *NVector) CircleOnSphere(sphereRadius, surfaceRadius float64, numPoints int) (result []*NVector) {
	result = make([]*NVector, 0, numPoints)
	t := v.TransformationMatrix()
	angle := surfaceRadius / sphereRadius
	centerToCircle := sphereRadius * math.Cos(angle)
	circleRadius := sphereRadius * math.Sin(angle)
	for i := 0; i < numPoints; i++ {
		angle := math.Pi * float64(2*i) / float64(numPoints)
		v := NewNVector(centerToCircle, circleRadius*math.Sin(angle), circleRadius*math.Cos(angle))
		result = append(result, t.Transform(v))
	}

	return
}
