package main

/* 
 * internal use 
 * you can run it through `go run .`
 */

import (
 	"github.com/barrack-obama/gonia"
	"fmt"
)

func main() {
	gonia := &gonia.Gonia{}

	if _, err := gonia.Parse("beatmap.osu", 1000000, 0); err != nil {
		return
	}

	gonia.CalculateStars()
	gonia.CalculatePP()

	fmt.Printf("PP: %f\n", gonia.PP.Total)
	fmt.Printf("Stars: %f\n", gonia.Stars)
}
