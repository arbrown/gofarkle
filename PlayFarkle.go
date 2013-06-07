package main

import (
	"fmt"
	"flag"
	"os"
	"github.com/arbrown/gofarkle/decider"
)

func main() {
	// Parse command line flags
	var verbose = flag.Bool("v", false, "whether game info should be printed")
	var debug = flag.Bool("d", false, "whether debug info should be printed")
	var numGames = flag.Int("games", 1, "how many games to play")
	var help = flag.Bool("help", false, "prints this help message")

	flag.Parse()

	// verbose output is a subset of debug output
	// damn shame I can't just or-equal them
	//*verbose |= *debug
	if *debug {
		*verbose = true
	}

	if *help {
		usage()
		os.Exit(2)
	}
	
	var ais = make([]decider.FarkleDecider, 0)

	var players = flag.Args()

	for _,s := range players {
		
		ai, err := getAi(s)

		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
		ais = append(ais, ai)
	}

	numPlayers := len(ais) 


	for i:=0;i<*numGames;i++ {
		if *verbose {
			var s string
			if numPlayers != 1 {
				s = "s";
			}

			// Lack of ternary... not cool
			fmt.Printf(" Starting game with %d player%s:\n", numPlayers, s)
			for i,p := range players {
				fmt.Printf("  %d.)\t%s\n", i,p)
			}
		}
	}




}

func getAi(name string) (ai decider.FarkleDecider, err error) {

	switch name {
	case "TerribleAi":
		return decider.TerribleAi { TargetScore:250 }, nil

	}
	
	return nil, fmt.Errorf("Error, '%s' is not a recognized Farkle AI\n", name)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "Following all flags specify which AIs are playing e.g.:")
	fmt.Fprintf(os.Stderr, "%s -v -games=10 TerribleAi DecentAi", os.Args[0])
}



