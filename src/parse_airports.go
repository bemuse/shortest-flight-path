package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	ABBREV_INDEX = 2
	NAME_INDEX   = 11
	LAT_INDEX    = 22
	LON_INDEX    = 24
)

type airport struct {
	latitude, longitude float64
	name, nickname      string
}

// The flag package provides a default help printer via -h switch
var inputFile *string = flag.String("i", "NfdcFacilities.csv", "Input CSV file name.")
var outputFile *string = flag.String("o", "usairports.in", "Output file name.")

func parseLatLon(in string) float64 {
	var degrees, minutes int
	var seconds float64
	var direction string
	_, err := fmt.Sscanf(in, "%d-%d-%7f%s", &degrees, &minutes, &seconds, &direction)
	if err != nil {
		panic(fmt.Sprintf("Error parsing latitude or longitude -- %s.", err))
	}

	coord := float64(degrees) + (float64(minutes)+seconds/60.0)/60.0
	switch direction {
	case "N", "E":
	case "S", "W":
		coord = -coord
	default:
		panic("unknown direction")
	}
	return coord
}

func main() {
	flag.Parse() // Scan the arguments list 

	in, err1 := os.Open(*inputFile)
	if err1 != nil {
		panic("couldn't open input file \"" + *inputFile + "\"")
	}
	defer func() { in.Close() }()

	csvReader := csv.NewReader(in)
	csvReader.TrailingComma = true

	airports := make([]airport, 0)

	// skip first line, which contains field names
	fields, err := csvReader.Read()
	for {
		fields, err = csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			break
		}

		latitude := parseLatLon(fields[LAT_INDEX])
		longitude := parseLatLon(fields[LON_INDEX])
		fmt.Printf("%s : %s : (%f, %f)\n", fields[ABBREV_INDEX], fields[NAME_INDEX], latitude, longitude)

		ap := airport{latitude, longitude, fields[NAME_INDEX], fields[ABBREV_INDEX]}
		airports = append(airports, ap)
	}

	out, err2 := os.Create(*outputFile)
	if err2 != nil {
		panic("couldn't open output file \"" + *outputFile + "\"")
	}
	defer func() { out.Close() }()

	fmt.Fprintf(out, "%d %d\n", len(airports), 750)
	for _, ap := range airports {
		fmt.Fprintf(out, "%f %f %q %q\n", ap.longitude, ap.latitude, ap.name, ap.nickname)
	}
	
	fmt.Fprintf(out, "1\n")
	fmt.Fprintf(out, "%q %q 1000\n", "LAX", "LGA")
}
