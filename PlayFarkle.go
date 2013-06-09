package main

import (
	"fmt"
	"flag"
	"os"
	"github.com/arbrown/gofarkle/farkle"
	"math/rand"
)

var debug, verbose bool

func main() {
	// Parse command line flags
	flag.BoolVar( &verbose, "v", false, "whether game info should be printed")
	flag.BoolVar(&debug, "d", false, "whether debug info should be printed")
	var numGames = flag.Int("games", 1, "how many games to play")
	var help = flag.Bool("help", false, "prints this help message")
	var randomOrder = flag.Bool("rand", false, "whether player order should be shuffled each game")
		var score = flag.Int("score", 10000, "target score threshold")

	flag.Parse()

	var score_threshold = *score

	// verbose output is a subset of debug output
	// damn shame I can't just or-equal them
	//*verbose |= *debug
	if debug {
		verbose = true
	}

	if *help {
		usage()
		os.Exit(2)
	}

	game := new(farkle.GameState)

	game.Players = make([]farkle.FarkleDecider, 0)
	game.PlayerNames = make([]string, 0)
	
	var playerNames = flag.Args()

	for _,s := range playerNames {
		
		ai, err := getAi(s)

		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
		game.Players = append(game.Players, ai)
		game.PlayerNames = append(game.PlayerNames,s)
	}

	numPlayers := len(game.Players) 


	for i:=0;i<*numGames;i++ {
		if verbose {
			var s string
			if numPlayers != 1 {
				s = "s";
			}

			// Lack of ternary... not cool
			fmt.Printf(" Starting game with %d player%s:\n", numPlayers, s)
			for i,p := range game.PlayerNames {
				fmt.Printf("  %d.)\t%s\n", i,p)
			}
		}


		// Shuffle players if necessary
		if *randomOrder {
			if debug {
				fmt.Println("Shuffling player order...")

			}
			n := numPlayers
			// Do the Knuth Shuffle!
			for n > 0 {
				k := rand.Intn(n)
				n--
				game.Players[n], game.Players[k] = game.Players[k], game.Players[n]
				game.PlayerNames[n], game.PlayerNames[k] = game.PlayerNames[k], game.PlayerNames[n]
			}
		}
		
		var player int

		n := numPlayers
		game.PlayerScores = make([]int, n)
		game.PlayerFarkles = make([]int, n)

		// While nobody has hit the target score
		for player = 0 ;checkScores(game.PlayerScores, score_threshold); player = (player+1) % n {
			takeTurn(game, player)
			if debug && player == numPlayers {
				fmt.Println("Current scores:")
				for j,k := range game.PlayerNames {
					fmt.Printf("%s\t%d", k, game.PlayerScores[j])
				}
			}
		}

		var firstOut = i
		// give everyone else one last chance
		for player = (player + 1) % n; i != firstOut; player = (player + 1) % n {
			takeTurn(game, player)
		}
	}

}

func getAi(name string) (ai farkle.FarkleDecider, err error) {

	switch name {
	case "TerribleAi":
		return farkle.TerribleAi { TargetScore:250 }, nil

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

func takeTurn(game *farkle.GameState, player int)  {
	name := game.PlayerNames[player]
	ai := game.Players[player]
	if verbose {
		fmt.Printf(" %s starting turn\n", name)
	}

	// Change score / farkles as necessary
	rollAgain := true
	numDice := 6
	var runScore int
	for rollAgain {
		var keepers []bool
		dice := rollDice(numDice)
		keepers, rollAgain = ai.FarkleDecide(dice, runScore, *game, player)
		rollScore := score(dice, keepers)
		switch {
		case rollScore == 0:
			game.PlayerFarkles[player] += 1
		case game.PlayerFarkles[player] >= 3:
			rollScore = -1000
			game.PlayerFarkles[player] = 0
			rollAgain = false
		}

		runScore += rollScore
	}

	game.PlayerScores[player] += runScore
	if verbose {
		fmt.Printf(" %s scored %d", name, runScore)
	}
}

func rollDice(numDice int) (dice []int)	 {
	dice = make([]int, numDice)
	for i:=0;i<numDice;i++ {
		dice[i] = rand.Intn(6) + 1
	}

	if verbose {
		fmt.Printf(" Rolled: %v\n", dice)
	}

	return dice
}

func score(dice []int, keepers []bool) (score int){

	var kept int
	for _, b := range keepers {
		if b{
			kept++
		}
	}

	// Ideally this function should verify that the 'kept'
	// dice were kept legally.  Since slices are passed by
	// reference, we can just modify the keepers slice
	// to reflect which dice were actually kept (i.e. legally)
	if verbose {
		fmt.Printf(" Keeping ")
		var count int
		for i,b := range keepers {
			switch {
			case b == false:
				break
			case b == true:
				count++
				fmt.Printf("%d", dice[i])
			case count < kept:
				fmt.Printf(", ")
			}
		}
	}

	//score = 0


	return
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "Following all flags specify which AIs are playing e.g.:")
	fmt.Fprintf(os.Stderr, "%s -v -games=10 TerribleAi DecentAi\n", os.Args[0])
}
