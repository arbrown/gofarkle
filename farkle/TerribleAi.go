// A terrible AI at Farkle
package farkle
import (
	"math"
)
type TerribleAi struct {
	TargetScore int
}

// Keep 1's and 5's until potentialScore is over TargetScore (probably 250)
func (t TerribleAi) FarkleDecide(dice []int, runScore int, game GameState, player int) (keep []bool, rollAgain bool) {

	diceRolled := len(dice)

	var keepers = make([]bool, diceRolled)

	var potentialScore int

	var ones, onesScore, fives, fivesScore int

	for index, roll := range dice {
		if roll == 1 {
			ones++
			keepers[index] = true
		}
		if roll == 5 {
			fives++
			keepers[index] = true
		}
	}

	switch ones {
	case 0,1,2:
		onesScore = ones * 100
	case 3,4,5,6:
		onesScore = (int)(1000 * (math.Exp2((float64)(ones - 3))))
	}

	switch fives {
	case 0,1,2:
		fivesScore = fives * 100
	case 3,4,5,6:
		fivesScore = (int)(500 * (math.Exp2((float64)(fives - 3))))
	}

	potentialScore = runScore + onesScore + fivesScore

	return keepers, potentialScore < t.TargetScore
}
