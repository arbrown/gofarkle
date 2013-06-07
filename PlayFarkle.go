package main

import (
	"fmt"
	"flag"
	"os"
	"github.com/arbrown/gofarkle/decider"
	"math/rand"
)

func main() {
	// Parse command line flags
	var verbose = flag.Bool("v", false, "whether game info should be printed")
	var debug = flag.Bool("d", false, "whether debug info should be printed")
	var numGames = flag.Int("games", 1, "how many games to play")
	var help = flag.Bool("help", false, "prints this help message")
	var randomOrder = flag.Bool("rand", false, "whether player order should be shuffled each game")
	var score = flag.Int("score", 10000, "target score threshold")

	flag.Parse()

	var score_threshold = *score

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

		// Shuffle players if necessary
		if *randomOrder {
			if *debug {
				fmt.Println("Shuffling player order...")

			}
			playerOrder := make([]decider.FarkleDecider, len(ais))
			perm := rand.Perm(len(ais))

			for i,v := range perm {
				playerOrder[v] = ais[i]
			}
			ais = playerOrder
		}
		
		// Current player id... lines got too long with 'currentPlayer' as the id.
		var i int

		n := numPlayers
		var playerScores = make([]int, n)
		var farkleCount = make([]int, n)

		// While nobody has hit the target score
		for i = 0 ;checkScores(playerScores, score_threshold); i = (i+1) % n {
			playerScore := playerScores[i]
			oppScores := append(playerScores[:i], playerScores[i+1:])
			farkles := farkleCount[i] % 3
			score, farkled := takeTurn(ais[i], playerScore, oppScores, farkles)
			playerScores[i] += score
			if (farkled) {
				farkleCount[i] += 1
			}
		}

		var firstOut = i
		// give everyone else one last chance
		for i = (i + 1) % n; i != firstOut; i = (i + 1) % n {
			playerScore := playerScores[i]
			oppScores := append(playerScores[:i], playerScores[i+1:])
			farkles := farkleCount[i] % 3
			score, farkled := takeTurn(ais[i], playerScore, oppScores, farkles)
			playerScores[i] += score
			if (farkled) {
				farkleCount[i] += 1
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

// Check the scores and return true if all are under the threshold
func checkScores(scores []int, threshold int) bool {
	for _, score:= range scores {
		if score >= threshold {
			return false
		}
	}
	return true
}

func takeTurn(player decider.FarkleDecider, totalScore int, oppScores []int, farkles int) (score int, farkled bool) {
	return 0,false
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "Following all flags specify which AIs are playing e.g.:")
	fmt.Fprintf(os.Stderr, "%s -v -games=10 TerribleAi DecentAi", os.Args[0])
}



