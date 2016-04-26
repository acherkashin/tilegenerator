package main

import (
	"fmt"
	"log"
	"github.com/paulsmith/gogeos/geos"
)

func main() {
	line, err := geos.FromWKT("POINT (10 10)")
	if err != nil {
		log.Fatal(err)
	}

	//buf, err := line.Buffer(2.5)
	//if err != nil {
	//	log.Fatal(err)
	//}

	fmt.Println(line.Coords())
}
