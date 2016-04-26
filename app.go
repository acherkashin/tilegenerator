package main

import (
	"fmt"
	"log"
	"github.com/terrafactory/tilegenerator/geo"
)

func main() {
	geom := geo.BaseGeometry{}
	line, err := geom.FromWKT("LINESTRING (0 0, 10 10, 20 20)")
	if err != nil {
		log.Fatal(err)
	}

	buf, err := line.Buffer(2.5)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(buf)
}
