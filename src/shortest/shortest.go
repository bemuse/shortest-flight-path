package main

import (
	"flag"
	"fmt"
	gsm "google_static_map"
	g "graph"
	ipolate "interpolate"
	"math"
	"os"
	"sphere"
)

const (
	EARTH_RADIUS_KM         = 6370.0
	INTERPOLATION_PRECISION = 0.00000001
	DEFAULT_INPUT_FILE      = "sample.in"

	// DEBUG Flags
	READ_AIRPORTS    = 1
	CONNECT_AIRPORTS = 2
	PRINT_ROUTE      = 4
	DEBUG            = 0 | PRINT_ROUTE // | CONNECT_AIRPORTS // | READ_AIRPORTS
)

var inputFileName *string = flag.String("f", DEFAULT_INPUT_FILE, "name of input file")
var verbose *bool = flag.Bool("v", false, "verbose output")
var readNames *bool = flag.Bool("r", false, "read airport names")
var googleMapsURL *bool = flag.Bool("gm", false, "generate Google Maps URL")

type locatable interface {
	Location() sphere.NVector
}

// NamedLocation

type Airport struct {
	sphere.NVector
	name string
}

func (a *Airport) String() string {
	return a.name
}

func (a *Airport) Location() sphere.NVector {
	return a.NVector
}

type AirportIntersection struct {
	sphere.NVector
	airports [2]*Airport
}

func (a *AirportIntersection) String() string {
	return a.airports[0].name + "/" + a.airports[1].name
}

func (a *AirportIntersection) Location() sphere.NVector {
	return a.NVector
}

// flightState

type flightState struct {
	remainingRange float64
	fullRange      float64
}

func newFlightState(remainingRange, fullRange float64) flightState {
	return flightState{remainingRange, fullRange}
}

func (fs flightState) TraverseStateHelper(v *g.Vertex) (newState g.PrivateTraverseState, ok bool) {
	var newFs flightState

	newFs.fullRange = fs.fullRange

	if v.Cost > fs.remainingRange {
		return fs, false
	}

	if _, isAirport := v.To.Record.(*Airport); isAirport {
		newFs.remainingRange = fs.fullRange
	} else if _, isIntersection := v.To.Record.(*AirportIntersection); isIntersection {
		newFs.remainingRange = fs.remainingRange - v.Cost
	} else {
		panic("unknown point in graph traversal")
	}

	return newFs, true
}

func createRoutes(graph *g.Graph, airportNode, intersectionNode *g.Node, midpointNodes *[]*g.Node, maxRadiusKm float64) {
	graph.ConnectBi(airportNode, intersectionNode, maxRadiusKm)
	intersection := intersectionNode.Record.(*AirportIntersection)
	intersectionNVec := &intersection.NVector

	for _, dest := range *midpointNodes {
		otherIntersection := dest.Record.(*AirportIntersection)
		distance := intersectionNVec.AngleBetween(&otherIntersection.NVector) * EARTH_RADIUS_KM
		graph.ConnectBi(dest, intersectionNode, distance)

		if *verbose {
			fmt.Printf("connecting \"%s\" and \"%s\"\n", intersection, otherIntersection)
		}
	}
}

func main() {
	flag.Parse() // Scan the arguments list 

	in, err := os.Open(*inputFileName)
	if err != nil {
		panic("couldn't open input file \"" + *inputFileName + "\"")
	}
	defer func() { in.Close() }()

	for caseNumber := 1; ; caseNumber++ {
		var airportCount int
		var maxRadiusKm float64

		_, err = fmt.Fscan(in, &airportCount, &maxRadiusKm)
		if nil != err {
			break
		}

		fmt.Printf("Case %d:\n", caseNumber)

		airportsByIndex := make([]*g.Node, airportCount+1)
		airportsByName := make(map[string]*g.Node)
		airportRadiusNodes := make(map[*g.Node]*[]*g.Node)
		graph := g.NewGraph()

		for i := 1; i <= airportCount; i++ {
			var lat, lon float64
			_, err = fmt.Fscan(in, &lon, &lat)
			if nil != err {
				panic("couldn't read lon-lat")
			}

			var name string
			names := make([]string, 0)
			if *readNames {
				for _, err = fmt.Fscanf(in, "%q", &name); err == nil; _, err = fmt.Fscanf(in, "%q", &name) {
					names = append(names, name)
				}
			} else {
				name = fmt.Sprintf("Airport %d", i)
				names = append(names, name)
			}

			ap := sphere.NewNVectorFromLatLongDeg(lat, lon)

			airport := Airport{*ap, names[0]}

			node := graph.NewNode(&airport)

			airportsByIndex[i] = node
			for _, name = range names {
				airportsByName[name] = node
			}
			sl := make([]*g.Node, 0)
			airportRadiusNodes[node] = &sl
		}

		radiusAngleRadians := maxRadiusKm / EARTH_RADIUS_KM
		circleRadiusKm := math.Sin(radiusAngleRadians) * EARTH_RADIUS_KM
		circleEarthRadiusKm := math.Cos(radiusAngleRadians) * EARTH_RADIUS_KM

		if *verbose {
			fmt.Printf("circle radius = %f; earth circle radius = %f\n", circleRadiusKm, circleEarthRadiusKm)
		}

		for airport1Node, midpoints1 := range airportRadiusNodes {
			airport1 := airport1Node.Record.(*Airport)
			if DEBUG&CONNECT_AIRPORTS != 0 {
				fmt.Printf("Connecting %q to other airports.\n", airport1.name)
			}
		inner:
			for airport2Node, midpoints2 := range airportRadiusNodes {
				airport2 := airport2Node.Record.(*Airport)

				if !airport1.NVector.LessThan(&airport2.NVector) {
					continue inner
				}

				airportAngle := airport1.NVector.AngleBetween(&airport2.NVector)
				distance := airportAngle * EARTH_RADIUS_KM
				if *verbose {
					fmt.Printf("airports %s and %s are %f km apart consider:%t.\n", airport1.name, airport2.name, distance, distance <= 2*maxRadiusKm)
				}
				if distance <= 2*maxRadiusKm {
					graph.ConnectBi(airport1Node, airport2Node, distance)

					perp := airport1.NVector.CrossProduct(&airport2.NVector).Normalize()
					toMeetV1 := perp.CrossProduct(&airport1.NVector)
					toMeetV2 := airport2.NVector.CrossProduct(perp)

					discCenter1 := airport1.ScaleTo(circleEarthRadiusKm)
					discCenter2 := airport2.ScaleTo(circleEarthRadiusKm)
					discMeetDistance := math.Tan(airportAngle/2) * circleEarthRadiusKm

					discMeetPoint := discCenter1.Add(toMeetV1.ScaleTo(discMeetDistance))
					discMeetPointAlt := discCenter2.Add(toMeetV2.ScaleTo(discMeetDistance)) // REMOVE

					if *verbose {
						fmt.Printf("%s and %s -> %s\n", airport1.name, airport2.name, toMeetV1.String())
						fmt.Printf("    %s==%s\n", discMeetPoint.String(), discMeetPointAlt.String())
					}

					sphereRadiusFunc := func(in float64) float64 {
						return perp.ScaleBy(in).Add(discMeetPoint).Magnitude()
					}

					intersectionFactor1, ok1 := ipolate.Interpolator(0.0, circleRadiusKm, EARTH_RADIUS_KM, INTERPOLATION_PRECISION, sphereRadiusFunc)
					intersectionFactor2, ok2 := ipolate.Interpolator(0.0, -circleRadiusKm, EARTH_RADIUS_KM, INTERPOLATION_PRECISION, sphereRadiusFunc)

					if ok1 && ok2 {
						airportPair := [2]*Airport{airport1Node.Record.(*Airport), airport2Node.Record.(*Airport)}

						intersection1 := perp.ScaleBy(intersectionFactor1).Add(discMeetPoint)
						// REMOVE name1 := fmt.Sprintf("intersection A of %s and %s", airport1.name, airport2.name)
						location1 := AirportIntersection{*intersection1, airportPair}
						node1 := graph.NewNode(&location1)

						intersection2 := perp.ScaleBy(intersectionFactor2).Add(discMeetPoint)
						// REMOVE name2 := fmt.Sprintf("intersection B of %s and %s", airport1.name, airport2.name)
						location2 := AirportIntersection{*intersection2, airportPair}
						node2 := graph.NewNode(&location2)

						createRoutes(graph, airport1Node, node1, midpoints1, maxRadiusKm)
						createRoutes(graph, airport1Node, node2, midpoints1, maxRadiusKm)

						createRoutes(graph, airport2Node, node1, midpoints2, maxRadiusKm)
						createRoutes(graph, airport2Node, node2, midpoints2, maxRadiusKm)

						*midpoints1 = append(*midpoints1, node1, node2)
						*midpoints2 = append(*midpoints2, node1, node2)
					} else if ok1 || ok2 {
						panic("only one point found with two nearby airports")
					} else {
						panic("no points found with two nearby airports")
					}
				} else {
					if *verbose {
						fmt.Println("too great")
					}
				}
			}
		}

		var flightCount int
		fmt.Fscan(in, &flightCount)

		for flight := 0; flight < flightCount; flight++ {
			var planeRange float64
			var airportFrom, airportTo *g.Node

			if *readNames {
				var airportFromName, airportToName string
				_, err = fmt.Fscanf(in, "%q %q %f", &airportFromName, &airportToName, &planeRange)
				if err != nil {
					panic("could not read names of airports on route")
				}
				airportFrom = airportsByName[airportFromName]
				airportTo = airportsByName[airportToName]
			} else {
				var airportFromIndex, airportToIndex int
				_, err = fmt.Fscanf(in, "%d %d %f", &airportFromIndex, &airportToIndex, &planeRange)
				if err != nil {
					panic("could not read indexes of airports on route")
				} else {
					airportFrom = airportsByIndex[airportFromIndex]
					airportTo = airportsByIndex[airportToIndex]
				}
			}

			if *verbose {
				fmt.Printf("from %s to %s with max plane range of %f\n", airportFrom.Record.(*Airport).String(), airportTo.Record.(*Airport).String(), planeRange)
			}

			fs := newFlightState(planeRange, planeRange)

			route, distance, ok := graph.Traverse(fs, airportFrom, airportTo)

			if ok {
				fmt.Printf("%0.3f\n", distance)
				if DEBUG&PRINT_ROUTE != 0 {
					for _, n := range route {
						fmt.Println(n.Record.String())
					}
				}
				if *googleMapsURL {
					gmap := gsm.NewMap(640, 640, 2)
					airportsSeen := make(map[*Airport]bool)
					flightPath := make([]sphere.NVector, 0, len(route))
					for _, n := range route {
						if airport, isAirport := n.Record.(*Airport); isAirport {
							airportsSeen[airport] = true
							lat, lon := airport.NVector.ToLatLonDegrees()
							gmap.AddMarker(gsm.NewPoint(lat, lon))
							flightPath = append(flightPath, airport.NVector)
						} else if intersection, isIntersection := n.Record.(*AirportIntersection); isIntersection {
							airportsSeen[intersection.airports[0]] = true
							airportsSeen[intersection.airports[1]] = true
							// lat, lon := intersection.NVector.ToLatLonDegrees()
							// gmap.AddMarker(gsm.NewPoint(lat, lon))
							flightPath = append(flightPath, intersection.NVector)
						}
					}
					for airport, _ := range airportsSeen {
						pathPoints := airport.NVector.CircleOnSphere(EARTH_RADIUS_KM, maxRadiusKm, 33)
						polyLine := gsm.NewPolyLine()
						polyLine.ClosePath = true
						polyLine.SetWeight(1)
						polyLine.SetColor("0x0000ffff")
						polyLine.SetFillColor("0x8080ff40")
						for _, pp := range pathPoints {
							lat, lon := pp.ToLatLonDegrees()
							polyLine.AddPointLatLon(lat, lon)
						}
						gmap.AddPath(polyLine)
					}
					flightPathPolyLine := makePolyLine(flightPath)
					flightPathPolyLine.SetWeight(1)
					flightPathPolyLine.SetColor("0xff0000ff")
					gmap.AddPath(flightPathPolyLine)
					fmt.Println(gmap.Encode(true))
				}
			} else {
				fmt.Println("impossible")
			}
		}
	}
}

func makePolyLine(points []sphere.NVector) *gsm.PolyLine {
	pl := gsm.NewPolyLine()
	for _, point := range points {
		lat, lon := point.ToLatLonDegrees()
		pl.AddPointLatLon(lat, lon)
	}
	return pl
}
