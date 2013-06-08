package game

import (
	"github.com/arbrown/gofarkle/decider"
)

type GameState struct {
	Players []decider.FarkleDecider
	PlayerNames []string
	PlayerScores []int
	PlayerFarkles []int	
}