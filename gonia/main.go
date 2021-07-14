package main

import (
	"os"
	"fmt"

	"github.com/barrack-obama/gonia"
)

func Usage(concrete error) {
	fmt.Println()
	fmt.Printf("\033[1;31m%s\033[0m", concrete)
	fmt.Println()
	fmt.Println("usage: ./gonia /path/to/beatmap.osu [score]s +[mods (int)]")
	fmt.Println("example: ./gonia /home/simon/beatmaps/2220863.osu 993344s +64")
	fmt.Println()
	fmt.Println("`path to beatmap` parameter is case sensitive, so if it says")
	fmt.Println("`parse error: file not found`, check if your spelling is correct.")
	fmt.Println("also score and mods are optional parameters, and the pp will display")
	fmt.Println("as a perfect score with no mods.")
	fmt.Println()
}

func main() {
	gonia := &gonia.Gonia{}

	err := gonia.ParseConf(os.Args[1:])

	if err != nil {
		Usage(err)
		return
	}

	_, err = gonia.Parse(gonia.Conf.Path, gonia.Conf.Score, gonia.Conf.Mods)

	if err != nil {
		Usage(err)
		return
	}

	gonia.CalculateStars()
	gonia.CalculatePP()

	fmt.Printf("Stars: %f\n", gonia.Stars)
	fmt.Printf("PP: %fpp\n", gonia.PP.Total)
}