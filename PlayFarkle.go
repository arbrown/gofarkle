package main

import (
	"fmt"
	"flag"
	"os"
	"github.com/arbrown/gofarkle/farkle"
	"github.com/arbrown/gofarkle/util"
	"github.com/arbrown/gofarkle/dice"
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
	var seed = flag.String("seed", "", "play with a specific game seed")

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
	originalPlayerId := make([]int, 0)
	wins := make([]int, 0)
	
	var playerNames = flag.Args()

	for i,s := range playerNames {
		
		ai, err := getAi(s)

		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
		game.Players = append(game.Players, ai)
		game.PlayerNames = append(game.PlayerNames,s)
		originalPlayerId = append(originalPlayerId, i)
		wins = append(wins, 0)
	}

	numPlayers := len(game.Players) 

	// if a seed was specified, randomize the dice just once
	if *seed != "" {
		rand.Seed(util.Hash(*seed))
	}

	for i:=0;i<*numGames;i++ {
		// Not sure if I REALLY need to reseed
		// random each game...
		if (*seed==""){
			dice.Randomize()
		}

		turn := 1

		if verbose {
			var s string
			if numPlayers != 1 {
				s = "s";
			}

			// Lack of ternary... not cool
			fmt.Printf(" Starting game #%d with %d player%s:\n", i, numPlayers, s)
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
				originalPlayerId[n], originalPlayerId[k] = originalPlayerId[k], originalPlayerId[n]
			}
		}
		
		var player int

		n := numPlayers
		game.PlayerScores = make([]int, n)
		game.PlayerFarkles = make([]int, n)

		// While nobody has hit the target score
		for player = 0 ;; player = (player+1) % n {
			if verbose && player == 0 {
				fmt.Printf(" Starting turn #%d\n", turn)
			}
			takeTurn(game, player)
			if debug && player == numPlayers -1 {
				fmt.Println("Current scores:")
				for j,k := range game.PlayerNames {
					fmt.Printf("%s\t%d\n", k, game.PlayerScores[j])
				}
			}
			if player == numPlayers - 1{
				turn++
			}

			if game.PlayerScores[player] > score_threshold {
				break
			}
		}

		if verbose {
			fmt.Printf(" %s has reached the target score with %d\n", game.PlayerNames[player], game.PlayerScores[player])
		}

		var firstOut = player
		if debug {
			fmt.Printf("firstOut: %s(%d)\n", game.PlayerNames[firstOut], firstOut)
		}
		// give everyone else one last chance
		for player = (player + 1) % n; player != firstOut; player = (player + 1) % n {
			if debug {
				fmt.Printf("%s taking last turn", game.PlayerNames[player])
			}
			takeTurn(game, player)
		}
		
		winner := util.Maxidx(game.PlayerScores)

		if debug {
			fmt.Printf("Player %d is the winner with %d\n", winner, game.PlayerScores[winner])
		}

		wins[originalPlayerId[winner]]++

		if verbose {
			fmt.Printf(" %s is the winner\n", game.PlayerNames[winner])
			for i,v := range game.PlayerNames {
				fmt.Printf(" %s\t%d\n", v, game.PlayerScores[i])
			}
		}
	}

	var s string
	if *numGames != 1 {
		s = "s"
	}
	fmt.Printf("Played %d game%s\n", *numGames, s)
	fmt.Printf("==================\n")
	fmt.Printf("%-15s%s\n", "Player", "Wins")
	for i, w := range wins{
		fmt.Printf("%-15s%d\n", playerNames[i], w)
	}
	

}

func getAi(name string) (ai farkle.FarkleDecider, err error) {

	switch name {
	case "TerribleAi":
		return farkle.TerribleAi { TargetScore:250 }, nil
	case "TerribleAi2":
		return farkle.TerribleAi { TargetScore:200}, nil
	case "GreedyAi":
		return farkle.TerribleAi { TargetScore:600}, nil

	}
	
	return nil, fmt.Errorf("Error, '%s' is not a recognized Farkle AI\n", name)
}

// Check the scores and return true if all are under the threshold
func checkScores(scores []int, threshold int) bool {
	for _, score:= range scores {
		if score >= threshold {
			if debug {
				fmt.Printf("Score %d over threshold %d", score, threshold)
			}
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
		dice := dice.RollDice(numDice)
		if verbose {
			fmt.Printf("  Rolled: %v\n", dice)
		}	
		keepers, rollAgain = ai.FarkleDecide(dice, runScore, *game, player)
		rollScore, kept := score(dice, keepers)
		numDice = numDice - kept
		if rollScore == 0{
			runScore = 0
			game.PlayerFarkles[player] += 1
			rollAgain = false
			if debug {
				fmt.Printf("%s farkle count: %d\n", game.PlayerNames[player], game.PlayerFarkles[player])
			}
		}else {
			game.PlayerFarkles[player] = 0
		}
		if game.PlayerFarkles[player] >= 3{
			runScore = -1000
			game.PlayerFarkles[player] = 0
		}
		if numDice == 0{
			numDice = 6
		}
		if debug {
			fmt.Printf("Adding %d to current run score of %d\n", rollScore, runScore)
		}
		runScore += rollScore

	}

	game.PlayerScores[player] += runScore
	if verbose {
		fmt.Printf(" %s scored %d\n", name, runScore)
	}
}

func score(dice []int, keepers []bool) (score int, diceKept int){
	// Ideally this function should verify that the 'kept'
	// dice were kept legally.  Since slices are passed by
	// reference, we can just modify the keepers slice
	// to reflect which dice were actually kept (i.e. legally)
	
	// We care about the number of dice kept
	var kept int
	// But we care more about the values of the kept dice
	// e.g. 3 1's and 2 5's 
	keptDiceCount := make([]int, 6)
	for i, b := range keepers {
		if b{
			kept++
			keptDiceCount[dice[i]-1]++
		}
	}

	if verbose {
		fmt.Printf("  Keeping ")
		var count int
		for i,b := range keepers {
			if b == true{
				count++
				fmt.Printf("%d", dice[i])
			} else {
				continue
			}
			if count < kept{
				fmt.Printf(", ")
			}
		}
		fmt.Printf("\n")
	}

	if debug {
		fmt.Printf("Kept Dice Count: %d\n", kept)
		fmt.Println(keptDiceCount)
	}

	switch kept {
	case 6:
		score = score6(keptDiceCount)
	case 5:
		score = score5(keptDiceCount)
	case 4:
		score = score4(keptDiceCount)
	case 3:
		score = score3(keptDiceCount)
	case 2:
		score = 100 * keptDiceCount[0] + 50 * keptDiceCount[4]
	case 1:
		score = 100 * keptDiceCount[0] + 50 * keptDiceCount[4]
	}

	diceKept = kept

	return
}

func score6(diceKeptCount []int) (score int) {
		switch {
		case util.Cmpslc(diceKeptCount,[]int{6,0,0,0,0,0}):
			score = 8000
		case util.Cmpslc(diceKeptCount,[]int{0,0,0,0,0,6}):
			score = 4800
		case util.Cmpslc(diceKeptCount,[]int{0,0,0,0,6,0}):
			score = 4000
		case util.Cmpslc(diceKeptCount, []int{1,1,1,1,1,1}):
			score = 3000
		case util.Cmpslc(diceKeptCount,[]int{0,0,0,6,0,0}):
			score = 3200
		case util.Cmpslc(diceKeptCount,[]int{0,0,6,0,0,0}):
			score = 2400
		case util.Cmpslc(diceKeptCount,[]int{0,6,0,0,0,0}):
			score = 1600
		default:
			// check for 3-pair
			if check3Pair(diceKeptCount) {
				score = 750
			} else {
				score = 0
			}
		}
	return score
}

func score5(diceKeptCount []int) (score int) {
	score = 0
	switch {
		case util.Cmpslc(diceKeptCount,[]int{5,0,0,0,0,0}):
			score += 4000
		case util.Cmpslc(diceKeptCount,[]int{0,0,0,0,0,5}):
			score = 2400
		case util.Cmpslc(diceKeptCount,[]int{0,0,0,0,5,0}):
			score += 2000
		case util.Cmpslc(diceKeptCount,[]int{0,0,0,5,0,0}):
			score += 1600
		case util.Cmpslc(diceKeptCount, []int{1,1,1,1,1,0}):
			score += 1500
		case util.Cmpslc(diceKeptCount, []int{0,1,1,1,1,1}):
			score += 1500
		case util.Cmpslc(diceKeptCount,[]int{0,0,5,0,0,0}):
			score += 1200
		case util.Cmpslc(diceKeptCount,[]int{0,5,0,0,0,0}):
			score += 800		
	}
	score += score4(diceKeptCount)
	return score
}

func score4(diceKeptCount []int) (score int) {
	score = 0
	switch {
		case util.Cmpslc(diceKeptCount,[]int{4,0,0,0,0,0}):
			score += 2000
		case util.Cmpslc(diceKeptCount,[]int{0,0,0,0,0,4}):
			score = 1200
		case util.Cmpslc(diceKeptCount,[]int{0,0,0,0,4,0}):
			score += 1000
		case util.Cmpslc(diceKeptCount,[]int{0,0,0,4,0,0}):
			score += 800
		case util.Cmpslc(diceKeptCount,[]int{0,0,4,0,0,0}):
			score += 600
		case util.Cmpslc(diceKeptCount,[]int{0,4,0,0,0,0}):
			score += 400		
	}
	score += score3(diceKeptCount)
	return score
}

func score3(diceKeptCount []int) (score int) {
	score = 0
	switch {
		case util.Cmpslc(diceKeptCount,[]int{3,0,0,0,0,0}):
			score += 1000
		case util.Cmpslc(diceKeptCount,[]int{0,0,0,0,0,3}):
			score = 600
		case util.Cmpslc(diceKeptCount,[]int{0,0,0,0,3,0}):
			score += 500
		case util.Cmpslc(diceKeptCount,[]int{0,0,0,3,0,0}):
			score += 400
		case util.Cmpslc(diceKeptCount,[]int{0,0,3,0,0,0}):
			score += 300
		case util.Cmpslc(diceKeptCount,[]int{0,3,0,0,0,0}):
			score += 200		
	}
	if (diceKeptCount[0] < 3){
		score += diceKeptCount[0] * 100
	}
	if (diceKeptCount[4] < 3){
		score += diceKeptCount[0] * 50
	}
	return score
}

func check3Pair(diceKeptCount []int) bool {
	var pairs int
	for _,v := range diceKeptCount {
		pairs += v/2
	}
	return pairs == 3
}


func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "Following all flags specify which AIs are playing e.g.:")
	fmt.Fprintf(os.Stderr, "%s -v -games=10 TerribleAi DecentAi\n", os.Args[0])
}
