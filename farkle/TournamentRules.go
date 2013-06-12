package farkle

import (
	"github.com/arbrown/gofarkle/dice"
	"github.com/arbrown/gofarkle/util"
	"fmt"
)

type TournamentRules struct {
	Verbose bool
	Debug bool
	Score int
	PlayerNames []string


}

func (tr TournamentRules) GamePlay(players []FarkleDecider) (winner_id, turns int) {
	numPlayers := len(players)

	game := new(GameState)

	game.Players = players

	if tr.PlayerNames != nil {
		game.PlayerNames = tr.PlayerNames
	}
	
	turn := 1

	if tr.Verbose {
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

	var player int

	n := numPlayers
	game.PlayerScores = make([]int, n)
	game.PlayerFarkles = make([]int, n)

	// While nobody has hit the target score
	for player = 0 ;; player = (player+1) % n {
		if tr.Verbose && player == 0 {
			fmt.Printf(" Starting turn #%d\n", turn)
		}
		tr.takeTurn(game, player)
		if tr.Debug && player == numPlayers -1 {
			fmt.Println("Current scores:")
			for j,k := range game.PlayerNames {
				fmt.Printf("%s\t%d\n", k, game.PlayerScores[j])
			}
		}
		if player == numPlayers - 1{
			turn++
		}

		if game.PlayerScores[player] > tr.Score {
			break
		}
	}

	if tr.Verbose {
		fmt.Printf(" %s has reached the target score with %d\n", game.PlayerNames[player], game.PlayerScores[player])
	}

	var firstOut = player
	turns = turn

	if tr.Debug {
		fmt.Printf("firstOut: %s(%d)\n", game.PlayerNames[firstOut], firstOut)
	}
	// give everyone else one last chance
	for player = (player + 1) % n; player != firstOut; player = (player + 1) % n {
		if tr.Debug {
			fmt.Printf("%s taking last turn", game.PlayerNames[player])
		}
		tr.takeTurn(game, player)
	}

	winner_id = util.Maxidx(game.PlayerScores)
	if winner_id != firstOut {
		turns++
	}

	if tr.Verbose {
		fmt.Printf(" %s is the winner\n", game.PlayerNames[winner_id])
		for i,v := range game.PlayerNames {
			fmt.Printf(" %s\t%d\n", v, game.PlayerScores[i])
		}
	}

	return
}

// Check the scores and return true if all are under the threshold
func (tr TournamentRules) checkScores(scores []int, threshold int) bool {
	for _, score:= range scores {
		if score >= threshold {
			if tr.Debug {
				fmt.Printf("Score %d over threshold %d", score, threshold)
			}
			return false
		}
	}
	return true
}

func (tr TournamentRules) takeTurn(game *GameState, player int)  {
	name := game.PlayerNames[player]
	ai := game.Players[player]
	if tr.Verbose {
		fmt.Printf(" %s starting turn\n", name)
	}

	// Change score / farkles as necessary
	rollAgain := true
	numDice := 6
	var runScore int
	for rollAgain {
		var keepers []bool
		dice := dice.RollDice(numDice)
		if tr.Verbose {
			fmt.Printf("  Rolled: %v\n", dice)
		}	
		keepers, rollAgain = ai.FarkleDecide(dice, runScore, *game, player)
		rollScore, kept := tr.score(dice, keepers)
		numDice = numDice - kept
		if rollScore == 0{
			runScore = 0
			game.PlayerFarkles[player] += 1
			rollAgain = false
			if tr.Debug {
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
		if tr.Debug {
			fmt.Printf("Adding %d to current run score of %d\n", rollScore, runScore)
		}
		runScore += rollScore

	}

	game.PlayerScores[player] += runScore
	if tr.Verbose {
		fmt.Printf(" %s scored %d\n", name, runScore)
	}
}


func (tr TournamentRules) score(dice []int, keepers []bool) (score int, diceKept int){
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

	if tr.Verbose {
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

	if tr.Debug {
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
	scoremore := true
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
			scoremore = false
		case util.Cmpslc(diceKeptCount, []int{0,1,1,1,1,1}):
			score += 1500
			scoremore = false
		case util.Cmpslc(diceKeptCount,[]int{0,0,5,0,0,0}):
			score += 1200
		case util.Cmpslc(diceKeptCount,[]int{0,5,0,0,0,0}):
			score += 800		
	}
	if (scoremore){
		score += score4(diceKeptCount)
	}
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
		score += diceKeptCount[4] * 50
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
