package farkle

type FarkleDecider interface {
	FarkleDecide(dice []int, runScore int, game GameState, player int) (keep []bool, rollAgain bool)
}
