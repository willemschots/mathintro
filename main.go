package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/willemschots/mathintro/polynomial"
)

func main() {
	pts, err := argsToPoints()
	if err != nil {
		log.Fatalf("failed to process input: %v", err)
	}

	p, err := polynomial.Interpolate(pts...)
	if err != nil {
		log.Fatalf("failed to interpolate polynomial: %v", err)
	}

	fmt.Println(p)
}

func argsToPoints() ([]polynomial.Point, error) {
	pts := os.Args[1:]
	if len(pts) == 0 {
		return nil, fmt.Errorf("expected atleast one point")
	}

	results := []polynomial.Point{}
	for _, pt := range pts {
		parts := strings.Split(pt, ",")
		coords := []float64{}
		for _, part := range parts {
			part = strings.TrimSpace(part)
			coord, err := strconv.ParseFloat(part, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse float: %v", part)
			}
			coords = append(coords, coord)
		}
		if len(coords) != 2 {
			return nil, fmt.Errorf("expected 2 elements per coordinate got: %v", len(coords))
		}

		results = append(results, polynomial.Point{X: coords[0], Y: coords[1]})
	}

	return results, nil
}
