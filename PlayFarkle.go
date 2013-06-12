package main

import (
	"fmt"
	"flag"
	"os"
	"github.com/arbrown/gofarkle/farkle"
	"github.com/arbrown/gofarkle/util"
	"github.com/arbrown/gofarkle/dice"
	"github.com/arbrown/gofarkle/ai"
	"math/rand"
)

var debug, verbose bool

func main() {
	// Parse command line flags
	flag.BoolVar( &verbose, "v", false, "whether game info should be printed")
	flag.BoolVar(&debug, "d", false, "whether debug info should be printed")
	var numGames = flag.Int("games", 1, "how many games to play")
	var help = flag.Bool("help", false, "prints this help message")
	var score = flag.Int("score", 10000, "target score threshold")
	var seed = flag.String("seed", "", "play with a specific game seed")

	flag.Parse()

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

	//////////////////////
	game.Players = make([]farkle.FarkleDecider, 0)
	players := make([]farkle.FarkleDecider, 0)
	//////////////////

	////////
	game.PlayerNames = make([]string, 0)
	names := make([]string, 0)
	////////

	wins := make([]int, 0)
	
	var playerNames = flag.Args()

	for _,s := range playerNames {
		
		ai, err := getAi(s)

		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
		game.Players = append(game.Players, ai)
		players = append(players, ai)

		//////////
		game.PlayerNames = append(game.PlayerNames,s)
		names = append(names, s)
		/////////
		wins = append(wins, 0)
	}


	// if a seed was specified, randomize the dice just once
	if *seed != "" {
		rand.Seed(util.Hash(*seed))
	} else {
		dice.Randomize()
	}

	var gameRules farkle.GamePlayer

	gameRules = farkle.TournamentRules { 
		Score: *score,
		Verbose: verbose,
		Debug: debug,
		PlayerNames:names,
	}

	for i:=0;i<*numGames;i++ {

		winner := gameRules.GamePlay(players)

		if debug {
			fmt.Printf("Player %d is the winner with %d\n", winner, game.PlayerScores[winner])
		}

		wins[winner]++
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

func getAi(name string) (player farkle.FarkleDecider, err error) {

	switch name {
	case "TerribleAi":
		return ai.TerribleAi { TargetScore:250 }, nil
	case "TerribleAi2":
		return ai.TerribleAi { TargetScore:200}, nil
	case "GreedyAi":
		return ai.TerribleAi { TargetScore:600}, nil
	case "Human":
		return ai.Human { PrintDice:!verbose }, nil

	}
	
	return nil, fmt.Errorf("Error, '%s' is not a recognized Farkle AI\n", name)
}



func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "Following all flags specify which AIs are playing e.g.:")
	fmt.Fprintf(os.Stderr, "%s -v -games=10 TerribleAi DecentAi\n", os.Args[0])
}
