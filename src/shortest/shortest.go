package main

import (
	"flag"
	"fmt"
	g "graph"
	itree "immutable_tree"
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

// NamedLocation

type Airport struct {
	sphere.NVector
	name string
}

func (a *Airport) String() string {
	return a.name
}

func (a1 *Airport) CompareTo(c2 itree.Comparable) int {
	if a2, ok := c2.(*Airport); !ok {
		panic("comparing airport against other type of data")
	} else if a1.LessThan(&a2.NVector) {
		return -1
	} else if a2.LessThan(&a1.NVector) {
		return 1
	}
	return 0
}

type AirportIntersection struct {
	sphere.NVector
	airports [2]*Airport
}

func (a *AirportIntersection) String() string {
	return a.airports[0].name + "/" + a.airports[1].name
}

// flightState

type flightState struct {
	remainingRange   float64
	fullRange        float64
	visitedAirports  *itree.Tree
	justSeenAirports [2]*Airport
}

func newFlightState(remainingRange, fullRange float64, tree *itree.Tree, airportsIn ...*Airport) flightState {
	if tree == nil {
		tree = itree.NewTree()
	}

	var airports [2]*Airport
	if len(airportsIn) >= 1 {
		airports[0] = airportsIn[0]
	}
	if len(airportsIn) == 2 {
		airports[1] = airportsIn[1]
	}

	return flightState{remainingRange, fullRange, tree, airports}
}

func (fs *flightState) pushAirports(newAirports ...*Airport) (ok bool) {
	for _, a := range newAirports {
		if fs.visitedAirports.HasValue(a) {
			return false
		}
	}

	for i, a := range fs.justSeenAirports {
		if a != nil {
			fs.visitedAirports = fs.visitedAirports.AddValue(a)
			fs.justSeenAirports[i] = nil
		}
	}

	for i, a := range newAirports {
		fs.justSeenAirports[i] = a
	}

	return true
}

func (fs flightState) TraverseStateHelper(v *g.Vertex) (newState g.PrivateTraverseState, ok bool) {
	var newFs flightState

	newFs.fullRange = fs.fullRange
	newFs.visitedAirports = fs.visitedAirports
	newFs.justSeenAirports = fs.justSeenAirports

	if v.Cost > fs.remainingRange {
		return fs, false
	}

	if airport, isAirport := v.To.Record.(*Airport); isAirport {
		newFs.remainingRange = fs.fullRange
		if ok := newFs.pushAirports(airport); !ok {
			return fs, false
		}
	} else if intersection, isIntersection := v.To.Record.(*AirportIntersection); isIntersection {
		newFs.remainingRange = fs.remainingRange - v.Cost
		if ok := newFs.pushAirports(intersection.airports[0], intersection.airports[1]); !ok {
			return fs, false
		}
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

	// reader := bufio.NewReader(in)

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

			fs := newFlightState(planeRange, planeRange, nil)

			route, distance, ok := graph.Traverse(fs, airportFrom, airportTo)

			if ok {
				fmt.Printf("%0.3f\n", distance)
				if DEBUG&PRINT_ROUTE != 0 {
					for _, n := range route {
						fmt.Println(n.Record.String())
					}
				}
			} else {
				fmt.Println("impossible")
			}
		}
	}
}
