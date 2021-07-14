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
	gonai := &gonia.Gonia{}

	if _, err := gonai.Parse("beatmap.osu", 1000000, 0); err != nil {
		return
	}

	gonai.CalculateStars()
	gonai.CalculatePP()

	fmt.Printf("PP: %f\n", gonai.PP.Total)
	fmt.Printf("Stars: %f\n", gonai.Stars)
}